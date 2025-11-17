package templatestate

type State int32

const (
	Changed State = iota
)

type Event struct {
	Slot  uint
	State State
}

func NewEvent(slot uint, state State) *Event {
	return &Event{Slot: slot, State: state}
}
