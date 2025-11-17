package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"

	PacketSerialFrameIn = "Serial Frame in"
	PacketTCPFrameIn    = "TCP Frame in"

	PacketTCPLogin     = "Login Message"
	PacketLoginAnswer  = "Login Answer"
	PacketTCPGarbage   = "Tcp Garbage"
	PacketLoginGarbage = "Login Garbage"

	PacketAdminMessage = "Admin Message"

	PacketTG41 = "Telegram 41"
	PacketTG60 = "Telegram 60"
	PacketTG67 = "Telegram 67"
	PacketTG70 = "Telegram 70"
	PacketTG71 = "Telegram 71"
	PacketTG80 = "Telegram 80"
	PacketTG00 = "Telegram 00"

	PacketTG10 = "Telegram 10"
	PacketTG20 = "Telegram 20"
	PacketTG40 = "Telegram 40"
	PacketTG72 = "Telegram 72"
	// In and Outgoing templates must have the same name! automate-> main

	PacketAnswerTG41 = "Telegram 41"
	PacketAnswerTG67 = "Telegram 67"
	PacketAnswerTG71 = "Telegram 71"
	PacketAnswerTG70 = "Telegram 70"
	PacketAnswerTG60 = "Telegram 60"
	PacketAnswerTG80 = "Telegram 80"
	PacketAnswerTG00 = "Telegram 00"

	PacketTG10Answer = "Telegram 10 Answer"
	PacketTG20Answer = "Telegram 20 Answer"
	PacketTG40Answer = "Telegram 40 Answer"
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

type TplTCPGarbagePacket struct {
	Data    []byte `byte:"len:*"`
	Header_ []byte `byte:"len:9,equal:0x446554655765434950"`
}
type TplTCPLoginGarbagePacket struct {
	Data    []byte `byte:"len:*"`
	Header_ []byte `byte:"len:13,equal:0x446554655765434941646d696e"`
}

// Outgoing

type TplTCPLoginPacket struct {
	Login_ []byte `byte:"len:26,equal:0x446554655765434941646d696e3a4c6f67696e3a333338333933"`
	Config []byte `byte:"len:*"`
}

type TplTCPAdminPacket struct {
	Header    []byte `byte:"len:14,equal:0x446554655765434941646d696e3a"`
	Telegrams []byte `byte:"len:*"`
}

// Incoming

type TplSerialFrameInPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	StagedPayload []byte `byte:"len:*"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplTCPAnswerFramePacket struct {
	Len           uint16 `byte:"len:2,endian:big"`
	StagedPayload []byte `byte:"len:{{.Len}}"`
}

type TplTCPLoginPositiveAnswerPacket struct {
	Login_ []byte `byte:"len:25,equal:0x446554655765434941646d696e3a43492d5374617465204f4b"`
}

type TplTCPLoginPositiveAnswerLenPacket struct {
	Login_ []byte `byte:"len:27,equal:0x0019446554655765434941646d696e3a43492d5374617465204f4b"`
}

//Telegrams

// Outgoing

type TplTestMessagePacket struct {
	CMD_    []byte `byte:"len:2,equal:0x3030"`
	Blank_  []byte `byte:"len:5,equal:0x2020202020"`
	Status_ []byte `byte:"len:1,equal:0x30"` // 0 or 9
}

type TplTg60Packet struct {
	CMD_        []byte `byte:"len:2,equal:0x3630"`
	Participant []byte `byte:"len:4"`
	Right       []byte `byte:"len:1"` // 0-8
	// OrderType   []byte `byte:"len:1"` // 0,1
}

type TplTg41Packet struct {
	CMD_        []byte `byte:"len:2,equal:0x3431"`
	Participant []byte `byte:"len:4"`
	DisplayType []byte `byte:"len:1"` // 0-9
	Text        []byte `byte:"len:16"`
	// Additions   []byte `byte:"len:*"` // Depending on Display-Type (just 2+)
}

type TplTg67Packet struct {
	CMD_        []byte `byte:"len:2,equal:0x3637"`
	Participant []byte `byte:"len:4"`
	Target_     []byte `byte:"len:4,equal:0x30303030"`
	Action_     []byte `byte:"len:1,equal:0x30"` // 0,1,3,5,8
	// FWDType        []byte `byte:"len:1"` // Forwarding Type 1,2,3,4
	// TargetAddition []byte `byte:"len:*"` // Langrufnummer Byte 5-24
}

type TplTg80Packet struct {
	CMD_        []byte `byte:"len:2,equal:0x3830"`
	Participant []byte `byte:"len:4"`
	Action      []byte `byte:"len:1"` // 0,1
}

type TplTg70Packet struct {
	CMD_        []byte `byte:"len:2,equal:0x3730"`
	Participant []byte `byte:"len:4"`
}

type TplTg71Packet struct {
	CMD_        []byte `byte:"len:2,equal:0x3731"`
	ControlCode []byte `byte:"len:1"`
	Participant []byte `byte:"len:4"`
	Year        []byte `byte:"len:4"`
	Month       []byte `byte:"len:2"`
	Day         []byte `byte:"len:2"`
	Hour        []byte `byte:"len:2"`
	Minute      []byte `byte:"len:2"`
}

// outgoing request telegrams

type TplTg10AnswerPacket struct {
	CMD_               []byte `byte:"len:2,equal:0x3130"`
	Blank_             []byte `byte:"len:3,equal:0x202020"`
	NotificationNumber []byte `byte:"len:2"`
	Result             []byte `byte:"len:1"`
}

type TplTg20AnswerPacket struct {
	CMD_               []byte `byte:"len:2,equal:0x3230"`
	Blanks0_           []byte `byte:"len:3,equal:0x202020"`
	NotificationNumber []byte `byte:"len:2"`
	Result             []byte `byte:"len:1"`
}

type TplTg40AnswerPacket struct {
	CMD_               []byte `byte:"len:2,equal:0x3430"`
	Blanks0_           []byte `byte:"len:3,equal:0x202020"`
	NotificationNumber []byte `byte:"len:2"`
	Result             []byte `byte:"len:1"`
}

// Incoming

type TplTestMessageAnswerPacket struct {
	CMD_    []byte `byte:"len:2,equal:0x3030"`
	Blank_  []byte `byte:"len:5,equal:0x2020202020"`
	Status_ []byte `byte:"len:1,equal:0x30"` // 0 or 9
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplTg60AnswerPacket struct {
	CMD         []byte `byte:"len:2,equal:0x3630"`
	Participant []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplTg70AnswerPacket struct {
	CMD         []byte `byte:"len:2,equal:0x3730"`
	Participant []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplTg80AnswerPacket struct {
	CMD         []byte `byte:"len:2,equal:0x3830"`
	Participant []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplTg41AnswerPacket struct {
	CMD         []byte `byte:"len:2,equal:0x3431"`
	Participant []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"` // 0,1,3
	Load        []byte `byte:"len:*"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplTg67AnswerPacket struct {
	CMD         []byte `byte:"len:2,equal:0x3637"`
	Participant []byte `byte:"len:4"`
	Target      []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"` // 1,2,3
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplTg71AnswerPacket struct {
	CMD         []byte `byte:"len:2,equal:0x3731"`
	Participant []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"`
	Year        []byte `byte:"len:4"`
	Month       []byte `byte:"len:2"`
	Day         []byte `byte:"len:2"`
	Hour        []byte `byte:"len:2"`
	Minute      []byte `byte:"len:2"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

// incoming request telegrams

type TplTg10Packet struct {
	CMD                []byte `byte:"len:2,equal:0x3130"`
	Blanks0_           []byte `byte:"len:3,equal:0x202020"`
	NotificationNumber []byte `byte:"len:2"`
	Blanks1_           []byte `byte:"len:1,equal:0x20"`
	Participant        []byte `byte:"len:4"`
	Blanks2_           []byte `byte:"len:1,equal:0x20"`
	SortSign           []byte `byte:"len:1"`
	Blanks3_           []byte `byte:"len:1,equal:0x20"`
	Date               []byte `byte:"len:6"` // yymmdd
	Blanks4_           []byte `byte:"len:1,equal:0x20"`
	Time               []byte `byte:"len:5"` // hh:mm
	Blanks5_           []byte `byte:"len:1,equal:0x20"`
	Duration           []byte `byte:"len:9"` // hhHmmMssS
	Blanks6_           []byte `byte:"len:1,equal:0x20"`
	SequenceNumber     []byte `byte:"len:4"`
	Blanks7_           []byte `byte:"len:1,equal:0x20"`
	TrunkLineNumber    []byte `byte:"len:3"`
	Blanks8_           []byte `byte:"len:1,equal:0x20"`
	PhoneNumber        []byte `byte:"len:20"`
	Blanks9_           []byte `byte:"len:1,equal:0x20"`
	Taxe               []byte `byte:"len:7"`
	Blanks10_          []byte `byte:"len:1,equal:0x20"`
	ID                 []byte `byte:"len:9"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type TplTg20Packet struct {
	CMD                []byte `byte:"len:2,equal:0x3230"`
	Blanks0_           []byte `byte:"len:3,equal:0x202020"`
	NotificationNumber []byte `byte:"len:2"`
	Blanks1_           []byte `byte:"len:1,equal:0x20"`
	Participant        []byte `byte:"len:4"`
	Blanks2_           []byte `byte:"len:1,equal:0x20"`
	RoomState          []byte `byte:"len:1"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type TplTg40Packet struct {
	CMD                []byte `byte:"len:2,equal:0x3430"`
	Blanks0_           []byte `byte:"len:3,equal:0x202020"`
	NotificationNumber []byte `byte:"len:2"`
	Blanks1_           []byte `byte:"len:1,equal:0x20"`
	Participant        []byte `byte:"len:4"`
	Blanks2_           []byte `byte:"len:1,equal:0x20"`
	Date               []byte `byte:"len:6"` // yymmdd
	Time               []byte `byte:"len:5"` // hh:mm
	Blanks3_           []byte `byte:"len:1,equal:0x20"`
	Data               []byte `byte:"len:20"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type TplTg72Packet struct {
	CMD         []byte `byte:"len:2,equal:0x3732"`
	ControlCode []byte `byte:"len:1"`
	Participant []byte `byte:"len:4"`
	Year        []byte `byte:"len:4"`
	Month       []byte `byte:"len:2"`
	Day         []byte `byte:"len:2"`
	Hour        []byte `byte:"len:2"`
	Minute      []byte `byte:"len:2"`
	Customer    []byte `byte:"len:4"`
	Result      []byte `byte:"len:1"`
}
