package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckIn    = "Check In"
	PacketCheckOut   = "Check Out"
	PacketDataChange = "Data Change"

	PacketKeyCreate       = "Create Key"
	PacketKeyCreateResult = "Key Result"

	PacketKeyRead       = "Key Read"
	PacketKeyReadAnswer = "Key Read Answer"

	PacketStatusRequest       = "Status Request"
	PacketStatusRequestAnswer = "Status Request Answer"

	PacketError = "Error"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"
)

// outgoing from PMS:
/*
CES = request status terminal
DIC = Issue GC
DRC = reading GC
*/

// incoming from miwa:
/*
RES = notify status terminal
PIC = Issue GC completed
PRC = reading GC completed
RER = failure
*/

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

type TplStatusRequestPacket struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:3,equal:0x434553"` // CES
	ID   []byte `byte:"len:2"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCreateKeyPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Cmd_        []byte `byte:"len:3,equal:0x444943"` // DIC
	ID          []byte `byte:"len:2"`
	IssueType   []byte `byte:"len:1"` // (MIWA can only produce new cards)
	CardType    []byte `byte:"len:2"`
	CheckIn     []byte `byte:"len:10"`
	CheckOut    []byte `byte:"len:10"`
	MainRoom    []byte `byte:"len:8"`
	ExtraRoom1  []byte `byte:"len:8"`
	ExtraRoom2  []byte `byte:"len:8"`
	ExtraRoom3  []byte `byte:"len:8"`
	ExtraRoom4  []byte `byte:"len:8"`
	ExtraRoom5  []byte `byte:"len:8"`
	Reserve1    []byte `byte:"len:32"`
	SpecialFlag []byte `byte:"len:40"`
	IssueNumber []byte `byte:"len:2"`
	POSInfo     []byte `byte:"len:37"`
	StaffCode   []byte `byte:"len:6"`
	Reserve2    []byte `byte:"len:4"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplKeyReadPacket struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:3,equal:0x445243"` // DRC
	ID   []byte `byte:"len:2"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

// incoming

type TplStatusRequestAnswerPacket struct {
	STX_   []byte `byte:"len:1,equal:0x02"`
	Cmd_   []byte `byte:"len:3,equal:0x524553"` // RES
	ID_    []byte `byte:"len:2"`
	Type   []byte `byte:"len:1"` // 1 - DCR , 2- MCR2, CCU, USB R/W, 3 - to be determined
	Status []byte `byte:"len:1"` // 0 - ready, 1 - busy (initiated by terminal), 2 - busy(by PC), 3 - busy (by PMS), 4 - failure
	ETX_   []byte `byte:"len:1,equal:0x03"`
}

type TplCreateKeySuccessPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Cmd_        []byte `byte:"len:3,equal:0x504943"` // PIC
	ID          []byte `byte:"len:2"`
	Result      []byte `byte:"len:1"` // 0: Normally completed 1: Failure in writing 2: Card being stuck 3: Time out 9: Other failures
	IssueNumber []byte `byte:"len:2"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplKeyReadResultPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:3,equal:0x505243"` // PRC
	ID                 []byte `byte:"len:2"`
	Result             []byte `byte:"len:1"`
	KeyCardInformation []byte `byte:"len:144"`
	POSInfo            []byte `byte:"len:37"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type TplErrorPacket struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Cmd_  []byte `byte:"len:3,equal:0x524552"` // RER
	ID_   []byte `byte:"len:2"`
	Error []byte `byte:"len:1"` // See info table
	ETX_  []byte `byte:"len:1,equal:0x03"`
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
