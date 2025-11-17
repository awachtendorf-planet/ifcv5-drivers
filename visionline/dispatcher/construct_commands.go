package visionline

// https://godoc.org/gopkg.in/go-playground/validator.v9

type def struct {
	Field    string
	Validate string
	Required bool
}

var (
	commands = map[byte][]def{

		// read status for request, alive heartbeat
		'C': {
			def{"EA", "required", true},         // encoder address
			def{"AM", "required,numeric", true}, // answer mode
		},

		// issue card key request
		'A': {
			def{"EA", "required", true},                     // encoder address
			def{"AM", "numeric,oneof=0 1", true},            // answer mode
			def{"CI", "omitempty,len=12", false},            // checkin time
			def{"CN", "", false},                            // counter number
			def{"CO", "len=12", true},                       // checkout time
			def{"CR", "", false},                            // common room (additional rooms)
			def{"GR", "", false},                            // guest room
			def{"IO", "", false},                            // issue operator id
			def{"JR", "omitempty,numeric,oneof=0 1", false}, // joiner
			def{"MK", "", false},                            // mobile key
			def{"NC", "numeric", false},                     // number of cards
			def{"NF", "", false},                            // notification on first entry
			def{"OF", "", false},                            // operator first name
			def{"OP", "", false},                            // operator password
			def{"OS", "", false},                            // operator surname
			def{"OT", "omitempty,numeric,oneof=0 1", false}, // one time
			def{"PD", "omitempty,hexadecimal", false},       // picture data
			def{"PF", "", false},                            // print field
			def{"PP", "", false},                            // pms private
			def{"SD", "omitempty,hexadecimal", false},       // smart card data
			def{"SF", "omitempty,numeric", false},           // smart card file
			def{"SN", "", false},                            // staff rooms, normal opening
			def{"SO", "", false},                            // staff rooms, on off opening
			def{"RS", "omitempty,numeric", false},           // rent safe
			def{"SR", "omitempty,hexadecimal|eq=?", false},  // serial number
			def{"T1", "omitempty,hexadecimal", false},       // track 1
			def{"T2", "", false},                            // track 2
			def{"TA", "", false},                            // terminal address
			def{"UF", "", false},                            // user frist name
			def{"UG", "", false},                            // user group
			def{"UI", "", false},                            // user id
			def{"US", "", false},                            // user surname
		},

		// read card
		'B': {
			def{"EA", "required", true},                    // encoder address
			def{"AM", "numeric,oneof=0 1", true},           // answer mode
			def{"AP", "omitempty,numeric", false},          // auto update print ?
			def{"CI", "omitempty,len=12", false},           // checkin time
			def{"CN", "", false},                           // counter number
			def{"CO", "len=12", true},                      // checkout time
			def{"CR", "", false},                           // common room (additional rooms)
			def{"CT", "omitempty,numeric", false},          // card type
			def{"GR", "", false},                           // guest room
			def{"OF", "", false},                           // operator first name
			def{"OP", "", false},                           // operator password
			def{"OS", "", false},                           // operator surname
			def{"PD", "omitempty,hexadecimal", false},      // picture data
			def{"PF", "", false},                           // print field
			def{"PP", "", false},                           // pms private
			def{"SD", "omitempty,hexadecimal", false},      // smart card data
			def{"SF", "omitempty,numeric", false},          // smart card file
			def{"SN", "", false},                           // staff rooms, normal opening
			def{"SO", "", false},                           // staff rooms, on off opening
			def{"SR", "omitempty,hexadecimal|eq=?", false}, // serial number
			def{"T1", "omitempty,hexadecimal", false},      // track 1
			def{"T2", "", false},                           // track 2
			def{"TA", "", false},                           // terminal address
			def{"TS", "", false},                           // time schedule
			def{"UF", "", false},                           // user frist name
			def{"UG", "", false},                           // user group
			def{"UI", "required", true},                    // user id
			def{"UM", "", false},                           // user message
			def{"US", "", false},                           // user surname
		},

		/*
			// auto-update
			'D': {
				def{"EA", "required", true},                     // encoder address
				def{"AU", "omitempty,numeric", false},           // auto update
				def{"IO", "required", true},                     // issue operator id
				def{"JR", "omitempty,numeric,oneof=0 1", false}, // joiner
				def{"NF", "omitempty,email|phone", false},       // notification on first entry, phone or email
				def{"NG", "omitempty,email|phone", false},       // notify guest, phone or email
				def{"NT", "", false},                            // notification text
				def{"OF", "", false},                            // operator first name
				def{"OS", "", false},                            // operator surname
				def{"PP", "", false},                            // pms private
				def{"SR", "omitempty,hexadecimal|eq=?", false},  // serial number
				def{"T2", "", false},                            // track 2
				def{"TA", "", false},                            // terminal address
				def{"UF", "", false},                            // user frist name
				def{"UG", "", false},                            // user group
				def{"UI", "", false},                            // user id
				def{"UM", "", false},                            // user message
				def{"US", "", false},                            // user surname
			},
		*/

		// change card
		'F': {
			def{"EA", "required", true},                     // encoder address
			def{"AM", "numeric,oneof=0 1", true},            // answer mode
			def{"CO", "omitempty,len=12", false},            // checkout time
			def{"GR", "", false},                            // guest room
			def{"IO", "", false},                            // issue operator id
			def{"JR", "omitempty,numeric,oneof=0 1", false}, // joiner
			def{"NF", "omitempty,email|phone", false},       // notification on first entry, phone or email
			def{"NG", "omitempty,email|phone", false},       // notify guest, phone or email
			def{"NR", "", false},                            // new room
			def{"NT", "", false},                            // notification text
			def{"OF", "", false},                            // operator first name
			def{"OS", "", false},                            // operator surname
			def{"OT", "omitempty,numeric,oneof=0 1", false}, // one time
			def{"PP", "", false},                            // pms private
			def{"RS", "omitempty,numeric", false},           // rent safe
			def{"TA", "", false},                            // terminal address
			def{"UI", "", false},                            // user id
		},

		// checkout guest
		'G': {
			def{"EA", "required", true},          // encoder address
			def{"AM", "numeric,oneof=0 1", true}, // answer mode
			def{"CR", "", false},                 // common room (additional rooms)
			def{"GR", "required", true},          // guest room
			def{"IO", "", false},                 // issue operator id
			def{"MB", "", false},                 // minibar
			def{"PP", "", false},                 // pms private
			def{"SL", "", false},                 // safe lock
			def{"TA", "", false},                 // terminal address
			def{"UI", "", false},                 // user id
		},
	}
)
