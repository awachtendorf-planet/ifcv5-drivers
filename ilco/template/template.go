package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"
	PacketSof = "SOF"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"

	PacketCodeCard = "Code Card"

	PacketCodeCardAnswer = "Code Card Answer"
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

type TplSOF struct {
	SOF_ []byte `byte:"len:1,equal:0x7f"`
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
type TplCodeCardPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Addr    []byte `byte:"len:*"`            // Address could have length 0 or 2
	Cmd_    []byte `byte:"len:1,equal:0x31"` // Command '1' for Key Create
	Len     []byte `byte:"len:4"`            // Message length
	Payload []byte `byte:"len:*"`            // payload as: Key Value Sep_
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

//Incoming
type TplAnswerDirectPacket struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:1,equal:0x34"` // Answer '4' to Command '1'
	Len   []byte `byte:"len:4"`            // Message length
	Key   []byte `byte:"len:1"`            // Key
	Value []byte `byte:"len:*"`            // Value
	Sep_  []byte `byte:"len:1,equal:0x1c"` // Separator
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplAnswerGatewayPacket struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Addr  []byte `byte:"len:2"`            // Address encoder gateway
	Cmd_  []byte `byte:"len:1,equal:0x34"` // Answer '4' to Command '1'
	Len   []byte `byte:"len:4"`            // Message length
	Key   []byte `byte:"len:1"`            // Key
	Value []byte `byte:"len:*"`            // Value
	Sep_  []byte `byte:"len:1,equal:0x1c"` // Separator
	ETX_  []byte `byte:"len:1,equal:0x03"`
}
