package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketRequest1 = "request 1"
	PacketRequest2 = "request 2"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketFramed          = "Framed Packet"

	// must exactly match the generated struct type's
	PacketLinkStart       = "LinkStart"
	PacketLinkAlive       = "LinkAlive"
	PacketLinkDescription = "LinkDescription"
	PacketPostInquiry     = "PostInquiry"
	PacketPostRequest     = "PostRequest"
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

type TplGarbage_Framing_1 struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x02"`
}

type TplGarbage_Framing_2 struct {
	STX_OVR_ []byte `byte:"len:1,equal:0x02"`
	Data_    []byte `byte:"len:*"`
	STX_     []byte `byte:"len:1,equal:0x02"`
}

type TplFramedPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Payload []byte `byte:"len:*"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}
