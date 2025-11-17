package syncstate

type State int32

const (
	Start State = iota
	End
	Error
)

type Event struct {
	Addr     string
	SyncType int
	State    State
	Error    error
}

func NewEvent(addr string, syncType int, state State, err error) *Event {
	return &Event{Addr: addr, SyncType: syncType, State: state, Error: err}
}
