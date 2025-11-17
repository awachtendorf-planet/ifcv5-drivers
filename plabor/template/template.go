package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"
	PacketEot = "EOT"

	PacketCheckIn          = "Check In"
	PacketDataChangeUpdate = "Data Change Update"
	PacketDataChangeMove   = "Data Change Move"
	PacketCheckOut         = "Check Out"

	PacketChargePayTV   = "Charge PayTV"
	PacketChargeMinibar = "Charge Minibar"

	PacketReqBill  = "Request Bill"
	PacketBillPart = "Part of Bill"
	PacketBalance  = "Balance / End of Bill"

	PacketMessageSignal       = "Message Signal"
	PacketMessageBlock        = "Message Block"
	PacketMessageEnd          = "Message End"
	PacketMessageDelete       = "Delete Message"
	PacketMessageRequest      = "Message Request"
	PacketMessageConfirmation = "Message Confirmation"

	PacketWakeupIndirect = "Wakeup Indirect"
	PacketWakeupDirect   = "Wakeup Direct"

	PacketWakeupIndirectk = "Wakeup Indirect Short"
	PacketWakeupDirectk   = "Wakeup Direct Short"

	PacketError = "NAK Error Message"

	PacketDBSyncRequest    = "DBSync Request"
	PacketDBSyncFrameStart = "DBSync Frame Start"
	PacketDBSyncFrameStopp = "DBSync Frame Stopp"

	PacketRoomstatus = "Roomstatus"

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

type TplEOT struct {
	EOT_ []byte `byte:"len:1,equal:0x04"`
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
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	STX2_ []byte `byte:"len:1,equal:0x02"`
}

//Incoming Requests

type TplRecord53Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3533"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon_     []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Value      []byte `byte:"len:8"`
	Channel    []byte `byte:"len:2"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord54PLaborPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Record_     []byte `byte:"len:2,equal:0x3534"`
	RoomNumber  []byte `byte:"len:4"`
	AccountID   []byte `byte:"len:7"`
	Day         []byte `byte:"len:2"`
	Dot0_       []byte `byte:"len:1,equal:0x2e"`
	Month       []byte `byte:"len:2"`
	Dot1_       []byte `byte:"len:1,equal:0x2e"`
	Year        []byte `byte:"len:4"`
	RoomNumberD []byte `byte:"len:4"`
	AccountIDD  []byte `byte:"len:7"`
	Hour        []byte `byte:"len:2"`
	Colon0_     []byte `byte:"len:1,equal:0x3a"`
	Minute      []byte `byte:"len:2"`
	MessageId   []byte `byte:"len:6"`
}

type TplRecord54ProDacPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3534"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord62Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3632"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	Balance    []byte `byte:"len:9"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord63Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3633"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord71Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3731"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Operator   []byte `byte:"len:2"`
	Group      []byte `byte:"len:1"`
	Status     []byte `byte:"len:2"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord19Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3139"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	RoomNumber []byte `byte:"len:4"`
	Operator   []byte `byte:"len:4"`
	Status     []byte `byte:"len:2"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord73PLaborPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Record_       []byte `byte:"len:2,equal:0x3733"`
	Day           []byte `byte:"len:2"`
	Dot0_         []byte `byte:"len:1,equal:0x2e"`
	Month         []byte `byte:"len:2"`
	Dot1_         []byte `byte:"len:1,equal:0x2e"`
	Year          []byte `byte:"len:4"`
	RoomNumber    []byte `byte:"len:4"`
	Hour          []byte `byte:"len:2"`
	Colon0_       []byte `byte:"len:1,equal:0x3a"`
	Minute        []byte `byte:"len:2"`
	Operator      []byte `byte:"len:2"`
	Group         []byte `byte:"len:1"`
	Amount        []byte `byte:"len:1"`
	ArticleNumber []byte `byte:"len:2"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplRecord73ProDacPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3733"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	RoomNumber []byte `byte:"len:4"`
	TransCode  []byte `byte:"len:5"`
	TransText  []byte `byte:"len:32"`
	Value      []byte `byte:"len:9"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord91Packet struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:2,equal:0x3931"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplRecord58PLaborPacket struct { // not supported
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3638"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Colon1_    []byte `byte:"len:1,equal:0x3a"`
	Second     []byte `byte:"len:2"`
	Result     []byte `byte:"len:1"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord58ProDacPacket struct { // not supported
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3638"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Result     []byte `byte:"len:1"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

//Incoming Reponse

