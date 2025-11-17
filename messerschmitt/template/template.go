package template

const (
	PacketAck0 = "ACK0" // outgoing
	PacketAck1 = "ACK1" // outgoing

	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"
	PacketEOT = "EOT" // end of transmission

	PacketWACK = "WACK"
	PacketTTD  = "TTD"

	PacketKeyRequest      = "Create Key"
	PacketKeyDelete       = "Delete Key"
	PacketSynchronisation = "Synchronisation"
	PacketFunctionCards   = "Function Cards"
	PacketSpecialReaders  = "Booking at Special Readers"
	PacketVendingMachines = "Booking at Vending Machines"
	PacketRoomEquipment   = "Booking at Room Equipment"

	PacketCommandAcknowledge = "Command Acknowledge"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"
)

type SerialACK0 struct {
	ACK_ []byte `byte:"len:2,equal:0x1030"`
}

type SerialACK1 struct {
	ACK_ []byte `byte:"len:2,equal:0x1031"`
}

type SerialWACK struct {
	WACK_ []byte `byte:"len:2,equal:0x103b"`
}

type SerialTTD struct {
	TTD_ []byte `byte:"len:2,equal:0x0205"`
}

type SerialNAK struct {
	NAK_ []byte `byte:"len:1,equal:0x15"`
}

type SerialENQ struct {
	ENQ_ []byte `byte:"len:1,equal:0x05"`
}

type SerialEOT struct {
	EOT_ []byte `byte:"len:1,equal:0x04"`
}

type SerialGarbage_ACK0 struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x1030"`
}

type SerialGarbage_ACK1 struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x1031"`
}

type SerialGarbage_NAK struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x15"`
}

type SerialGarbage_ENQ struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x05"`
}

type SerialUnknownPacket1 struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type SerialUnknownPacket2 struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	ETB_  []byte `byte:"len:1,equal:0x17"`
}

type SerialGarbage_Framing_1 struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x02"`
}

type SerialGarbage_Framing_2 struct {
	STX_OVR_ []byte `byte:"len:1,equal:0x02"`
	Data_    []byte `byte:"len:*"`
	STX_     []byte `byte:"len:1,equal:0x02"`
}

// incoming serial

type SerialSynchronisation1 struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	ID_  []byte `byte:"len:2,equal:0x3231"` // 21
	Date []byte `byte:"len:6"`              // ddmmyy
	Time []byte `byte:"len:6"`              // hhmmss
	CR_  []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_ []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialSynchronisation2 struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	ID_  []byte `byte:"len:2,equal:0x3231"` // 21
	Date []byte `byte:"len:6"`              // ddmmyy
	Time []byte `byte:"len:6"`              // hhmmss
	CR_  []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_ []byte `byte:"len:1,equal:0x17"`   // ETB
}

type SerialCommandAcknowledge1 struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	ID_         []byte `byte:"len:2,equal:0x3232"` // 22
	Room        []byte `byte:"len:4"`              // room number
	Encoder     []byte `byte:"len:1"`              // 1..9, A..F
	Status      []byte `byte:"len:1"`              // answer status
	GuestIndex2 []byte `byte:"len:8"`              // 2. block guest index
	TrackData   []byte `byte:"len:37"`             // card data
	GuestIndex1 []byte `byte:"len:8"`              // 1. block guest index
	CR_         []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_        []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialCommandAcknowledge2 struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	ID_         []byte `byte:"len:2,equal:0x3232"` // 22
	Room        []byte `byte:"len:4"`              // room number
	Encoder     []byte `byte:"len:1"`              // 1..9, A..F
	Status      []byte `byte:"len:1"`              // answer status
	GuestIndex2 []byte `byte:"len:8"`              // 2. block guest index
	TrackData   []byte `byte:"len:37"`             // card data
	GuestIndex1 []byte `byte:"len:8"`              // 1. block guest index
	CR_         []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_        []byte `byte:"len:1,equal:0x17"`   // ETB
}

