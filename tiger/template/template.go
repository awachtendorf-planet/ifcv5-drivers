package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketRoomBasedCall      = "Room Based Call Detail"
	PacketMiniBarBilling     = "Mini Bar Billing"
	PacketOtherChargePosting = "Other Charge Posting"

	PacketMessageWaitingGuest       = "Message Waiting (Guest)"
	PacketMessageWaitingReservation = "Message Waiting (Reservation)"

	PacketRoomBasedCheckIn  = "Room Based Check In"
	PacketRoomBasedCheckOut = "Room Based Check Out"
	PacketAdditionalGuest   = "Additional Guest Check In"

	PacketWakeUpSet         = "Wake Up Set"
	PacketWakeUpClear       = "Wake Up Clear"
	PacketInformationUpdate = "Information Update"
	PacketVideoRights       = "Video Rights"

	PacketRoomstatus   = "Room Status"
	PacketRoomtransfer = "Room Transfer"

	PacketDND = "Do Not Disturb"

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
	STX_  []byte `byte:"len:1,equal:0x02"`
	Data_ []byte `byte:"len:*"`
	STX2_ []byte `byte:"len:1,equal:0x02"`
}

//Incoming Requests

type TplRoomBasedCallDetailPacket struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x4352"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	Date          []byte `byte:"len:6"`
	Space0_       []byte `byte:"len:1,equal:0x20"`
	Time          []byte `byte:"len:4"`
	Space1_       []byte `byte:"len:1,equal:0x20"`
	RoomNumber    []byte `byte:"len:*"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	Duration      []byte `byte:"len:6"`
	Space3_       []byte `byte:"len:1,equal:0x20"`
	DialledNumber []byte `byte:"len:18"`
	Space4_       []byte `byte:"len:1,equal:0x20"`
	Cost          []byte `byte:"len:10"`
	Space5_       []byte `byte:"len:1,equal:0x20"`
	Sequence      []byte `byte:"len:5"`
	Space6_       []byte `byte:"len:1,equal:0x20"`
	ConditionCode []byte `byte:"len:2"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplRoomBasedCallDetailPacketk struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x4352"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	Date          []byte `byte:"len:4"`
	Space0_       []byte `byte:"len:1,equal:0x20"`
	Time          []byte `byte:"len:4"`
	Space1_       []byte `byte:"len:1,equal:0x20"`
	RoomNumber    []byte `byte:"len:*"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	Duration      []byte `byte:"len:6"`
	Space3_       []byte `byte:"len:1,equal:0x20"`
	DialledNumber []byte `byte:"len:18"`
	Space4_       []byte `byte:"len:1,equal:0x20"`
	Cost          []byte `byte:"len:10"`
	Space5_       []byte `byte:"len:1,equal:0x20"`
	Sequence      []byte `byte:"len:5"`
	Space6_       []byte `byte:"len:1,equal:0x20"`
	ConditionCode []byte `byte:"len:2"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplMiniBarBillingPacket struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record_    []byte `byte:"len:2,equal:0x4D42"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	Date       []byte `byte:"len:6"`
	Space0_    []byte `byte:"len:1,equal:0x20"`
	Time       []byte `byte:"len:4"`
	Space1_    []byte `byte:"len:1,equal:0x20"`
	RoomNumber []byte `byte:"len:*"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	ID         []byte `byte:"len:3"`
	Space3_    []byte `byte:"len:1,equal:0x20"`
	PPID       []byte `byte:"len:4"`
	Space4_    []byte `byte:"len:1,equal:0x20"`
	Quantity   []byte `byte:"len:3"`
	Space5_    []byte `byte:"len:1,equal:0x20"`
	Cost       []byte `byte:"len:10"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplMiniBarBillingPacketk struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record_    []byte `byte:"len:2,equal:0x4D42"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	Date       []byte `byte:"len:4"`
	Space0_    []byte `byte:"len:1,equal:0x20"`
	Time       []byte `byte:"len:4"`
	Space1_    []byte `byte:"len:1,equal:0x20"`
	RoomNumber []byte `byte:"len:*"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	ID         []byte `byte:"len:3"`
	Space3_    []byte `byte:"len:1,equal:0x20"`
	PPID       []byte `byte:"len:4"`
	Space4_    []byte `byte:"len:1,equal:0x20"`
	Quantity   []byte `byte:"len:3"`
	Space5_    []byte `byte:"len:1,equal:0x20"`
	Cost       []byte `byte:"len:10"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplOtherChargingPacket struct {
	SOH_        []byte `byte:"len:1,equal:0x01"`
	Record_     []byte `byte:"len:2,equal:0x4350"`
	STX_        []byte `byte:"len:1,equal:0x02"`
	Date        []byte `byte:"len:6"`
	Space0_     []byte `byte:"len:1,equal:0x20"`
	Time        []byte `byte:"len:4"`
	Space1_     []byte `byte:"len:1,equal:0x20"`
	PinNumber   []byte `byte:"len:5"`
	Space2_     []byte `byte:"len:1,equal:0x20"`
	RoomNumber  []byte `byte:"len:*"`
	Space3_     []byte `byte:"len:1,equal:0x20"`
	ID          []byte `byte:"len:3"`
	Space4_     []byte `byte:"len:1,equal:0x20"`
	PPID        []byte `byte:"len:4"`
	Space5_     []byte `byte:"len:1,equal:0x20"`
	Description []byte `byte:"len:18"`
	Space6_     []byte `byte:"len:1,equal:0x20"`
	Quantity    []byte `byte:"len:3"`
	Space7_     []byte `byte:"len:1,equal:0x20"`
	Cost        []byte `byte:"len:10"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplOtherChargingPacketk struct {
	SOH_        []byte `byte:"len:1,equal:0x01"`
	Record_     []byte `byte:"len:2,equal:0x4350"`
	STX_        []byte `byte:"len:1,equal:0x02"`
	Date        []byte `byte:"len:4"`
	Space0_     []byte `byte:"len:1,equal:0x20"`
	Time        []byte `byte:"len:4"`
	Space1_     []byte `byte:"len:1,equal:0x20"`
	PinNumber   []byte `byte:"len:5"`
	Space2_     []byte `byte:"len:1,equal:0x20"`
	RoomNumber  []byte `byte:"len:*"`
	Space3_     []byte `byte:"len:1,equal:0x20"`
	ID          []byte `byte:"len:3"`
	Space4_     []byte `byte:"len:1,equal:0x20"`
	PPID        []byte `byte:"len:4"`
	Space5_     []byte `byte:"len:1,equal:0x20"`
	Description []byte `byte:"len:18"`
	Space6_     []byte `byte:"len:1,equal:0x20"`
	Quantity    []byte `byte:"len:3"`
	Space7_     []byte `byte:"len:1,equal:0x20"`
	Cost        []byte `byte:"len:10"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplMessageWaitingReservationPacket struct {
	SOH_                []byte `byte:"len:1,equal:0x01"`
	Record_             []byte `byte:"len:2,equal:0x4D52"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	Date                []byte `byte:"len:6"`
	Space0_             []byte `byte:"len:1,equal:0x20"`
	Time                []byte `byte:"len:4"`
	Space1_             []byte `byte:"len:1,equal:0x20"`
	ReservationNumber   []byte `byte:"len:10"`
	Space2_             []byte `byte:"len:1,equal:0x20"`
	ActionCode          []byte `byte:"len:1"`
	Space3_             []byte `byte:"len:1,equal:0x20"`
	MessageWaitingCount []byte `byte:"len:2"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

// Outgoing

type TplCheckInOutRoomBasedPacket struct {
	SOH_                []byte `byte:"len:1,equal:0x01"`
	Record_             []byte `byte:"len:2,equal:0x4C4E"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	Date                []byte `byte:"len:6"`
	Space0_             []byte `byte:"len:1,equal:0x20"`
	Time                []byte `byte:"len:4"`
	Space1_             []byte `byte:"len:1,equal:0x20"`
	RoomNumber          []byte `byte:"len:*"`
	Space2_             []byte `byte:"len:1,equal:0x20"`
	ActionCode          []byte `byte:"len:1"`
	Space3_             []byte `byte:"len:1,equal:0x20"`
	VipStatus           []byte `byte:"len:2"`
	Space4_             []byte `byte:"len:1,equal:0x20"`
	GroupID             []byte `byte:"len:4"`
	Space5_             []byte `byte:"len:1,equal:0x20"`
	AccountNumber       []byte `byte:"len:6"`
	Space6_             []byte `byte:"len:1,equal:0x20"`
	ReservationNumber   []byte `byte:"len:10"`
	Space7_             []byte `byte:"len:1,equal:0x20"`
	LanguageDescription []byte `byte:"len:14"`
	LanguageCode        []byte `byte:"len:2"`
	Space8_             []byte `byte:"len:1,equal:0x20"`
	GuestName           []byte `byte:"len:*"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type TplAdditionalGuestPacket struct {
	SOH_                []byte `byte:"len:1,equal:0x01"`
	Record_             []byte `byte:"len:2,equal:0x4147"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	Date                []byte `byte:"len:6"`
	Space0_             []byte `byte:"len:1,equal:0x20"`
	Time                []byte `byte:"len:4"`
	Space1_             []byte `byte:"len:1,equal:0x20"`
	PinNumber           []byte `byte:"len:5"`
	Space2_             []byte `byte:"len:1,equal:0x20"`
	RoomNumber          []byte `byte:"len:*"`
	Space3_             []byte `byte:"len:1,equal:0x20"`
	ActionCode          []byte `byte:"len:1"`
	Space4_             []byte `byte:"len:1,equal:0x20"`
	VipStatus           []byte `byte:"len:2"`
	Space5_             []byte `byte:"len:1,equal:0x20"`
	GroupID             []byte `byte:"len:4"`
	Space6_             []byte `byte:"len:1,equal:0x20"`
	AccountNumber       []byte `byte:"len:6"`
	Space7_             []byte `byte:"len:1,equal:0x20"`
	ReservationNumber   []byte `byte:"len:10"`
	Space8_             []byte `byte:"len:1,equal:0x20"`
	LanguageDescription []byte `byte:"len:14"`
	LanguageCode        []byte `byte:"len:2"`
	Space9_             []byte `byte:"len:1,equal:0x20"`
	GuestName           []byte `byte:"len:*"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type TplWakeUpPacket struct {
	SOH_         []byte `byte:"len:1,equal:0x01"`
	Record_      []byte `byte:"len:2,equal:0x574B"`
	STX_         []byte `byte:"len:1,equal:0x02"`
	DateStamp    []byte `byte:"len:6"`
	Space0_      []byte `byte:"len:1,equal:0x20"`
	TimeStamp    []byte `byte:"len:4"`
	Space1_      []byte `byte:"len:1,equal:0x20"`
	DateCall     []byte `byte:"len:6"`
	Space2_      []byte `byte:"len:1,equal:0x20"`
	TimeCall     []byte `byte:"len:4"`
	Space3_      []byte `byte:"len:1,equal:0x20"`
	PinNumber    []byte `byte:"len:5"`
	Space4_      []byte `byte:"len:1,equal:0x20"`
	RoomNumber   []byte `byte:"len:*"`
	Space5_      []byte `byte:"len:1,equal:0x20"`
	GroupID      []byte `byte:"len:4"`
	Space6_      []byte `byte:"len:1,equal:0x20"`
	LanguageCode []byte `byte:"len:2"`
	Space7_      []byte `byte:"len:1,equal:0x20"`
	ActionCode   []byte `byte:"len:1"`
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

type TplRoomTransferPacket struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x5452"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	Date          []byte `byte:"len:6"`
	Space0_       []byte `byte:"len:1,equal:0x20"`
	Time          []byte `byte:"len:4"`
	Space1_       []byte `byte:"len:1,equal:0x20"`
	RoomNumberOld []byte `byte:"len:*"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	RoomNumberNew []byte `byte:"len:*"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplInformationUpdatePacket struct {
	SOH_                 []byte `byte:"len:1,equal:0x01"`
	Record_              []byte `byte:"len:2,equal:0x5352"`
	STX_                 []byte `byte:"len:1,equal:0x02"`
	Date                 []byte `byte:"len:6"`
	Space0_              []byte `byte:"len:1,equal:0x20"`
	Time                 []byte `byte:"len:4"`
	Space1_              []byte `byte:"len:1,equal:0x20"`
	RoomNumber           []byte `byte:"len:*"`
	Space2_              []byte `byte:"len:1,equal:0x20"`
	VipStatus            []byte `byte:"len:2"`
	Space3_              []byte `byte:"len:1,equal:0x20"`
	GroupID              []byte `byte:"len:4"`
	Space4_              []byte `byte:"len:1,equal:0x20"`
	AccountNumber        []byte `byte:"len:6"`
	Space5_              []byte `byte:"len:1,equal:0x20"`
	ReservationNumber    []byte `byte:"len:10"`
	Space6_              []byte `byte:"len:1,equal:0x20"`
	FirstWake            []byte `byte:"len:4"`
	Space7_              []byte `byte:"len:1,equal:0x20"`
	SecondWake           []byte `byte:"len:4"`
	Space8_              []byte `byte:"len:1,equal:0x20"`
	LanguageDescription  []byte `byte:"len:14"`
	LanguageCode         []byte `byte:"len:2"`
	Space9_              []byte `byte:"len:1,equal:0x20"`
	Lvl9Occupancy        []byte `byte:"len:1"`
	Space10_             []byte `byte:"len:1,equal:0x20"`
	MessageWaitingAction []byte `byte:"len:1"`
	Space11_             []byte `byte:"len:1,equal:0x20"`
	RoomStatusCode       []byte `byte:"len:2"`
	Space12_             []byte `byte:"len:1,equal:0x20"`
	GuestName            []byte `byte:"len:*"`
	ETX_                 []byte `byte:"len:1,equal:0x03"`
}

type TplVideoRightsPacket struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x5652"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	Date          []byte `byte:"len:6"`
	Space0_       []byte `byte:"len:1,equal:0x20"`
	Time          []byte `byte:"len:4"`
	Space1_       []byte `byte:"len:1,equal:0x20"`
	RoomNumber    []byte `byte:"len:*"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	AccountNumber []byte `byte:"len:6"`
	Space3_       []byte `byte:"len:1,equal:0x20"`
	ViewBill      []byte `byte:"len:1"`
	Space4_       []byte `byte:"len:1,equal:0x20"`
	VideoCheckOut []byte `byte:"len:1"`
	Space5_       []byte `byte:"len:1,equal:0x20"`
	CFlag         []byte `byte:"len:1"`
	Space6_       []byte `byte:"len:1,equal:0x20"`
	DFlag         []byte `byte:"len:1"`
	Space7_       []byte `byte:"len:1,equal:0x20"`
	EFlag         []byte `byte:"len:1"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplDoNotDisturbPacket struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record_    []byte `byte:"len:2,equal:0x444E"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	Date       []byte `byte:"len:6"`
	Space0_    []byte `byte:"len:1,equal:0x20"`
	Time       []byte `byte:"len:4"`
	Space1_    []byte `byte:"len:1,equal:0x20"`
	RoomNumber []byte `byte:"len:*"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	DND        []byte `byte:"len:1"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

// General

type TplRoomStatusPacket struct {
	SOH_           []byte `byte:"len:1,equal:0x01"`
	Record_        []byte `byte:"len:2,equal:0x5253"`
	STX_           []byte `byte:"len:1,equal:0x02"`
	Date           []byte `byte:"len:6"`
	Space0_        []byte `byte:"len:1,equal:0x20"`
	Time           []byte `byte:"len:4"`
	Space1_        []byte `byte:"len:1,equal:0x20"`
	RoomNumber     []byte `byte:"len:*"`
	Space2_        []byte `byte:"len:1,equal:0x20"`
	RoomStatusCode []byte `byte:"len:2"`
	Space3_        []byte `byte:"len:1,equal:0x20"`
	ID             []byte `byte:"len:3"`
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplMessageWaitingGuestPacket struct {
	SOH_                []byte `byte:"len:1,equal:0x01"`
	Record_             []byte `byte:"len:2,equal:0x4D57"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	Date                []byte `byte:"len:6"`
	Space0_             []byte `byte:"len:1,equal:0x20"`
	Time                []byte `byte:"len:4"`
	Space1_             []byte `byte:"len:1,equal:0x20"`
	RoomNumber          []byte `byte:"len:*"`
	Space2_             []byte `byte:"len:1,equal:0x20"`
	ActionCode          []byte `byte:"len:1"`
	Space3_             []byte `byte:"len:1,equal:0x20"`
	MessageWaitingCount []byte `byte:"len:2"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 															Opera Special Case															//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//Incoming Requests

type TplRoomBasedCallDetailPacketOpera struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x4352"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	Time          []byte `byte:"len:4"`
	Space1_       []byte `byte:"len:1,equal:0x20"`
	RoomNumber    []byte `byte:"len:5"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	Duration      []byte `byte:"len:*"`
	Space3_       []byte `byte:"len:1,equal:0x20"`
	DialledNumber []byte `byte:"len:18"`
	Space4_       []byte `byte:"len:1,equal:0x20"`
	Cost          []byte `byte:"len:6"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplMessageWaitingPacketOpera struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record_    []byte `byte:"len:2,equal:0x4D52"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	RoomNumber []byte `byte:"len:5"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	ActionCode []byte `byte:"len:1"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusPacketOpera struct {
	SOH_           []byte `byte:"len:1,equal:0x01"`
	Record_        []byte `byte:"len:2,equal:0x5253"`
	STX_           []byte `byte:"len:1,equal:0x02"`
	Time           []byte `byte:"len:4"`
	Space1_        []byte `byte:"len:1,equal:0x20"`
	RoomNumber     []byte `byte:"len:5"`
	Space2_        []byte `byte:"len:1,equal:0x20"`
	RoomStatusCode []byte `byte:"len:2"`
	Space3_        []byte `byte:"len:1,equal:0x20"`
	ID             []byte `byte:"len:3"`
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

// Outgoing

type TplCheckInRoomBasedPacketOpera0 struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record     []byte `byte:"len:2"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	RoomNumber []byte `byte:"len:5"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	ActionCode []byte `byte:"len:1"`
	Space3_    []byte `byte:"len:1,equal:0x20"`
	GuestName  []byte `byte:"len:*"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplCheckInRoomBasedPacketOpera12 struct {
	SOH_                []byte `byte:"len:1,equal:0x01"`
	Record              []byte `byte:"len:2"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	RoomNumber          []byte `byte:"len:5"`
	Space2_             []byte `byte:"len:1,equal:0x20"`
	ActionCode          []byte `byte:"len:1"`
	Space3_             []byte `byte:"len:1,equal:0x20"`
	ReservationNumber   []byte `byte:"len:*"`
	Space4_             []byte `byte:"len:1,equal:0x20"`
	LanguageDescription []byte `byte:"len:14"`
	LanguageCode        []byte `byte:"len:2"`
	Space5_             []byte `byte:"len:1,equal:0x20"`
	GuestName           []byte `byte:"len:*"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type TplCheckOutRoomBasedPacketOpera0 struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record     []byte `byte:"len:2"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	RoomNumber []byte `byte:"len:5"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	ActionCode []byte `byte:"len:1,equal:0x32"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplCheckOutRoomBasedPacketOpera12 struct {
	SOH_                []byte `byte:"len:1,equal:0x01"`
	Record              []byte `byte:"len:2"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	RoomNumber          []byte `byte:"len:5"`
	Space2_             []byte `byte:"len:1,equal:0x20"`
	ActionCode          []byte `byte:"len:1,equal:0x32"`
	Space3_             []byte `byte:"len:1,equal:0x20"`
	ReservationNumber   []byte `byte:"len:*"`
	Space4_             []byte `byte:"len:1,equal:0x20"`
	LanguageDescription []byte `byte:"len:14"`
	LanguageCode        []byte `byte:"len:2"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
}

type TplWakeUpPacketOpera struct {
	SOH_         []byte `byte:"len:1,equal:0x01"`
	Record_      []byte `byte:"len:2,equal:0x574B"`
	STX_         []byte `byte:"len:1,equal:0x02"`
	TimeCall     []byte `byte:"len:4"`
	Space3_      []byte `byte:"len:1,equal:0x20"`
	RoomNumber   []byte `byte:"len:5"`
	Space5_      []byte `byte:"len:1,equal:0x20"`
	LanguageCode []byte `byte:"len:2"`
	Space7_      []byte `byte:"len:1,equal:0x20"`
	ActionCode   []byte `byte:"len:1,"`
	ETX_         []byte `byte:"len:1,equal:0x03"`
}

type TplRoomTransferPacketOpera2 struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x5453"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	RoomNumberOld []byte `byte:"len:5"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	ReservationID []byte `byte:"len:8"`
	Space3_       []byte `byte:"len:1,equal:0x20"`
	RoomNumberNew []byte `byte:"len:5"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplRoomTransferPacketOpera01 struct {
	SOH_          []byte `byte:"len:1,equal:0x01"`
	Record_       []byte `byte:"len:2,equal:0x5452"`
	STX_          []byte `byte:"len:1,equal:0x02"`
	RoomNumberOld []byte `byte:"len:5"`
	Space2_       []byte `byte:"len:1,equal:0x20"`
	RoomNumberNew []byte `byte:"len:5"`
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplInformationUpdatePacketOpera0 struct {
	SOH_                 []byte `byte:"len:1,equal:0x01"`
	Record               []byte `byte:"len:2"`
	STX_                 []byte `byte:"len:1,equal:0x02"`
	RoomNumber           []byte `byte:"len:5"`
	Space2_              []byte `byte:"len:1,equal:0x20"`
	Magic_               []byte `byte:"len:9,equal:0x393939392039393939"`
	Space3_              []byte `byte:"len:1,equal:0x20"`
	LanguageCode         []byte `byte:"len:2"`
	Space9_              []byte `byte:"len:1,equal:0x20"`
	Lvl9Occupancy        []byte `byte:"len:1"` // 0unbarroomoccupied // 1barroomoccupied
	Space10_             []byte `byte:"len:1,equal:0x20"`
	MessageWaitingAction []byte `byte:"len:1"`
	Space12_             []byte `byte:"len:1,equal:0x20"`
	GuestName            []byte `byte:"len:*"`
	ETX_                 []byte `byte:"len:1,equal:0x03"`
}

type TplInformationUpdatePacketOpera12 struct {
	SOH_                 []byte `byte:"len:1,equal:0x01"`
	Record               []byte `byte:"len:2"`
	STX_                 []byte `byte:"len:1,equal:0x02"`
	RoomNumber           []byte `byte:"len:5"`
	Space0_              []byte `byte:"len:1,equal:0x20"`
	ReservationNumber    []byte `byte:"len:*"`
	Space1_              []byte `byte:"len:1,equal:0x20"`
	Magic_               []byte `byte:"len:9,equal:0x393939392039393939"`
	Space2_              []byte `byte:"len:1,equal:0x20"`
	LanguageDescription  []byte `byte:"len:14"`
	LanguageCode         []byte `byte:"len:2"`
	Space3_              []byte `byte:"len:1,equal:0x20"`
	Lvl9Occupancy        []byte `byte:"len:1"` // 0unbarroomoccupied // 1barroomoccupied
	Space4_              []byte `byte:"len:1,equal:0x20"`
	MessageWaitingAction []byte `byte:"len:1"`
	Space5_              []byte `byte:"len:1,equal:0x20"`
	RoomStatus_          []byte `byte:"len:2,equal:0x3939"`
	Space6_              []byte `byte:"len:1,equal:0x20"`
	GuestName            []byte `byte:"len:*"`
	ETX_                 []byte `byte:"len:1,equal:0x03"`
}

type TplDoNotDisturbPacketOpera struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record_    []byte `byte:"len:2,equal:0x444E"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	RoomNumber []byte `byte:"len:*"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	DND        []byte `byte:"len:1"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}

type TplMessageWaitingGuestPacketOpera struct {
	SOH_       []byte `byte:"len:1,equal:0x01"`
	Record_    []byte `byte:"len:2,equal:0x4D57"`
	STX_       []byte `byte:"len:1,equal:0x02"`
	RoomNumber []byte `byte:"len:5"`
	Space2_    []byte `byte:"len:1,equal:0x20"`
	ActionCode []byte `byte:"len:1"`
	ETX_       []byte `byte:"len:1,equal:0x03"`
}
