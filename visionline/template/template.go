package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckout   = "Checkout"
	PacketCodeCard   = "Code Card"
	PacketCardUpdate = "Update Card"
	PacketReadKey    = "Read Card"
	PacketAlive      = "Alive"

	PacketGeneric         = "Generic Packet"
	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
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

// socket layer

type TplCodeCardSocket struct {
	Cmd_  []byte `byte:"len:3,equal:0x434341"` // CCA
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	EOT_  int    `byte:"len:2,equal:0x0d0a"`
}

type TplReadCardSocket struct {
	Cmd_  []byte `byte:"len:3,equal:0x434342"` // CCB
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	EOT_  int    `byte:"len:2,equal:0x0d0a"`
}

type TplAliveSocket struct {
	Cmd_  []byte `byte:"len:3,equal:0x434343"` // CCC
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	EOT_  int    `byte:"len:2,equal:0x0d0a"`
}

// type TplCardAutoUpdateSocket struct {
// 	Cmd_  []byte `byte:"len:3,equal:0x434344"` // CCD
// 	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
// 	Data  []byte `byte:"len:*"`
// 	EOT_  int    `byte:"len:2,equal:0x0d0a"`
// }

type TplCardUpdateSocket struct {
	Cmd_  []byte `byte:"len:3,equal:0x434346"` // CCF
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	EOT_  int    `byte:"len:2,equal:0x0d0a"`
}

type TplCheckoutSocket struct {
	Cmd_  []byte `byte:"len:3,equal:0x434347"` // CCG
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	EOT_  int    `byte:"len:2,equal:0x0d0a"`
}

type TplGenericPacketSocket struct {
	Cmd_  []byte `byte:"len:2,equal:0x4343"` // CC
	Cmd   int    `byte:"len:1"`              // ?
	Pipe_ int    `byte:"len:1,equal:0x3b"`   // ;
	Data  []byte `byte:"len:*"`
	EOT_  int    `byte:"len:2,equal:0x0d0a"`
}

type TplGarbageSocket struct {
	Garbage []byte `byte:"len:*"`
	Cmd_    []byte `byte:"len:2,equal:0x4343"` // CC
}

// serial layer

type TplCodeCardSerial struct {
	STX_  uint8  `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:3,equal:0x434341"` // CCA
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplReadCardSerial struct {
	STX_  uint8  `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:3,equal:0x434342"` // CCB
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplAliveSerial struct {
	STX_  uint8  `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:3,equal:0x434343"` // CCC
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

// type TplCardAutoUpdateSerial struct {
// 	STX_  uint8  `byte:"len:1,equal:0x02"`
// 	Cmd_  []byte `byte:"len:3,equal:0x434344"` // CCD
// 	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
// 	Data  []byte `byte:"len:*"`
// 	ETX_  []byte `byte:"len:1,equal:0x03"`
// }

type TplCardUpdateSerial struct {
	STX_  uint8  `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:3,equal:0x434346"` // CCF
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplCheckoutSerial struct {
	STX_  uint8  `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:3,equal:0x434347"` // CCG
	Pipe_ int    `byte:"len:1,equal:0x3b"`     // ;
	Data  []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplGenericPacketSerial struct {
	STX_  uint8  `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:2,equal:0x4343"` // CC
	Cmd   int    `byte:"len:1"`              // ?
	Pipe_ int    `byte:"len:1,equal:0x3b"`   // ;
	Data  []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplUnknownPacketSerial struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplGarbageSerial_1 struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x02"`
}

type TplGarbageSerial_2 struct {
	STX_OVR_ []byte `byte:"len:1,equal:0x02"`
	Data_    []byte `byte:"len:*"`
	STX_     []byte `byte:"len:1,equal:0x02"`
}
