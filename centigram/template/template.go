package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckIn                    = "Check In"
	PacketCheckOut                   = "Check Out"
	PacketRoomChange                 = "Room Change"
	PacketChangeName                 = "Change Display Name"
	PacketChangeFCOS                 = "Change FCOS"
	PacketChangeVoiceMail            = "Change Voice Mail"
	PacketMessageLampOn              = "Change Message Lamp On"
	PacketSwapEnquiry                = "Database Swap Enquiry"
	PacketMessageWaitingStatus       = "Message Waiting Status"
	PacketRefusion                   = "Refusion"
	PacketMessageWaitingStatusOnSwap = "Message Waiting Status (Swap)"

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

type TplHISCheckin struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x32"` // 2
	Extension []byte `byte:"len:6"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISCheckout struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x33"` // 3
	Extension []byte `byte:"len:6"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISMessageLampOn struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x36"` // 6
	Extension []byte `byte:"len:6"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISRoomChange struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x37"` // 7
	Source    []byte `byte:"len:6"`            // source
	Extension []byte `byte:"len:6"`            // destination
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISChangeFCOS struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x38"` // 8
	Extension []byte `byte:"len:6"`            //
	Code      []byte `byte:"len:2"`            // 00-64
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISChangeVoiceMail struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x39"` // 9
	Extension []byte `byte:"len:6"`            //
	Code      []byte `byte:"len:2"`            // 00-99
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISChangeName struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x44"` // D
	Extension []byte `byte:"len:6"`
	Name      []byte `byte:"len:6"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreCheckin struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x32"` // 2
	Extension []byte `byte:"len:6"`
	PAD_      []byte `byte:"len:6,equal:0x202020202020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreCheckout struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x33"` // 3
	Extension []byte `byte:"len:6"`            //
	PAD_      []byte `byte:"len:6,equal:0x202020202020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreMessageLampOn struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x36"` // 6
	Extension []byte `byte:"len:6"`            //
	PAD_      []byte `byte:"len:6,equal:0x202020202020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreRoomChange struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x37"` // 7
	Source    []byte `byte:"len:6"`            // source
	Extension []byte `byte:"len:6"`            // destination
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreChangeFCOS struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x38"` // 8
	Extension []byte `byte:"len:6"`            //
	Code      []byte `byte:"len:2"`            // 00-64
	PAD_      []byte `byte:"len:4,equal:0x20202020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreChangeVoiceMail struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x39"` // 9
	Extension []byte `byte:"len:6"`            //
	Code      []byte `byte:"len:2"`            // 00-99
	PAD_      []byte `byte:"len:4,equal:0x20202020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreChangeName struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x44"` // D
	Extension []byte `byte:"len:6"`
	Name      []byte `byte:"len:6"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// incomming

type TplHISSwapEnquiry struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:1,equal:0x31"` // 1
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplHISMessageWaitingStatus struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x34"` // 4
	Extension []byte `byte:"len:6"`            //
	Unread    []byte `byte:"len:2"`            //
	Urgent    []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplHISBadMailboxAddress struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x35"` // 5
	Extension []byte `byte:"len:6"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreSwapEnquiry struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:1,equal:0x31"` // 1
	OVR_ []byte `byte:"len:12,equal:0x202020202020202020202020"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreMessageWaitingStatus struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x34"` // 4
	Extension []byte `byte:"len:6"`            //
	Unread    []byte `byte:"len:2"`            //
	Urgent    []byte `byte:"len:2"`            //
	OVR_      []byte `byte:"len:2,equal:0x2020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEncoreBadMailboxAddress struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:1,equal:0x35"` // 5
	Extension []byte `byte:"len:6"`            //
	OVR_      []byte `byte:"len:6,equal:0x202020202020"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}
