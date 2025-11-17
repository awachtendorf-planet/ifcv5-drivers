package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketRequest       = "Request"
	PacketSwapEnquiry   = "Swap Request"
	PacketCheckIn       = "Check In"
	PacketCheckOut      = "Check Out"
	PacketDataChange    = "Data Change"
	PacketRoomMove      = "Room Move"
	PacketWakeupRequest = "Wakeup Request"
	PacketWakeupClear   = "Wakeup Clear"
	PacketRoomStatus    = "Room Status"

	PacketCheckInSwap       = "Check In (Swap)"
	PacketCheckOutSwap      = "Check Out (Swap)"
	PacketWakeupRequestSwap = "Wakeup Request (Swap)"
	PacketWakeupClearSwap   = "Wakeup Clear (Swap)"
	PacketRoomStatusSwap    = "Room Status (Swap)"

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

// incoming
type TplRequestPacket struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:1,equal:0x52"` // R
	Room []byte `byte:"len:*"`            // Room Number
	BS_  []byte `byte:"len:1,equal:0x5c"` // backslash
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplSwapRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:1,equal:0x46"` // F
	BS_  []byte `byte:"len:1,equal:0x5c"` // backslash
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

// outgoing
type TplOutgoing struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:1,equal:0x52"` // R
	Room []byte `byte:"len:*"`            // Room Number
	BS_  []byte `byte:"len:1,equal:0x5c"` // backslash
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplOutgoingMove struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Cmd_    []byte `byte:"len:1,equal:0x43"` // C
	OldRoom []byte `byte:"len:*"`            // Old Room Number
	UNDL1_  []byte `byte:"len:1,equal:0x5f"` //
	Room    []byte `byte:"len:*"`            // Destination Room Number
	UNDL2_  []byte `byte:"len:1,equal:0x5f"` //
	Data    []byte `byte:"len:*"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplOutgoingSwap struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	SwapStart_ []byte `byte:"len:2,equal:0x5a2b"` // Z+
	BS1_       []byte `byte:"len:1,equal:0x5c"`   // backslash
	Cmd_       []byte `byte:"len:1,equal:0x52"`   // R
	Room       []byte `byte:"len:*"`              // Room Number
	BS2_       []byte `byte:"len:1,equal:0x5c"`   // backslash
	Data       []byte `byte:"len:*"`
	SwapEnd_   []byte `byte:"len:2,equal:0x5a2d"` // Z-
	BS3_       []byte `byte:"len:1,equal:0x5c"`   // backslash
	ETX_       []byte `byte:"len:1,equal:0x03"`
}