type SerialFunctionCards1 struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	ID_      []byte `byte:"len:2,equal:0x3233"` // 23
	Room     []byte `byte:"len:4"`              // room number
	CardType []byte `byte:"len:2"`              // card type of function
	Data     []byte `byte:"len:8"`              // application data
	Date     []byte `byte:"len:6"`              // ddmmyy
	Time     []byte `byte:"len:6"`              // hhmmss
	CR_      []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_     []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialFunctionCards2 struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	ID_      []byte `byte:"len:2,equal:0x3233"` // 23
	Room     []byte `byte:"len:4"`              // room number
	CardType []byte `byte:"len:2"`              // card type of function
	Data     []byte `byte:"len:8"`              // application data
	Date     []byte `byte:"len:6"`              // ddmmyy
	Time     []byte `byte:"len:6"`              // hhmmss
	CR_      []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_     []byte `byte:"len:1,equal:0x17"`   // ETB
}

type SerialSpecialReaders1 struct {
	STX_             []byte `byte:"len:1,equal:0x02"`
	ID_              []byte `byte:"len:2,equal:0x3234"` // 24
	Reader           []byte `byte:"len:4"`              // reader number the booking was done
	Date             []byte `byte:"len:6"`              // ddmmyy
	Time             []byte `byte:"len:6"`              // hhmmss
	CardType_        []byte `byte:"len:2"`              // card type 51 guest
	CompanyCode_     []byte `byte:"len:5"`              // company code on the card
	Room             []byte `byte:"len:4"`              // room number
	E_               []byte `byte:"len:1"`              // 0 adult, 1 child, 2 young person
	CardIndex_       []byte `byte:"len:1"`              // card index
	Data1_           []byte `byte:"len:2"`              // application data
	CodedAssignment_ []byte `byte:"len:5"`              // coded assignment number on the card
	Data2_           []byte `byte:"len:6"`              // application data
	GuestIndex       []byte `byte:"len:16"`             // guest index
	CR_              []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_             []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialSpecialReaders2 struct {
	STX_             []byte `byte:"len:1,equal:0x02"`
	ID_              []byte `byte:"len:2,equal:0x3234"` // 24
	Reader           []byte `byte:"len:4"`              // reader number the booking was done
	Date             []byte `byte:"len:6"`              // ddmmyy
	Time             []byte `byte:"len:6"`              // hhmmss
	CardType_        []byte `byte:"len:2"`              // card type 51 guest
	CompanyCode_     []byte `byte:"len:5"`              // company code on the card
	Room             []byte `byte:"len:4"`              // room number
	E_               []byte `byte:"len:1"`              // 0 adult, 1 child, 2 young person
	CardIndex_       []byte `byte:"len:1"`              // card index
	Data1_           []byte `byte:"len:2"`              // application data
	CodedAssignment_ []byte `byte:"len:5"`              // coded assignment number on the card
	Data2_           []byte `byte:"len:6"`              // application data
	GuestIndex       []byte `byte:"len:16"`             // guest index
	CR_              []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_             []byte `byte:"len:1,equal:0x17"`   // ETB
}

type SerialVendingMachines1 struct {
	STX_             []byte `byte:"len:1,equal:0x02"`
	ID_              []byte `byte:"len:2,equal:0x3235"` // 25
	Reader           []byte `byte:"len:4"`              // reader number the booking was done
	Date             []byte `byte:"len:6"`              // ddmmyy
	Time             []byte `byte:"len:6"`              // hhmmss
	CardType_        []byte `byte:"len:2"`              // card type 51 guest
	CompanyCode_     []byte `byte:"len:5"`              // company code on the card
	Room             []byte `byte:"len:4"`              // room number
	E_               []byte `byte:"len:1"`              // 0 adult, 1 child, 2 young person
	CardIndex_       []byte `byte:"len:1"`              // card index
	Data1_           []byte `byte:"len:2"`              // application data
	CodedAssignment_ []byte `byte:"len:5"`              // coded assignment number on the card
	Article          []byte `byte:"len:2"`              // acrticle number
	Price            []byte `byte:"len:4"`              // article price
	GuestIndex       []byte `byte:"len:16"`             // guest index
	CR_              []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_             []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialVendingMachines2 struct {
	STX_             []byte `byte:"len:1,equal:0x02"`
	ID_              []byte `byte:"len:2,equal:0x3235"` // 25
	Reader           []byte `byte:"len:4"`              // reader number the booking was done
	Date             []byte `byte:"len:6"`              // ddmmyy
	Time             []byte `byte:"len:6"`              // hhmmss
	CardType_        []byte `byte:"len:2"`              // card type 51 guest
	CompanyCode_     []byte `byte:"len:5"`              // company code on the card
	Room             []byte `byte:"len:4"`              // room number
	E_               []byte `byte:"len:1"`              // 0 adult, 1 child, 2 young person
	CardIndex_       []byte `byte:"len:1"`              // card index
	Data1_           []byte `byte:"len:2"`              // application data
	CodedAssignment_ []byte `byte:"len:5"`              // coded assignment number on the card
	Article          []byte `byte:"len:2"`              // acrticle number
	Price            []byte `byte:"len:4"`              // article price
	GuestIndex       []byte `byte:"len:16"`             // guest index
	CR_              []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_             []byte `byte:"len:1,equal:0x17"`   // ETB
}

type SerialRoomEquipment1 struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	ID_        []byte `byte:"len:2,equal:0x3236"` // 26
	Room       []byte `byte:"len:4"`              // reader number the booking was done
	Date       []byte `byte:"len:6"`              // ddmmyy
	Time       []byte `byte:"len:6"`              // hhmmss
	Device     []byte `byte:"len:1"`              // device address in the room
	Data_      []byte `byte:"len:4"`              // application data
	Article    []byte `byte:"len:2"`              // acrticle number
	Price      []byte `byte:"len:4"`              // article price
	GuestIndex []byte `byte:"len:16"`             // guest index
	CR_        []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_       []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialRoomEquipment2 struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	ID_        []byte `byte:"len:2,equal:0x3236"` // 26
	Room       []byte `byte:"len:4"`              // reader number the booking was done
	Date       []byte `byte:"len:6"`              // ddmmyy
	Time       []byte `byte:"len:6"`              // hhmmss
	Device     []byte `byte:"len:1"`              // device address in the room
	Data_      []byte `byte:"len:4"`              // application data
	Article    []byte `byte:"len:2"`              // acrticle number
	Price      []byte `byte:"len:4"`              // article price
	GuestIndex []byte `byte:"len:16"`             // guest index
	CR_        []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_       []byte `byte:"len:1,equal:0x17"`   // ETB
}

// outgoing serial

type SerialSynchronisationReply struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	ID_  []byte `byte:"len:2,equal:0x3034"` // 04
	Date []byte `byte:"len:6"`              // ddmmyy
	Time []byte `byte:"len:6"`              // hhmmss
	CR_  []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_ []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialKeyCreate struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	ID_           []byte `byte:"len:2,equal:0x3031"` // 01
	Room          []byte `byte:"len:4"`              // room number
	CardType_     []byte `byte:"len:2,equal:0x3531"` // 51 guest card, 56, 57, 99
	Assignment    []byte `byte:"len:1"`              // 0 new, 1 old
	E_            []byte `byte:"len:1,equal:0x30"`   // 0 adult, 1 child, 2 young person
	Suite         []byte `byte:"len:6"`              // suite bits
	AccessPoints  []byte `byte:"len:27"`             // access points
	Cards         []byte `byte:"len:1"`              // number of cards 1..9
	Encoder       []byte `byte:"len:1"`              // 0 or 1..9, A..F
	DepartureDate []byte `byte:"len:6"`              // date of departure ddmmjj
	GuestIndex2   []byte `byte:"len:8"`              // 2. block guest index
	DepartureTime []byte `byte:"len:3"`              // time of departure hhm (minutes in 10' steps)
	Data          []byte `byte:"len:19"`             // optional, free text field
	GuestIndex1   []byte `byte:"len:8"`              // 1. block guest index
	CR_           []byte `byte:"len:1,equal:0x0d"`   // CR
	ETX_          []byte `byte:"len:1,equal:0x03"`   // ETX
}

