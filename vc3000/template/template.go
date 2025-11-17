package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketRegister    = "Register"
	PacketUnregister  = "Unregister"
	PacketRegisterAck = "Register Ack"
	PacketRegisterNak = "Register Nak"

	PacketCheckout       = "Checkout Lcl"
	PacketCheckoutRmt    = "Checkout Rmt"
	PacketCodeCard       = "Code Card"
	PacketCodeCardModify = "Code Card Modify"
	PacketReadKey        = "Read Card"

	PacketAnswer     = "Answer"
	PacketAnswerData = "Answer Data"

	PacketGeneric         = "Generic Packet"
	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown Framed Packet"
)

// outgoing 62 bytes
type TplPacketRegister struct {
	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version int    `byte:"len:2,endian:little"`
	Cmd_    int    `byte:"len:1,equal:0x01"` // PMSifRegister
	Pad_    int    `byte:"len:3,equal:0x000000"`
	Len     int    `byte:"len:4,equal:0x2c000000,endian:little"`
	License []byte `byte:"len:19"`
	Pad1_   int    `byte:"len:1,equal:0x00"`
	AppName []byte `byte:"len:19"`
	Pad2_   int    `byte:"len:1,equal:0x00"`
	Ret_    int    `byte:"len:4,equal:0x00000000"`
}

// incoming
type TplPacketRegisterAck struct {
	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version int    `byte:"len:2,endian:little"`
	Cmd_    int    `byte:"len:1,equal:0x01"` // PMSifRegister
	Pad_    int    `byte:"len:3,equal:0x000000"`
	Len_    int    `byte:"len:4,endian:little"`
	License []byte `byte:"len:19"`
	Pad1_   int    `byte:"len:1,equal:0x00"`
	AppName []byte `byte:"len:19"`
	Pad2_   int    `byte:"len:1,equal:0x00"`
	Ret     int    `byte:"len:4,equal:0x00000000"`
}

// incoming
type TplPacketRegisterNak struct {
	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version int    `byte:"len:2,endian:little"`
	Cmd_    int    `byte:"len:1,equal:0x01"` // PMSifRegister
	Pad_    int    `byte:"len:3,equal:0x000000"`
	Len_    int    `byte:"len:4,endian:little"`
	License []byte `byte:"len:19"`
	Pad1_   int    `byte:"len:1,equal:0x00"`
	AppName []byte `byte:"len:19"`
	Pad2_   int    `byte:"len:1,equal:0x00"`
	Ret     int    `byte:"len:4,equal:0xFFFFFFFF"`
}

// incoming
type TplPacketUnregister struct {
	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version int    `byte:"len:2,endian:little"`
	Cmd_    int    `byte:"len:1,equal:0x02"` // PMSifUnregister
	Pad_    int    `byte:"len:3,equal:0x000000"`
	Len     int    `byte:"len:4,endian:little"`
	Data    []byte `byte:"len:{{.Len}}"`
}

// outgoing 583 bytes
type TplCodeCardSocket struct {
	Sync1_      []byte `byte:"len:4,equal:0x55555555"`
	Sync2_      []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version     int    `byte:"len:2,endian:little"`
	Cmd_        int    `byte:"len:1,equal:0x04"` // PMSifEncoderKcdRmt
	Pad_        int    `byte:"len:3,equal:0x000000"`
	Len         int    `byte:"len:4,endian:little"`
	FF          int    `byte:"len:1"`   // ff
	Data        []byte `byte:"len:511"` // dta
	Pad1_       int    `byte:"len:1,equal:0x00"`
	Destination int    `byte:"len:2,endian:little"` // DD
	Pad2_       int    `byte:"len:1,equal:0x00"`
	Source      int    `byte:"len:2,endian:little"` // SS
	Pad3_       int    `byte:"len:1,equal:0x00"`
	Debug_      int    `byte:"len:4,equal:0x00000000"`
	OpID        []byte `byte:"len:9"`
	Pad4_       int    `byte:"len:1,equal:0x00"`
	OpFirst     []byte `byte:"len:15"`
	Pad5_       int    `byte:"len:1,equal:0x00"`
	OpLast      []byte `byte:"len:15"`
	Pad6_       int    `byte:"len:1,equal:0x00"`
}

