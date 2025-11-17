package vc3000

// https://godoc.org/gopkg.in/go-playground/validator.v9

type def struct {
	Field    string
	Validate string
	Required bool
}

var (
	cmd2800 = map[byte][]def{

		// check in guest
		'A': {
			def{"R", "required,max=7", true},
			def{"U", "required,max=15,numeric", true},
			def{"O", "required,len=12", true},
			def{"A", "oneof=0 1", false},
		},

		// check out guest
		'B': {
			def{"R", "required,max=7", true},
		},

		// verify guest card
		'E': {
			def{"?", "max=15", false},
			def{"B", "max=40,printascii", false},
		},
	}

	cmdVision = map[byte][]def{

		// check in guest, use -> G/I/H
		'A': {
			def{"R", "omitempty,max=7", true},
			def{"L", "omitempty,max=23", true},
			def{"T", "required,max=15", true},
			def{"U", "required,max=15", true},
			def{"D", "required,len=12", true},
			def{"O", "required,len=12", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"A", "max=53,accesspoint", false},
			def{"C", "max=2,numeric", false},
			def{"P", "max=15", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
		},

		// check out guest
		'B': {
			def{"R", "required,max=7", true},
			def{"T", "required,max=15", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"P", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"?", "max=15", false},
		},

		// change guest status
		'C': {
			def{"R", "required,max=7", true},
			def{"T", "required,max=15", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"P", "max=15", false},
			def{"U", "required,max=15", true},
			def{"D", "required,len=12", true},
			def{"O", "required,len=12", true},
			def{"A", "max=53,accesspoint", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
		},

		// replace guest card
		'D': {
			def{"R", "required,max=7", true},
			def{"T", "required,max=15", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"P", "max=15", false},
			def{"A", "max=53,accesspoint", false},
			def{"C", "max=2,numeric", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
		},

		// verify/read guest card
		'E': {
			def{"?", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false}, // not defined in spec, but must requested (empty) to get the serial number
		},

		// replace guest id
		'F': {
			def{"R", "required,max=7", true},
			def{"T", "required,max=15", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"P", "max=15", false},
			def{"A", "max=53,oneof=0 1", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
		},

		// pre check in guest
		'G': {
			def{"R", "omitempty,max=7", true},
			def{"L", "omitempty,max=23", true},
			def{"T", "required,max=15", true},
			def{"U", "required,max=15", true},
			def{"D", "required,len=12", true},
			def{"O", "required,len=12", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"A", "max=53,accesspoint", false},
			def{"C", "max=2,numeric", false},
			def{"P", "max=15", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
			def{"W", "omitempty,min=2,max=21", false},
		},

		// add guest
		'H': {
			def{"R", "required,max=7", true},
			def{"T", "required,max=15", true},
			def{"U", "required,max=15", true},
			def{"B", "max=40,printascii", false},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"P", "max=15", false},
			def{"O", "len=12", false},
			def{"A", "max=53,accesspoint", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
			def{"W", "omitempty,min=2,max=21", false},
		},

		// check out old guest, check in new guest
		'I': {
			def{"R", "omitempty,max=7", true},
			def{"L", "omitempty,max=23", true},
			def{"T", "required,max=15", true},
			def{"U", "required,max=15", true},
			def{"D", "required,len=12", true},
			def{"O", "required,len=12", true},
			def{"F", "max=15", false},
			def{"N", "max=15", false},
			def{"A", "max=53,accesspoint", false},
			def{"C", "max=2,numeric", false},
			def{"P", "max=15", false},
			def{"1", "max=76", false},
			def{"2", "max=37", false},
			def{"I", "max=128", false},
			def{"?", "max=15", false},
			def{"B", "max=40,printascii", false},
			def{"J", "omitempty,len=1,oneof=1 2 3 4 5", false},
			def{"S", "omitempty,min=8,max=20,hexadecimal", false},
			def{"V", "omitempty,len=8,hexadecimal", false},
			def{"W", "omitempty,min=2,max=21", false},
		},
	}

	commands = map[uint]map[byte][]def{}
)

func init() {
	commands = make(map[uint]map[byte][]def)
	commands[0] = cmdVision
	commands[1] = cmd2800
}
