package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketAlive          = "Alive"
	PacketBeaconRequest  = "Optional beacon request"
	PacketBeaconResponse = "Optional beacon response"
	PacketKeyRequest     = "Key request"
	PacketKeyDelete      = "Key delete"
	PacketSuccess        = "Transaction success"
	PacketError          = "Transaction aborted"

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

// alive
type TplPFC_10 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3130"`   // PFC 10
	To_       []byte `byte:"len:2,equal:0x3030"`   // to interface station number
	From_     []byte `byte:"len:2,equal:0x3030"`   // from interface station number
	Terminal_ []byte `byte:"len:3,equal:0x303030"` // pms terminal 000
	Option_   []byte `byte:"len:2,equal:0x4645"`   // wants to receive SRC 55 linked response
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// PFC20 response
type TplPFC_62_000_00 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3632"`   // PFC 62
	To_       []byte `byte:"len:2,equal:0x3030"`   // to interface station number
	From_     []byte `byte:"len:2,equal:0x3030"`   // from interface station number
	Terminal_ []byte `byte:"len:3,equal:0x303030"` // pms terminal 000
	Status_   []byte `byte:"len:2,equal:0x3030"`   // transaction completed successfully
	Option_   []byte `byte:"len:1,equal:0x30"`     // future use
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// checkin / key request
type TplPFC_20_00X struct {
	STX_             []byte `byte:"len:1,equal:0x02"`
	Cmd_             []byte `byte:"len:2,equal:0x3230"` // PFC 20
	To_              []byte `byte:"len:2,equal:0x3030"` // to interface station number
	From_            []byte `byte:"len:2,equal:0x3030"` // from interface station number
	Terminal         []byte `byte:"len:3"`              // pms terminal
	TXC              []byte `byte:"len:3"`              // TXC 001 = new key, 003 = duplicate key
	Password         []byte `byte:"len:7"`              // Saflok Password
	KeyNumber        []byte `byte:"len:*"`              // Room
	KeyLevel         []byte `byte:"len:1"`              //
	EncoderStation   []byte `byte:"len:2"`              // Encoder Number
	EncoderLED_      []byte `byte:"len:2,equal:0x4646"` //
	KeyCount         []byte `byte:"len:2"`              //
	CheckoutDate     []byte `byte:"len:6"`              // MMDDYY
	CheckoutTime     []byte `byte:"len:4"`              // HHMM
	KeyExpireDate    []byte `byte:"len:6"`              // MMDDYY
	KeyExpireTime    []byte `byte:"len:4"`              // HHMM
	PassNumberOption []byte `byte:"len:1"`              //
	PassNumber       []byte `byte:"len:12"`             // Accesspoints
	TrackData        []byte `byte:"len:*"`              //
	ETX_             []byte `byte:"len:1,equal:0x03"`
}

// checkout / key delete
type TplPFC_20_015 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3230"`   // PFC 20
	To_       []byte `byte:"len:2,equal:0x3030"`   // to interface station number
	From_     []byte `byte:"len:2,equal:0x3030"`   // from interface station number
	Terminal  []byte `byte:"len:3"`                // pms terminal
	TXC_      []byte `byte:"len:3,equal:0x303135"` // TXC 015 -> check out
	Password  []byte `byte:"len:7"`                // Saflok Password
	KeyNumber []byte `byte:"len:*"`                // Room
	KeyLevel  []byte `byte:"len:1"`                //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// incoming

// response to PFC 10
type TplSRC_55 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3535"`   // SRC 55
	To_       []byte `byte:"len:2,equal:0x3030"`   // to interface station number
	From_     []byte `byte:"len:2,equal:0x3030"`   // from interface station number
	Terminal_ []byte `byte:"len:3,equal:0x303030"` // pms terminal 000
	Option_   []byte `byte:"len:2,equal:0x4645"`   // wants to receive SRC 55 linked response
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// TCP optional beacon 135
type TplSRC_20_135 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3230"`   // SRC 20
	To_       []byte `byte:"len:2,equal:0x3030"`   // to interface station number
	From_     []byte `byte:"len:2,equal:0x3030"`   // from interface station number
	Terminal_ []byte `byte:"len:3,equal:0x303030"` // pms terminal 000
	TXC_      []byte `byte:"len:3,equal:0x313335"` // TXC 135 -> TCP optional beacon message (check connection)
	Data_     []byte `byte:"len:*"`                // overread
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// transaction success
type TplSRC_62_00 struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	Cmd_     []byte `byte:"len:2,equal:0x3632"` // SRC 62
	To_      []byte `byte:"len:2,equal:0x3030"` // to interface station number
	From_    []byte `byte:"len:2,equal:0x3030"` // from interface station number
	Terminal []byte `byte:"len:3"`              // pms terminal
	Status   []byte `byte:"len:2,equal:0x3030"` // transaction completed successfully
	KeyCount []byte `byte:"len:3"`              // number of keys made
	Options  []byte `byte:"len:*"`              // futur use
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

// transaction error
type TplSRC_62_03 struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Cmd_         []byte `byte:"len:2,equal:0x3632"` // SRC 62
	To_          []byte `byte:"len:2,equal:0x3030"` // to interface station number
	From_        []byte `byte:"len:2,equal:0x3030"` // from interface station number
	Terminal     []byte `byte:"len:3"`              // pms terminal
	Status       []byte `byte:"len:2,equal:0x3033"` // non-recoverable procedural error
	ResponseCode []byte `byte:"len:3"`              // number of keys made
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

// transaction error
type TplSRC_62_04 struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Cmd_         []byte `byte:"len:2,equal:0x3632"` // SRC 62
	To_          []byte `byte:"len:2,equal:0x3030"` // to interface station number
	From_        []byte `byte:"len:2,equal:0x3030"` // from interface station number
	Terminal     []byte `byte:"len:3"`              // pms terminal
	Status       []byte `byte:"len:2,equal:0x3034"` // non-recoverable system error
	ResponseCode []byte `byte:"len:3"`              // number of keys made
	Options      []byte `byte:"len:*"`              // futur use
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

// transaction error
type TplSRC_62_05 struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Cmd_         []byte `byte:"len:2,equal:0x3632"` // SRC 62
	To_          []byte `byte:"len:2,equal:0x3030"` // to interface station number
	From_        []byte `byte:"len:2,equal:0x3030"` // from interface station number
	Terminal     []byte `byte:"len:3"`              // pms terminal
	Status       []byte `byte:"len:2,equal:0x3035"` // non-recoverable system error
	ResponseCode []byte `byte:"len:3"`              // number of keys made
	Options      []byte `byte:"len:*"`              // futur use
	ETX_         []byte `byte:"len:1,equal:0x03"`
}
