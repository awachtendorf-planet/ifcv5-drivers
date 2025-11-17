package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketRestart = "Restart Message"
	PacketStartup = "Startup Message"
	PacketUpdate  = "Update Message"

	//PacketChange  = "Change Message"

	PacketCheckIn   = "Check In"
	PacketCheckOut  = "Check Out"
	PacketLockBar   = "Lock Bar"
	PacketUnlockBar = "Unlock Bar"

	PacketUpdateCheckIn   = "Update Check In"
	PacketUpdateCheckOut  = "Update Check Out"
	PacketUpdateLockBar   = "Update Lock Bar"
	PacketUpdateUnlockBar = "Update Unlock Bar"

	PacketResult = "Result Message"
	PacketSale   = "Sale Message"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"
)

type TplRestartPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:1,equal:0x2a"` // *
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplStartupPacket struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Record_        []byte `byte:"len:1,equal:0x53"`   // S
	Version_       []byte `byte:"len:2,equal:0x3032"` // 02
	SequenceNumber []byte `byte:"len:5"`              // ddddd or bbbbb
	Date           []byte `byte:"len:10"`             // YYMMDDhhmm
	Optional_      []byte `byte:"len:1,equal:0x52"`   // R (optional)
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplUpdatePacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:1,equal:0x55"` // U
	Room    []byte `byte:"len:*"`            //
	Command []byte `byte:"len:1"`            // I,O (undocumented L,U)
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplChangePacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:1,equal:0x43"` // C
	Room    []byte `byte:"len:*"`            //
	Command []byte `byte:"len:1"`            // I,O,L,U
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

// type TplResultPacket struct {
// 	STX_    []byte `byte:"len:1,equal:0x02"`
// 	Record_ []byte `byte:"len:1,equal:0x52"` // R
// 	Room    []byte `byte:"len:*"`            //
// 	Command []byte `byte:"len:1"`            // I,O,L,U
// 	Status  []byte `byte:"len:1"`
// 	ETX_    []byte `byte:"len:1,equal:0x03"`
// }

type TplResultPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:1,equal:0x52"` // R
	Data    []byte `byte:"len:*"`            // room,command,state
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

// type TplSalePacket struct {
// 	STX_           []byte `byte:"len:1,equal:0x02"`
// 	Record_        []byte `byte:"len:1,equal:0x56"` // V
// 	SequenceNumber []byte `byte:"len:5"`            // ddddd or bbbbb
// 	Room           []byte `byte:"len:*"`            //
// 	Index          []byte `byte:"len:2"`            //
// 	Price          []byte `byte:"len:7"`            //
// 	Description    []byte `byte:"len:20"`           //
// 	Date           []byte `byte:"len:10"`           // YYMMDDhhmm
// 	ETX_           []byte `byte:"len:1,equal:0x03"`
// }

type TplSalePacket struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Record_        []byte `byte:"len:1,equal:0x56"` // V
	SequenceNumber []byte `byte:"len:5"`            // ddddd or bbbbb
	Data           []byte `byte:"len:*"`            // room,index,price,description,date
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

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