type SerialKeyDelete struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	ID_         []byte `byte:"len:2,equal:0x3032"`           // 02
	Room        []byte `byte:"len:4"`                        // room number
	CardType_   []byte `byte:"len:2,equal:0x3531"`           // 51 guest card
	Encoder     []byte `byte:"len:1"`                        // 0 or 1..9, A..F
	Data_       []byte `byte:"len:7,equal:0x30303030303030"` // application data 0000000
	Invalidity_ []byte `byte:"len:1,equal:0x30"`             // 0 immediately
	CR_         []byte `byte:"len:1,equal:0x0d"`             // CR
	ETX_        []byte `byte:"len:1,equal:0x03"`             // ETX
}

// incoming tcp

type SocketSynchronisation struct {
	ID_  []byte `byte:"len:2,equal:0x3731"` // 71
	Date []byte `byte:"len:8"`              // ddmmyyyy
	Time []byte `byte:"len:6"`              // hhmmss
}

type SocketCommandAcknowledge struct {
	ID_        []byte `byte:"len:2,equal:0x3732"` // 72
	Room       []byte `byte:"len:10"`             // room number
	Encoder    []byte `byte:"len:2"`              // encoder 0 or 01..32
	Status     []byte `byte:"len:1"`              // answer status
	GuestIndex []byte `byte:"len:16"`             // guest index
	TrackData  []byte `byte:"len:37"`             // card data
}

