package syncstate

type State int32

const (
	Start State = iota
	End
	Error
	Cancel
)

type Event struct {
	Addr  string
	State State
	Error error
}

func NewEvent(addr string, state State, err error) *Event {
	return &Event{Addr: addr, State: state, Error: err}
}
