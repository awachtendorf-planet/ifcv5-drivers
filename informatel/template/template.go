package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckIn        = "Check In"
	PacketCheckOut       = "Check Out"
	PacketRoomMove       = "Room Move"
	PacketNameChange     = "Name Change"
	PacketClassOfService = "Class of Service"

	PacketRoomStatus = "Room Status"

	PacketMessageLightStatus = "Message Light Status"
	PacketDND                = "Do not Disturb"

	PacketWakeupOutgoingCall = "WakeupCall Outgoing"
	PacketWakeupIncomingCall = "WakeupCall Incoming"

	PacketChargeRecord  = "Charge Record"
	PacketMinibarCharge = "Minibar Charge"
	PacketVoicemail     = "VoiceMail"

	PacketDBSwap = "DBSwap"

	PacketVirtualNumberAssign   = "Virtual Number Assign"   // Not Implemented
	PacketVirtualNumberUnassign = "Virtual Number Unassign" // Not Implemented

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

// Outgoing

// CheckIn Example
// "INP1324               AF EQUIPAGES                          F                                                                                                                                                               y"

type CheckInPacket struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Type_        []byte `byte:"len:3,equal:0x494E50"` // INP
	Extension    []byte `byte:"len:10"`               // left justified
	Filler1_     []byte `byte:"len:9,equal:0x202020202020202020"`
	Name         []byte `byte:"len:20"`
	Filler2_     []byte `byte:"len:15,equal:0x202020202020202020202020202020"`
	VIP          []byte `byte:"len:3"`
	LanguageCode []byte `byte:"len:1"`
	LastFiller_  []byte `byte:"len:159,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"` // 221-62
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

type CheckOutPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Type_       []byte `byte:"len:3,equal:0x4F5554"`                                                                                                                                                                                                                                                                                                                                                                                                                           // OUT
	Extension   []byte `byte:"len:10"`                                                                                                                                                                                                                                                                                                                                                                                                                                         // left justified
	LastFiller_ []byte `byte:"len:207,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"` // 221-14
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type RoomMovePacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Type_         []byte `byte:"len:3,equal:0x44454C"` // DEL
	OldRoomnumber []byte `byte:"len:10"`               // left justified
	Filler1_      []byte `byte:"len:6,equal:0x202020202020"`
	Sharer_       []byte `byte:"len:1,equal:0x4E"` // No Sharer
	Filler2_      []byte `byte:"len:102,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	Extension     []byte `byte:"len:8"`
	LastFiller_   []byte `byte:"len:90,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"` // 221 - 131
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type NameChangePacket struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Type_        []byte `byte:"len:3,equal:0x4D4F44"` // MOD
	Extension    []byte `byte:"len:10"`               // left justified
	Filler1_     []byte `byte:"len:9,equal:0x202020202020202020"`
	Name         []byte `byte:"len:20"`
	Filler2_     []byte `byte:"len:15,equal:0x202020202020202020202020202020"`
	VIP          []byte `byte:"len:3"`
	LanguageCode []byte `byte:"len:1"`
	LastFiller_  []byte `byte:"len:159,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"` // 221-62
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

type VirtualNumberAssignPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Type_         []byte `byte:"len:3,equal:0x534441"` // SDA
	Extension     []byte `byte:"len:10"`               // left justified
	Filler_       []byte `byte:"len:109,equal:0x20202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	VirtualNumber []byte `byte:"len:8"`
	LastFiller_   []byte `byte:"len:90,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type VirtualNumberUnassignPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Type_         []byte `byte:"len:3,equal:0x534453"` // SDS
	Extension     []byte `byte:"len:10"`               // left justified
	Filler_       []byte `byte:"len:109,equal:0x20202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	VirtualNumber []byte `byte:"len:8"`
	LastFiller_   []byte `byte:"len:90,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type MessagePacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Type_         []byte `byte:"len:3,equal:0x4D5347"` // MSG
	Extension     []byte `byte:"len:10"`               // left justified
	Filler_       []byte `byte:"len:83,equal:0x2020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	MessageStatus []byte `byte:"len:3"`
	LastFiller_   []byte `byte:"len:121,equal:0x20202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type DNDPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Type_       []byte `byte:"len:3,equal:0x4E5044"` // NPD
	Extension   []byte `byte:"len:10"`               // left justified
	Filler_     []byte `byte:"len:83,equal:0x2020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	DND         []byte `byte:"len:3"`
	LastFiller_ []byte `byte:"len:121,equal:0x20202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type WakeupOutPacket struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Type_        []byte `byte:"len:3,equal:0x524556"` // REV
	Extension    []byte `byte:"len:10"`               // left justified
	Filler1_     []byte `byte:"len:67,equal:0x20202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	WakeupStatus []byte `byte:"len:3"`
	Filler2_     []byte `byte:"len:7,equal:0x20202020202020"`
	WakeupTime   []byte `byte:"len:4"` // HHMM
	LastFiller_  []byte `byte:"len:126,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

type ClassOfServicePacket struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Type_          []byte `byte:"len:3,equal:0x45544C"` // ETL
	Extension      []byte `byte:"len:10"`               // left justified
	Filler_        []byte `byte:"len:99,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ClassOfService []byte `byte:"len:3"`
	LastFiller_    []byte `byte:"len:105,equal:0x202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020"`
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

