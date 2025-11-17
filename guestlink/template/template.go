package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketVerify = "Command Acknowledge"
	PacketError  = "Command Refusion"

	PacketStart               = "Start Request"
	PacketTest                = "Test Line Request"
	PacketHelo                = "Helo"
	PacketInit                = "Init Request"
	PacketCheckIn             = "Check In"
	PacketCheckOut            = "Check Out"
	PacketWakeupSet           = "Wakeup Set"
	PacketWakeupClear         = "Wakeup Clear"
	PacketWakeupResult        = "Wakeup Result"
	PacketRoomStatus          = "Room Status"
	PacketGuestMessageRead    = "Guest Message Read"
	PacketGuestMessageStatus  = "Guest Message Status"
	PacketGuestMessageRequest = "Guest Message Request"
	PacketPostCharge          = "Post Charge"
	PacketCheckoutRequest     = "Express Checkout"
	PacketDisplayRequest      = "Display Request"
	PacketLookupRequest       = "Lookup Request"
	PacketStatusRequest       = "Status Request"
	PacketNameReply           = "Name Reply"
	PacketInfoReply           = "Info Reply"
	PacketItemReply           = "Item Reply"
	PacketBalanceReply        = "Balance Reply"
	PacketGuestMessageHeader  = "Guest Message Header"
	PacketGuestMessageCaller  = "Guest Message Caller"
	PacketGuestMessageText    = "Guest Message Text"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknownCommand  = "Unknown Command"
	PacketUnknown         = "Unknown Framed Packet"
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

type TplUnknownPacket struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplUnknownCommand struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4,equal:0x39393939"`
	Command_    []byte `byte:"len:4"`
	Data_       []byte `byte:"len:*"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
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