type TplRecord17PLaborPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3137"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Error      []byte `byte:"len:2"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord17ProDacPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3137"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Error      []byte `byte:"len:2"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord57Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3537"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Error      []byte `byte:"len:2"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord64ProLabPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3634"`
	AccountID  []byte `byte:"len:7"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Colon1_    []byte `byte:"len:1,equal:0x3a"`
	Second     []byte `byte:"len:2"`
	RoomNumber []byte `byte:"len:4"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord64ProDacPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3634"`
	AccountID  []byte `byte:"len:7"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	Colon1_    []byte `byte:"len:1,equal:0x3a"`
	Second     []byte `byte:"len:2"`
	RoomNumber []byte `byte:"len:4"`
	MessageId  []byte `byte:"len:6"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord75Packet struct { // not supported
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3735"`
	Day        []byte `byte:"len:2"`
	Dot0_      []byte `byte:"len:1,equal:0x2e"`
	Month      []byte `byte:"len:2"`
	Dot1_      []byte `byte:"len:1,equal:0x2e"`
	Year       []byte `byte:"len:4"`
	RoomNumber []byte `byte:"len:4"`
	Hour       []byte `byte:"len:2"`
	Colon0_    []byte `byte:"len:1,equal:0x3a"`
	Minute     []byte `byte:"len:2"`
	TVNumber   []byte `byte:"len:6"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

// Outgoing

type TplRecord21Packet struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Record_            []byte `byte:"len:2,equal:0x3231"`
	CIDate             []byte `byte:"len:10"`
	RoomNumber         []byte `byte:"len:4"`
	AccountID          []byte `byte:"len:7"`
	GuestName          []byte `byte:"len:40"`
	CODate             []byte `byte:"len:10"`
	RemoteCheckoutFlag []byte `byte:"len:1"`
	BillViewFlag       []byte `byte:"len:1"`
	TVProgrammFlag     []byte `byte:"len:1"`
	StandardVideoFlag  []byte `byte:"len:1"`
	AdultVideoFlag     []byte `byte:"len:1"`
	SeminarFlag        []byte `byte:"len:1"`
	SwitchOnTVFlag     []byte `byte:"len:1"`
	CodeNrHighFlag     []byte `byte:"len:1"`
	CodeNrLowFlag      []byte `byte:"len:1"`
	LanguageFlag       []byte `byte:"len:1"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type TplRecord31Packet struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Record_            []byte `byte:"len:2,equal:0x3331"`
	CIDate             []byte `byte:"len:10"`
	RoomNumber         []byte `byte:"len:4"`
	AccountID          []byte `byte:"len:7"`
	GuestName          []byte `byte:"len:40"`
	CODate             []byte `byte:"len:10"`
	RemoteCheckoutFlag []byte `byte:"len:1"`
	BillViewFlag       []byte `byte:"len:1"`
	TVProgrammFlag     []byte `byte:"len:1"`
	StandardVideoFlag  []byte `byte:"len:1"`
	AdultVideoFlag     []byte `byte:"len:1"`
	SeminarFlag        []byte `byte:"len:1"`
	SwitchOnTVFlag     []byte `byte:"len:1"`
	CodeNrHighFlag     []byte `byte:"len:1"`
	CodeNrLowFlag      []byte `byte:"len:1"`
	LanguageFlag       []byte `byte:"len:1"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type TplRecord22Packet struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Record_       []byte `byte:"len:2,equal:0x3232"`
	RoomNumberOld []byte `byte:"len:4"`
	AccountID     []byte `byte:"len:7"`
	RoomNumberNew []byte `byte:"len:4"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplRecord12Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3132"`
	CODate     []byte `byte:"len:10"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord23Packet struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Record_       []byte `byte:"len:2,equal:0x3233"`
	Date          []byte `byte:"len:10"`
	RoomNumber    []byte `byte:"len:4"`
	AccountID     []byte `byte:"len:7"`
	ArticleNumber []byte `byte:"len:2"`
	OrderText     []byte `byte:"len:32"`
	Value         []byte `byte:"len:9"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplRecord33Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3333"`
	Date       []byte `byte:"len:10"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	Total      []byte `byte:"len:9"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord14PLaborPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3134"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

// type TplRecord14ProDacPacket struct {
// 	STX_       []byte `byte:"len:1,equal:0x02"`
// 	Record_    []byte `byte:"len:2,equal:0x3134"`
// 	RoomNumber []byte `byte:"len:4"`
// 	AccountID  []byte `byte:"len:7"`
// 	Date       []byte `byte:"len:10"`
// 	Time       []byte `byte:"len:5"`
// 	MessageId  []byte `byte:"len:6"`
// 	ETX_       []byte `byte:"len:1,equal:0x03"`
// }

type TplRecord24Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3234"`
	Date       []byte `byte:"len:10"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	Time       []byte `byte:"len:8"`
	Text       []byte `byte:"len:40"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord34Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3334"`
	Date       []byte `byte:"len:10"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	Time       []byte `byte:"len:8"`
	Text       []byte `byte:"len:40"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord44Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3434"`
	RoomNumber []byte `byte:"len:4"`
	AccountID  []byte `byte:"len:7"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord18Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3138"`
	Date       []byte `byte:"len:10"`
	RoomNumber []byte `byte:"len:4"`
	Time       []byte `byte:"len:5"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}
type TplRecord18Packetk struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3138"`
	Date       []byte `byte:"len:5"`
	RoomNumber []byte `byte:"len:4"`
	Time       []byte `byte:"len:5"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord08Packet struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3038"`
	Date       []byte `byte:"len:10"`
	RoomNumber []byte `byte:"len:4"`
	Time       []byte `byte:"len:5"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord08Packetk struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Record_    []byte `byte:"len:2,equal:0x3038"`
	Date       []byte `byte:"len:5"`
	RoomNumber []byte `byte:"len:4"`
	Time       []byte `byte:"len:5"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRecord41FrameStartPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:2,equal:0x3431"`
	State   []byte `byte:"len:1,equal:0x31"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplRecord41FrameStoppPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:2,equal:0x3431"`
	State   []byte `byte:"len:1,equal:0x30"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

// not supported

type TplRecord06Packet struct { // time syncrho resp
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:2,equal:0x3036"`
	Day     []byte `byte:"len:2"`
	Dot0_   []byte `byte:"len:1,equal:0x2e"`
	Month   []byte `byte:"len:2"`
	Dot1_   []byte `byte:"len:1,equal:0x2e"`
	Year    []byte `byte:"len:4"`
	Hour    []byte `byte:"len:2"`
	Colon0_ []byte `byte:"len:1,equal:0x3a"`
	Minute  []byte `byte:"len:2"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}

type TplRecord56Packet struct { // time synchro req
	STX_    []byte `byte:"len:1,equal:0x02"`
	Record_ []byte `byte:"len:2,equal:0x3536"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}
