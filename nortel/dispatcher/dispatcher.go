package nortel

import (
	"regexp"
	"strings"
	"sync"
	"time"

	ifc "github.com/weareplanet/ifcv5-main/ifc/base"
	"github.com/weareplanet/ifcv5-main/ifc/defines"
	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/ifc/record"
	"github.com/weareplanet/ifcv5-main/log"
	tinylru "github.com/weareplanet/ifcv5-main/utils/lru"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// Dispatcher represent the driver dispatcher with generic dispatcher
type Dispatcher struct {
	*dispatcher.Dispatcher
	deviceCache     tinylru.LRU     // addr <-> device type
	regexpCache     tinylru.LRU     // cache for regexp.Compile (station <-> DialledNumberRegexp)
	management      map[uint64]bool // station <-> have management device
	managementGuard sync.RWMutex    // mutex for management
	wake            map[uint64]bool // station <-> wake
	cos             map[uint64]bool // station <-> class of service
}

// NewDispatcher create a new driver dispatcher
func NewDispatcher(network *ifc.Network) *Dispatcher {
	return &Dispatcher{
		dispatcher.NewDispatcher(network),
		tinylru.LRU{},
		tinylru.LRU{},
		make(map[uint64]bool),
		sync.RWMutex{},
		make(map[uint64]bool),
		make(map[uint64]bool),
	}
}

// Startup start the driver and generic dispatcher
func (d *Dispatcher) Startup() {
	log.Info().Msgf("startup %T", d)
	d.deviceCache.Resize(1024)
	d.regexpCache.Resize(1024)
	d.setDefaults()
	d.initOverwrite()
	d.StartupDispatcher()
}

// Close close the driver and generic dispatcher
func (d *Dispatcher) Close() {
	d.CloseDispatcher()
	log.Info().Msgf("shutdown %T", d)
}

// LogIncomingBytes log incoming byte stream
func (d *Dispatcher) LogIncomingBytes(addr string, data []byte) {
	log.Debug().Msgf("rx:%s %d bytes", addr, len(data))
	d.DumpHex(data)
}

// LogOutgoingBytes log outgoing byte stream
func (d *Dispatcher) LogOutgoingBytes(addr string, data []byte) {
	log.Debug().Msgf("tx:%s %d bytes", addr, len(data))
	d.DumpHex(data)
}

// GetIncomingParserSlot returns the required parser slot
// 1 = management device with normal framed serial mode
// 2 = management device with background terminal mode
// 3 = call logging
func (d *Dispatcher) GetIncomingParserSlot(addr string) uint {
	dt := d.getDeviceType(addr)
	if dt.Management {
		if !dt.BackgroundTerminal {
			return 1
		}
		return 2
	}
	return 3
}

// GetOutgoingParserSlot returns the required parser slot
// 1 = management device, ProcessOutgoingBytes recognizes which mode it is
// 0 = no slot defined
func (d *Dispatcher) GetOutgoingParserSlot(addr string) uint {
	dt := d.getDeviceType(addr)
	if dt.Management {
		return 1
	}
	return 0
}

// IsMove ...
func (d *Dispatcher) IsMove(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok {
		if data, exist := guest.GetGeneric(defines.OldRoomName); exist {
			oldRoom := cast.ToString(data)
			if oldRoom != guest.Reservation.RoomNumber {
				return true
			}
		}
	}
	return false
}

// IsSharer
func (d *Dispatcher) IsSharer(context interface{}) bool {
	if guest, ok := context.(*record.Guest); ok && guest.Reservation.SharedInd {
		return true
	}
	return false
}

// IsBackgroundTerminalMode ...
func (d *Dispatcher) IsBackgroundTerminalMode(addr string) bool {
	dt := d.getDeviceType(addr)
	return dt.BackgroundTerminal
}

// GetBackgroundTerminalTimeout ...
func (d *Dispatcher) GetBackgroundTerminalTimeout(_ string) time.Duration {
	return bgdTimeout
}

// GetAnswerTimeout ...
func (d *Dispatcher) GetAnswerTimeout(_ string) time.Duration {
	return packetAnswerTimeout
}

// GetPacketTimeout ...
func (d *Dispatcher) GetPacketTimeout(_ string) time.Duration {
	return packetTimeout
}

// PostInternalCall ...
func (d *Dispatcher) PostInternalCall(station uint64) bool {
	data := d.GetConfig(station, "PostInternalCall", "false")
	return cast.ToBool(data)
}

// PostAllRecords ...
func (d *Dispatcher) PostAllRecords(station uint64) bool {
	data := d.GetConfig(station, "PostAllRecords", "true")
	return cast.ToBool(data)
}

// PostLastRecordOnly ...
func (d *Dispatcher) PostLastRecordOnly(station uint64) bool {
	data := d.GetConfig(station, "PostLastRecordOnly", "false")
	return cast.ToBool(data)
}

// MergeRecords ...
func (d *Dispatcher) MergeRecords(station uint64) bool {
	data := d.GetConfig(station, "MergeRecords", "false")
	return cast.ToBool(data)
}

// SendGuestName ...
func (d *Dispatcher) SendGuestName(station uint64) bool {
	data := d.GetConfig(station, "SendGuestName", "true")
	return cast.ToBool(data)
}

// SendGuestLanguage ...
func (d *Dispatcher) SendGuestLanguage(station uint64) bool {
	data := d.GetConfig(station, "SendGuestLanguage", "true")
	return cast.ToBool(data)
}

// SendClassOfService ...
func (d *Dispatcher) SendClassOfService(station uint64) bool {
	data := d.GetConfig(station, "SendClassOfService", "true")
	return cast.ToBool(data)
}

// GetDialledNumberRegex ...
func (d *Dispatcher) GetDialledNumberRegex(station uint64) (*regexp.Regexp, error) {

	if cache, exist := d.regexpCache.Get(station); exist {
		if re, ok := cache.(*regexp.Regexp); ok {
			return re, nil
		}
		return nil, errors.New("no regexp defined")
	}

	data := d.GetConfig(station, "DialledNumberRegexp", "")
	if len(data) == 0 {
		d.regexpCache.Set(station, nil)
		return nil, errors.New("no regexp defined")
	}

	pattern, err := regexp.Compile(data)
	if err != nil {
		log.Warn().Msgf("station: %d dialled number regexp '%s' failed, err=%s", station, data, err)
	}

	d.regexpCache.Set(station, pattern)
	return pattern, err
}

// SendGuestVIPState ...
func (d *Dispatcher) SendGuestVIPState(station uint64) bool {
	data := d.GetConfig(station, "SendGuestVIPState", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) handleWakeup(station uint64) bool {
	data := d.GetConfig(station, "HandleWakeup", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) sendGuestNameLength(station uint64) bool {
	data := d.GetConfig(station, "SendGuestNameLength", "true")
	return cast.ToBool(data)
}

func (d *Dispatcher) sendRoomVacant(station uint64) bool {
	data := d.GetConfig(station, "SendRoomVacant", "false")
	return cast.ToBool(data)
}

func (d *Dispatcher) getGuestNameLength(station uint64) int {
	data := d.GetConfig(station, "GuestNameLength", "23")
	length := d.getNumeric(data)
	return length
}

func (d *Dispatcher) getExtensionWidth(station uint64) int {
	data := d.GetConfig(station, "ExtensionWidth", "7")
	width := d.getNumeric(data)
	if width < 2 {
		width = 2
	}
	return width
}

func (d *Dispatcher) getManagementDeviceNumber(station uint64) int {
	data := d.GetConfig(station, "ManagmentDevice", "0")
	devNumber := d.getNumeric(data)
	return devNumber
}

func (d *Dispatcher) getNumeric(data string) int {
	data = strings.Trim(data, " ")
	data = strings.TrimLeft(data, "0")
	value := cast.ToInt(data)
	return value
}
