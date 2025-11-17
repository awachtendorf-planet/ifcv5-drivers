package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"

	PacketCodeCard   = "Code Card"
	PacketCodeCard14 = "Code Card 14"
	PacketCodeCard15 = "Code Card 15"
	PacketCheckout   = "Checkout"

	PacketCodeCardAnswer = "Code Card Answer"
	PacketError          = "Error"
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
type TplCodeCardPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	Cmd_          []byte `byte:"len:1,equal:0x43"` // C
	CardType      []byte `byte:"len:1"`            // N C A
	CardCount     []byte `byte:"len:*"`            // Number of Cards Nothing or 1-9
	Sep00_        []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // Number of Encoder 1-9
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Ejection_     []byte `byte:"len:1,equal:0x45"` // (E)jection (R)etenetion or T: Ejection by the rear side
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room          []byte `byte:"len:*"`            // Room max 7 characters
	Sep2_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room2         []byte `byte:"len:*"`            // Room max 7 characters
	Sep3_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room3         []byte `byte:"len:*"`            // Room max 7 characters
	Sep4_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room4         []byte `byte:"len:*"`            // Room max 7 characters
	Sep5_         []byte `byte:"len:1,equal:0xb3"` // Separator
	AA            []byte `byte:"len:*"`            // Assigned Authorisations
	Sep6_         []byte `byte:"len:1,equal:0xb3"` // Separator
	AD            []byte `byte:"len:*"`            // Denied Authorisations
	Sep7_         []byte `byte:"len:1,equal:0xb3"` // Separator
	InitalDate    []byte `byte:"len:*"`            // HHDDMMYY
	Sep8_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ExpireDate    []byte `byte:"len:*"`            // HHDDMMYY
	Sep9_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Operator      []byte `byte:"len:*"`            // Operator Data
	Sep10_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Track1        []byte `byte:"len:*"`            // ISO standard accepts 65 characters on track 1, 62 alphanumeric characters + 3 control characters
	Sep11_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Track2        []byte `byte:"len:*"`            // ISO standard accepts 17 characters on track 2, 14 numeric characters + 3 control characters
	Sep12_        []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplCodeCardPacket14 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	Cmd_          []byte `byte:"len:1,equal:0x43"` // C
	CardType      []byte `byte:"len:1"`            // N C A
	CardCount     []byte `byte:"len:*"`            // Number of Cards Nothing or 1-9
	Sep00_        []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // Number of Encoder 1-9
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Ejection_     []byte `byte:"len:1,equal:0x45"` // (E)jection (R)etenetion or T: Ejection by the rear side
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room          []byte `byte:"len:*"`            // Room max 7 characters
	Sep2_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room2         []byte `byte:"len:*"`            // Room max 7 characters
	Sep3_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room3         []byte `byte:"len:*"`            // Room max 7 characters
	Sep4_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room4         []byte `byte:"len:*"`            // Room max 7 characters
	Sep5_         []byte `byte:"len:1,equal:0xb3"` // Separator
	AA            []byte `byte:"len:*"`            // Assigned Authorisations
	Sep6_         []byte `byte:"len:1,equal:0xb3"` // Separator
	AD            []byte `byte:"len:*"`            // Denied Authorisations
	Sep7_         []byte `byte:"len:1,equal:0xb3"` // Separator
	InitalDate    []byte `byte:"len:*"`            // HHDDMMYY
	Sep8_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ExpireDate    []byte `byte:"len:*"`            // HHDDMMYY
	Sep9_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Operator      []byte `byte:"len:*"`            // Operator Data
	Sep10_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Track1        []byte `byte:"len:*"`            // ISO standard accepts 65 characters on track 1, 62 alphanumeric characters + 3 control characters
	Sep11_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Track2        []byte `byte:"len:*"`            // ISO standard accepts 17 characters on track 2, 14 numeric characters + 3 control characters
	Sep12_        []byte `byte:"len:1,equal:0xb3"` // Separator
	KeyID         []byte `byte:"len:*"`            // KeyID
	Sep13_        []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplCodeCardPacket15 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	Cmd_          []byte `byte:"len:1,equal:0x43"` // C
	CardType      []byte `byte:"len:1"`            // N C A
	CardCount     []byte `byte:"len:*"`            // Number of Cards Nothing or 1-9
	Sep00_        []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // Number of Encoder 1-9
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Ejection_     []byte `byte:"len:1,equal:0x45"` // (E)jection (R)etenetion or T: Ejection by the rear side
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room          []byte `byte:"len:*"`            // Room max 7 characters
	Sep2_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room2         []byte `byte:"len:*"`            // Room max 7 characters
	Sep3_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room3         []byte `byte:"len:*"`            // Room max 7 characters
	Sep4_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Room4         []byte `byte:"len:*"`            // Room max 7 characters
	Sep5_         []byte `byte:"len:1,equal:0xb3"` // Separator
	AA            []byte `byte:"len:*"`            // Assigned Authorisations
	Sep6_         []byte `byte:"len:1,equal:0xb3"` // Separator
	AD            []byte `byte:"len:*"`            // Denied Authorisations
	Sep7_         []byte `byte:"len:1,equal:0xb3"` // Separator
	InitalDate    []byte `byte:"len:*"`            // HHDDMMYY
	Sep8_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ExpireDate    []byte `byte:"len:*"`            // HHDDMMYY
	Sep9_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Operator      []byte `byte:"len:*"`            // Operator Data
	Sep10_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Track1        []byte `byte:"len:*"`            // ISO standard accepts 65 characters on track 1, 62 alphanumeric characters + 3 control characters
	Sep11_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Track2        []byte `byte:"len:*"`            // ISO standard accepts 17 characters on track 2, 14 numeric characters + 3 control characters
	Sep12_        []byte `byte:"len:1,equal:0xb3"` // Separator
	KeyID         []byte `byte:"len:*"`            // KeyID
	Sep13_        []byte `byte:"len:1,equal:0xb3"` // Separator
	Extra_        []byte `byte:"len:1,equal:0x31"` // Extra
	Sep14_        []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplCheckoutPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"`   // Separator
	Cmd_          []byte `byte:"len:2,equal:0x434f"` // CO
	Sep0_         []byte `byte:"len:1,equal:0xb3"`   // Separator
	EncoderNumber []byte `byte:"len:*"`              // EncoderNumber
	Sep1_         []byte `byte:"len:1,equal:0xb3"`   // Separator
	Room          []byte `byte:"len:*"`              // Room max 7 characters
	Sep2_         []byte `byte:"len:1,equal:0xb3"`   // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

