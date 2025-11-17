package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketKeyRequest = "Create Key"
	PacketKeyDelete  = "Delete Key"

	PacketCodeCardAnswer = "Code Card Answer"
	PacketError          = "Code Card Error"

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

type TplCodeCard struct {
	STX_             []byte `byte:"len:1,equal:0x02"`
	Sep1_            []byte `byte:"len:1,equal:0xB3"`
	Encoder          []byte `byte:"len:*"`            // encoder number
	Sep2_            []byte `byte:"len:1,equal:0xB3"` //
	C_               []byte `byte:"len:1,equal:0x43"` // C
	Action           []byte `byte:"len:1"`            // I/G
	KeyCount         []byte `byte:"len:*"`            // key cards count, optional
	Sep3_            []byte `byte:"len:1,equal:0xB3"` //
	Room             []byte `byte:"len:*"`            // max 18 characters
	Sep4_            []byte `byte:"len:1,equal:0xB3"` //
	ActivationDate   []byte `byte:"len:*"`            // DD/MM/YYYY, optional
	Sep5_            []byte `byte:"len:1,equal:0xB3"` //
	ActivationTime   []byte `byte:"len:*"`            // HH:MM, optional
	Sep6_            []byte `byte:"len:1,equal:0xB3"` //
	ExpirationDate   []byte `byte:"len:*"`            // DD/MM/YYYY, optional
	Sep7_            []byte `byte:"len:1,equal:0xB3"` //
	ExpirationTime   []byte `byte:"len:*"`            // HH:MM, optional
	Sep8_            []byte `byte:"len:1,equal:0xB3"` //
	Grants           []byte `byte:"len:*"`            // Pool,Library ... optional
	Sep9_            []byte `byte:"len:1,equal:0xB3"` //
	KeyPad           []byte `byte:"len:*"`            // pin code 4..6 numeric, optional
	Sep10_           []byte `byte:"len:1,equal:0xB3"` //
	CardOperation    []byte `byte:"len:*"`            // EF/RP/ED, optional, default EF
	Sep11_           []byte `byte:"len:1,equal:0xB3"` //
	Operator         []byte `byte:"len:*"`            // max 10 characters, optional
	Sep12_           []byte `byte:"len:1,equal:0xB3"` //
	TesaHotelEncoder []byte `byte:"len:*"`            // 0..6, optional ?
	Sep13_           []byte `byte:"len:1,equal:0xB3"` //
	Track1           []byte `byte:"len:*"`            //
	Sep14_           []byte `byte:"len:1,equal:0xB3"` //
	Track2           []byte `byte:"len:*"`            //
	Sep15_           []byte `byte:"len:1,equal:0xB3"` //
	Technology       []byte `byte:"len:*"`            // M/C/K/P ?
	Sep16_           []byte `byte:"len:1,equal:0xB3"` //
	Room2            []byte `byte:"len:*"`            // additional room
	Sep17_           []byte `byte:"len:1,equal:0xB3"` //
	Room3            []byte `byte:"len:*"`            // additional room
	Sep18_           []byte `byte:"len:1,equal:0xB3"` //
	Room4            []byte `byte:"len:*"`            // additional room
	Sep19_           []byte `byte:"len:1,equal:0xB3"` //
	ReturnCardID     []byte `byte:"len:*"`            // set to 1 answer will include card id
	Sep20_           []byte `byte:"len:1,equal:0xB3"` //
	CardID           []byte `byte:"len:*"`            // Tesa Hotel Encoder must be seet to 0 or 2
	Sep21_           []byte `byte:"len:1,equal:0xB3"` //
	CardType         []byte `byte:"len:*"`            // 1/2/3, optional
	Sep22_           []byte `byte:"len:1,equal:0xB3"` //
	PhoneNumber      []byte `byte:"len:*"`            // optional
	Sep23_           []byte `byte:"len:1,equal:0xB3"` //
	Mail             []byte `byte:"len:*"`            // optional
	Sep24_           []byte `byte:"len:1,equal:0xB3"` //
	Mail2            []byte `byte:"len:*"`            // optional
	Sep25_           []byte `byte:"len:1,equal:0xB3"` //
	Mail3            []byte `byte:"len:*"`            // optional
	Sep26_           []byte `byte:"len:1,equal:0xB3"` //
	Mail4            []byte `byte:"len:*"`            // optional
	Sep27_           []byte `byte:"len:1,equal:0xB3"` //
	ETX_             []byte `byte:"len:1,equal:0x03"`
}

type TplDeleteCard struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	Sep1_    []byte `byte:"len:1,equal:0xB3"`
	Encoder  []byte `byte:"len:*"`
	Sep2_    []byte `byte:"len:1,equal:0xB3"`
	C_       []byte `byte:"len:1,equal:0x43"` // C
	Action_  []byte `byte:"len:1,equal:0x4f"` // O
	Sep3_    []byte `byte:"len:1,equal:0xB3"`
	Room     []byte `byte:"len:*"`            // max 18 characters
	Sep4_    []byte `byte:"len:1,equal:0xB3"` //
	Operator []byte `byte:"len:*"`            // max 10 characters, optional
	Sep5_    []byte `byte:"len:1,equal:0xB3"`
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

// incoming

type TplCodeAnswer struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	Sep1_    []byte `byte:"len:1,equal:0xB3"`
	Encoder  []byte `byte:"len:*"`
	Sep2_    []byte `byte:"len:1,equal:0xB3"`
	C_       []byte `byte:"len:1,equal:0x43"` // C
	Action   []byte `byte:"len:1"`            // I/G/O
	KeyCount []byte `byte:"len:*"`            // key cards count, optional
	Sep3_    []byte `byte:"len:1,equal:0xB3"`
	CardID   []byte `byte:"len:*"`
	Sep4_    []byte `byte:"len:1,equal:0xB3"`
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

type TplCodeAnswerCO struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Sep1_   []byte `byte:"len:1,equal:0xB3"`
	Encoder []byte `byte:"len:*"`
	Sep2_   []byte `byte:"len:1,equal:0xB3"`
	C_      []byte `byte:"len:1,equal:0x43"` // C
	Action  []byte `byte:"len:1,equal:0x4f"` // O
	Sep3_   []byte `byte:"len:1,equal:0xB3"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplErrorPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Sep1_   []byte `byte:"len:1,equal:0xB3"`
	Encoder []byte `byte:"len:*"`
	Sep2_   []byte `byte:"len:1,equal:0xB3"`
	E_      []byte `byte:"len:1,equal:0x45"` // E
	Error   []byte `byte:"len:*"`
	Sep3_   []byte `byte:"len:1,equal:0xB3"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}
