package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckIn        = "Check In"
	PacketCheckOut       = "Check Out"
	PacketChangeName     = "Change Display Name"
	PacketMessageLamp    = "Set Message Lamp"
	PacketSetRestriction = "Set Restriction"
	PacketRoomStatus     = "Room Status"
	PacketWakeupSet      = "Wakeup Set"
	PacketWakeupClear    = "Wakeup clear"
	PacketCallCounter    = "Call Counter"

	PacketSwapEnquiry = "Database Swap Enquiry"
	PacketSwapStart   = "Database Swap Start"
	PacketSwapEnd     = "Database Swap End"
	PacketAlive       = "Alive"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"
)

type TplACK struct {
	ACK_ []byte `byte:"len:1,equal:0x06"`
}

type TplNAK struct {
	NAK_ []byte `byte:"len:1,equal:0x15"`
}

type TplENQ struct {
	ENQ_ []byte `byte:"len:1,equal:0x05"`
}

type TplGarbage_ACK struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x06"`
}

type TplGarbage_NAK struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x15"`
}

type TplGarbage_ENQ struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x05"`
}

type TplUnknownPacket struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplGarbage_Framing_1 struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x02"`
}

type TplGarbage_Framing_2 struct {
	STX_OVR_ []byte `byte:"len:1,equal:0x02"`
	Data_    []byte `byte:"len:*"`
	STX_     []byte `byte:"len:1,equal:0x02"`
}

// outgoing
type TplAlive struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:10,equal:0x41524559555448455245"` // AREYUTHERE
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCheckin struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x43484b"` // CHK
	Status_   []byte `byte:"len:2,equal:0x3120"`   // 1
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplCheckout struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x43484b"` // CHK
	Status_   []byte `byte:"len:2,equal:0x3020"`   // 0
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplChangeName struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x4e414d"` // NAM
	Status    []byte `byte:"len:1"`                // 1 = add, 2 = replace, 3 = delete
	Padding_  []byte `byte:"len:1,equal:0x20"`
	Name      []byte `byte:"len:21"`
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplMessageLamp struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x4d5720"` // MW
	Status    []byte `byte:"len:1"`                // 1 = on, 0 = off
	Padding_  []byte `byte:"len:1,equal:0x20"`
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplSetRestriction struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x525354"` // RST
	Status    []byte `byte:"len:1"`                // class of service
	Padding_  []byte `byte:"len:1,equal:0x20"`
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplWakeup struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x574b50"` // WKP
	Status    []byte `byte:"len:4"`                // time 24 hour, spaces to delete
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplSwapStart struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	Cmd_     []byte `byte:"len:3,equal:0x475253"` // GRS
	Padding_ []byte `byte:"len:7,equal:0x20202020202020"`
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

type TplSwapEnd struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	Cmd_     []byte `byte:"len:3,equal:0x454e44"` // END
	Padding_ []byte `byte:"len:7,equal:0x20202020202020"`
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

// incomming

type TplSwapEnquiry struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:5,equal:0x5251494e5a"` // RQINZ
	OVR_ []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatus struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x535453"` // STS
	Status    []byte `byte:"len:1"`                // condition
	Padding_  []byte `byte:"len:1,equal:0x20"`
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplCallCounter struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:3,equal:0x4d5220"` // MR
	Status    []byte `byte:"len:4"`                // counter
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}