//Incoming
type TplAnswerPositivePacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	Cmd_          []byte `byte:"len:1,equal:0x43"` // C
	CardType      []byte `byte:"len:1"`            // N C A O
	CardCount     []byte `byte:"len:*"`            // Number of Cards Nothing or 1-9
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // EncoderNumber
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplAnswerPositiveAdvancedPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	Cmd_          []byte `byte:"len:1,equal:0x43"` // C
	CardType      []byte `byte:"len:1"`            // N C A O
	CardCount     []byte `byte:"len:*"`            // Number of Cards Nothing or 1-9
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // EncoderNumber
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	Data          []byte `byte:"len:*"`            // Additional Data
	Sep2_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplAnswerSimpleErrorPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	ErrorCode     []byte `byte:"len:2"`            //
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // EncoderNumber
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplAnswerAdvancedErrorPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Sep_          []byte `byte:"len:1,equal:0xb3"` // Separator
	Ovr_          []byte `byte:"len:1,equal:0x5f"` // underscore
	ErrorCode     []byte `byte:"len:4"`            // Exxx
	Sep0_         []byte `byte:"len:1,equal:0xb3"` // Separator
	EncoderNumber []byte `byte:"len:*"`            // EncoderNumber
	Sep1_         []byte `byte:"len:1,equal:0xb3"` // Separator
	ETX_          []byte `byte:"len:1,equal:0x03"`
}