type SocketFunctionCards struct {
	ID_      []byte `byte:"len:2,equal:0x3733"` // 73
	Room     []byte `byte:"len:10"`             // room number
	CardType []byte `byte:"len:2"`              // card type of function
	Date     []byte `byte:"len:8"`              // ddmmyyyy
	Time     []byte `byte:"len:6"`              // hhmmss
}

type SocketSpecialReaders struct {
	ID_        []byte `byte:"len:2,equal:0x3734"` // 74
	Reader     []byte `byte:"len:10"`             // special area reader
	Date       []byte `byte:"len:8"`              // ddmmyyyy
	Time       []byte `byte:"len:6"`              // hhmmss
	Room       []byte `byte:"len:10"`             // room number
	CardIndex_ []byte `byte:"len:1"`              // card index
	GuestIndex []byte `byte:"len:16"`             // guest index
}

// outgoing tcp

type SocketSynchronisationReply struct {
	ID_  []byte `byte:"len:2,equal:0x3434"` // 44
	Date []byte `byte:"len:8"`              // ddmmyyyy
	Time []byte `byte:"len:6"`              // hhmmss
}

type SocketKeyCreate struct {
	ID_           []byte `byte:"len:2,equal:0x3431"` // 41
	Room          []byte `byte:"len:10"`             // room number
	Encoder       []byte `byte:"len:2"`              // encoder 0 or 01..32
	Cards         []byte `byte:"len:1"`              // number of cards 1..9
	Assignment    []byte `byte:"len:1"`              // 0 new, 1 duplicate
	DepartureDate []byte `byte:"len:8"`              // date of departure ddmmjjjj
	DepartureTime []byte `byte:"len:2"`              // time of departure hh (full hours only)
	C_            []byte `byte:"len:1,equal:0x30"`   // connected room, 0 no connected room, 1 connected room
	GuestIndex    []byte `byte:"len:16"`             // guest index
	AccessPoints  []byte `byte:"len:40"`             // access points
	// optional
	// CheckinDate []byte `byte:"len:8"` // valid date ddmmjjjj
	// CheckinTime []byte `byte:"len:2"` // valid time hh (full hours only)

}

type SocketKeyDelete struct {
	ID_         []byte `byte:"len:2,equal:0x3432"` // 42
	Room        []byte `byte:"len:10"`             // room number
	Encoder     []byte `byte:"len:2"`              // encoder 0 or 01..32
	Invalidity_ []byte `byte:"len:1,equal:0x30"`   // 0 immediately
}