// Incoming

type ChargeRecordPacket struct {
	STX_                []byte `byte:"len:1,equal:0x02"`
	Type_               []byte `byte:"len:3,equal:0x544943"` // TIC
	Extension           []byte `byte:"len:10"`               // left justified
	RoomStatus          []byte `byte:"len:3"`
	Filler_             []byte `byte:"len:122"`
	DateTime            []byte `byte:"len:10"` // YYMMDDHHMM
	DialledDigits       []byte `byte:"len:16"`
	Units               []byte `byte:"len:5"`
	Duration            []byte `byte:"len:5"` // SSSSS
	LastFiller_         []byte `byte:"len:36"`
	TransactionDateTime []byte `byte:"len:10"` // YYMMDDHHMM
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type MinibarChargePacket struct {
	STX_                []byte `byte:"len:1,equal:0x02"`
	Type_               []byte `byte:"len:3,equal:0x424152"` // BAR
	Extension           []byte `byte:"len:10"`               // left justified
	Filler1_            []byte `byte:"len:6"`
	ArticleAmount       []byte `byte:"len:1"`
	ArticleNumber       []byte `byte:"len:3"`
	Filler2_            []byte `byte:"len:187"`
	TransactionDateTime []byte `byte:"len:10"` // YYMMDDHHMM
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type RoomStatusPacket struct {
	STX_       []byte `byte:"len:1,equal:0x02"`
	Type_      []byte `byte:"len:3,equal:0x455447"` // ETG
	Extension  []byte `byte:"len:10"`               // left justified
	Filler1_   []byte `byte:"len:3"`
	RoomStatus []byte `byte:"len:3"`
	Filler2_   []byte `byte:"len:*"` // Fkin Reality
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type VoiceMailPacket struct {
	STX_                []byte `byte:"len:1,equal:0x02"`
	Type_               []byte `byte:"len:3,equal:0x564F58"` // VOX
	Extension           []byte `byte:"len:10"`               // left justified`
	Filler1_            []byte `byte:"len:83"`
	VoiceMail           []byte `byte:"len:1"`
	Filler_             []byte `byte:"len:113"`
	TransactionDateTime []byte `byte:"len:10"` // YYMMDDHHMM
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type WakeupInPacket struct {
	STX_         []byte `byte:"len:1,equal:0x02"`
	Type_        []byte `byte:"len:3,equal:0x524556"` // REV
	Extension    []byte `byte:"len:10"`               // left justified
	Filler1_     []byte `byte:"len:67"`
	WakeupStatus []byte `byte:"len:2"`
	Filler2_     []byte `byte:"len:8"`
	WakeupTime   []byte `byte:"len:4"` // HHMM ! Handling if Call for Today or tomorrow
	LastFiller_  []byte `byte:"len:126"`
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

type DBSwapRequestPacket struct {
	STX_    []byte `byte:"len:1,equal:0x02"`
	Type_   []byte `byte:"len:3,equal:0x524553"` // RES
	Filler_ []byte `byte:"len:217"`
	ETX_    []byte `byte:"len:1,equal:0x03"`
}