// outgoing 577 bytes
type TplCheckoutSocket struct {
	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version int    `byte:"len:2,endian:little"`
	Cmd_    int    `byte:"len:1,equal:0x03"` // PMSifEncoderKcdLcl
	Pad_    int    `byte:"len:3,equal:0x000000"`
	Len     int    `byte:"len:4,endian:little"`
	FF      int    `byte:"len:1"`   // ff
	Data    []byte `byte:"len:511"` // dta
	Pad1_   int    `byte:"len:1,equal:0x00"`
	Debug_  int    `byte:"len:4,equal:0x00000000"`
	OpID    []byte `byte:"len:9"`
	Pad4_   int    `byte:"len:1,equal:0x00"`
	OpFirst []byte `byte:"len:15"`
	Pad5_   int    `byte:"len:1,equal:0x00"`
	OpLast  []byte `byte:"len:15"`
	Pad6_   int    `byte:"len:1,equal:0x00"`
}

// incoming fallback
type TplGenericPacketSocket struct {
	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
	Version int    `byte:"len:2,endian:little"`
	Cmd     int    `byte:"len:4,endian:little"`
	Len     int    `byte:"len:4,endian:little"`
	Data    []byte `byte:"len:{{.Len}}"`
}

// incoming
type TplGarbageSocket struct {
	Data   []byte `byte:"len:*"`
	Sync1_ []byte `byte:"len:4,equal:0x55555555"`
	Sync2_ []byte `byte:"len:4,equal:0xAAAAAAAA"`
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

// outgoing
type TplCodeCardSerial struct {
	STX_        uint8  `byte:"len:1,equal:0x02"`
	Destination int    `byte:"len:2"` // dd
	Source      int    `byte:"len:2"` // ss
	FF          int    `byte:"len:1"` // ff
	Data        []byte `byte:"len:*"` // data
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

// outgoing
type TplReadCardSerial struct {
	STX_        uint8  `byte:"len:1,equal:0x02"`
	Destination int    `byte:"len:2"` // dd
	Source      int    `byte:"len:2"` // ss
	FF          int    `byte:"len:1"` // ff
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

// incoming
type TplAnswerSerial struct {
	STX_        uint8  `byte:"len:1,equal:0x02"`
	Destination int    `byte:"len:2"` // dd
	Source      int    `byte:"len:2"` // ss
	FF          int    `byte:"len:1"` // ff
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

// incoming
type TplAnswerWithDataSerial struct {
	STX_        uint8  `byte:"len:1,equal:0x02"`
	Destination int    `byte:"len:2"` // dd
	Source      int    `byte:"len:2"` // ss
	FF          int    `byte:"len:1"` // ff
	Data        []byte `byte:"len:*"` // data
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

// incoming
type TplUnknownSerial struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

// incoming
type TplGarbageSerial struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x02"`
}

// outgoing
// type TplPacketRegister struct {
// 	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
// 	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
// 	Version int    `byte:"len:2,endian:little"`
// 	Cmd_    int    `byte:"len:1,equal:0x01"` // PMSifRegister
// 	Pad_    int    `byte:"len:3,equal:0x000000"`
// 	Len     int    `byte:"len:4,endian:little"`
// 	Data    []byte `byte:"len:{{.Len}}"`
// }

// type TplCodeCardSocket struct {
// 	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
// 	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
// 	Version int    `byte:"len:2,endian:little"`
// 	Cmd_    int    `byte:"len:1,equal:0x04"` // PMSifEncoderKcdRmt
// 	Pad_    int    `byte:"len:3,equal:0x000000"`
// 	Len     int    `byte:"len:4,endian:little"`
// 	Data    []byte `byte:"len:{{.Len}}"`
// }

// type TplCheckoutSocket struct {
// 	Sync1_  []byte `byte:"len:4,equal:0x55555555"`
// 	Sync2_  []byte `byte:"len:4,equal:0xAAAAAAAA"`
// 	Version int    `byte:"len:2,endian:little"`
// 	Cmd_    int    `byte:"len:1,equal:0x03"` // PMSifEncoderKcdLcl
// 	Pad_    int    `byte:"len:3,equal:0x000000"`
// 	Len     int    `byte:"len:4,endian:little"`
// 	Data    []byte `byte:"len:{{.Len}}"`
// }
