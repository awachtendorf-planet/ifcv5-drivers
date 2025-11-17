package guestlink

import (
	"fmt"
	"strings"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/log"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Transaction struct {
	Identifier string
	Sequence   int
}

func (t *Transaction) IsLastSequence() bool {
	return t.Sequence == 9999
}

func (d *Dispatcher) GetTransaction(in *ifc.LogicalPacket) Transaction {

	transaction := Transaction{}

	if in == nil {
		return transaction
	}

	data := in.Data()

	if tr, exist := data[guestlink_tan]; exist {
		transaction.Identifier = cast.ToString(tr)
	}

	if sq, exist := data[guestlink_seq]; exist {
		seq := cast.ToString(sq)
		transaction.Sequence = cast.ToInt(seq)
	}

	return transaction

}

func (d *Dispatcher) NewTransaction(addr string) (Transaction, error) {

	transaction := Transaction{
		Sequence: 9999,
	}

	station, err := d.GetStationAddr(addr)

	if err == nil {

		tan := d.nextTAN(station)

		transaction.Identifier = fmt.Sprintf("M%03d", tan)

	}

	return transaction, err
}

func (d *Dispatcher) nextTAN(station uint64) uint16 {

	d.tanGuard.Lock()

	tan, exist := d.tan[station]

	if !exist {
		counter := d.ReadValue(station, defines.Transaction, "0")
		tan = cast.ToUint16(counter)
	}

	if tan >= 999 || tan < 0 {
		tan = 0
	}

	tan++

	d.tan[station] = tan

	d.StoreValue(station, defines.Transaction, tan)

	d.tanGuard.Unlock()

	return tan

}

func (d *Dispatcher) RegisterTransaction(owner interface{}, addr string, transaction Transaction) error {

	d.transactionGuard.Lock()

	name := d.getOwnerIndex(owner, addr)
	if lastTransaction, exist := d.owner[name]; exist {
		d.transactionGuard.Unlock()
		return errors.Errorf("previous transaction '%s' is still registered", lastTransaction.Identifier)
	}

	idx := d.getTransactionIndex(addr, transaction)

	if _, exist := d.transaction[idx]; exist {
		d.transactionGuard.Unlock()
		return errors.Errorf("transaction '%s' still registered", transaction.Identifier)
	}

	d.owner[name] = transaction
	d.transaction[idx] = owner

	log.Debug().Msgf("%T register transaction '%s' for '%T', addr '%s'", d, transaction.Identifier, owner, addr)

	d.transactionGuard.Unlock()

	return nil
}

func (d *Dispatcher) UnregisterTransaction(owner interface{}, addr string, transaction Transaction) error {

	d.transactionGuard.Lock()

	name := d.getOwnerIndex(owner, addr)

	lastTransaction, exist := d.owner[name]
	if !exist || lastTransaction.Identifier != transaction.Identifier {
		d.transactionGuard.Unlock()
		return errors.Errorf("the transaction '%s' is not yours", transaction.Identifier)
	}

	idx := d.getTransactionIndex(addr, transaction)

	delete(d.owner, name)

	delete(d.transaction, idx)

	log.Debug().Msgf("%T unregister transaction '%s' for '%T', addr '%s'", d, transaction.Identifier, owner, addr)

	d.transactionGuard.Unlock()

	return nil

}

func (d *Dispatcher) Owner(addr string, transaction Transaction) (interface{}, bool) {

	idx := d.getTransactionIndex(addr, transaction)

	d.transactionGuard.RLock()
	owner, exist := d.transaction[idx]
	d.transactionGuard.RUnlock()

	return owner, exist
}

func (d *Dispatcher) LastTransaction(owner interface{}, addr string) (Transaction, bool) {

	name := d.getOwnerIndex(owner, addr)

	d.transactionGuard.RLock()
	transaction, exist := d.owner[name]
	d.transactionGuard.RUnlock()

	return transaction, exist
}

func (d *Dispatcher) FreeTransaction(owner interface{}, addr string) {

	var free []string

	d.transactionGuard.RLock() // read lock

	for k, v := range d.transaction {

		// owner = datachange.Plugin
		// addr  = pcon-DESTR1:2000:1:COM4
		// k 	 = pcon-DESTR1:2000:1:COM4:M010 (addr:transaction)

		if v == owner && len(addr) > 0 && strings.HasPrefix(k, addr) {
			free = append(free, k)
		}
	}

	d.transactionGuard.RUnlock()

	d.transactionGuard.Lock() // write lock

	for i := range free {
		log.Warn().Msgf("%T free transaction '%s' for '%T', addr '%s'", d, free[i], owner, addr)
		delete(d.transaction, free[i])
	}

	name := d.getOwnerIndex(owner, addr)

	delete(d.owner, name)

	d.transactionGuard.Unlock()

}

func (d *Dispatcher) RegisterRoute(owner interface{}, addr string, route *ifc.Route) {

	if route == nil {
		return
	}

	name := d.getOwnerIndex(owner, addr)

	d.routerGuard.Lock()

	d.router[name] = route

	d.routerGuard.Unlock()

}

func (d *Dispatcher) GetRoute(owner interface{}, addr string) (*ifc.Route, bool) {

	name := d.getOwnerIndex(owner, addr)

	d.routerGuard.RLock()
	route, exist := d.router[name]
	d.routerGuard.RUnlock()

	return route, exist
}

func (d *Dispatcher) UnregisterRoute(owner interface{}, addr string) {

	name := d.getOwnerIndex(owner, addr)

	d.routerGuard.Lock()

	delete(d.router, name)

	d.routerGuard.Unlock()

}

func (d *Dispatcher) UnregisterAndCloseRoute(owner interface{}, addr string) {

	name := d.getOwnerIndex(owner, addr)

	d.routerGuard.Lock()

	route, exist := d.router[name]
	delete(d.router, name)

	defer d.routerGuard.Unlock()

	if exist && route != nil {

		defer func() {
			recover() // already closed?
		}()

		if !route.Shared {

			if route.C != nil {
				close(route.C)
			}

		}

	}

}

func (d *Dispatcher) getTransactionIndex(addr string, transaction Transaction) string {
	idx := fmt.Sprintf("%s:%s", addr, transaction.Identifier)
	return idx
}

func (d *Dispatcher) getOwnerIndex(owner interface{}, addr string) string {
	idx := fmt.Sprintf("%T:%s", owner, addr)
	return idx
}
