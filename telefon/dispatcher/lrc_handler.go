package telefon

import (
	"strings"

	"github.com/pkg/errors"
)

type calculateFn func(seed byte, data []byte, dataLength int) ([]byte, error)

func (d *Dispatcher) initLRCHandler() {
	d.registerLRCHandler("XOR", d.calculateXOR)
}

func (d *Dispatcher) registerLRCHandler(name string, cb calculateFn) {

	name = strings.ToUpper(name)

	d.lrcGuard.Lock()
	d.lrcHandler[name] = cb
	d.lrcGuard.Unlock()
}

func (d *Dispatcher) getLRCHandler(name string) calculateFn {

	name = strings.ToUpper(name)

	d.lrcGuard.RLock()
	handler := d.lrcHandler[name]
	d.lrcGuard.RUnlock()

	return handler
}

func (d *Dispatcher) calculateXOR(seed byte, data []byte, length int) ([]byte, error) {

	ret := []byte{}

	if length > len(data) {
		return ret, errors.Errorf("data length: %d is smaller than the length specification: %d", len(data), length)
	}

	calculated := seed

	for i := 0; i < length; i++ {
		char := data[i]
		calculated ^= char
	}

	ret = append(ret, calculated)

	return ret, nil
}
