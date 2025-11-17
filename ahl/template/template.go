package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckOut           = "Check Out"
	PacketCheckIn            = "Check In"
	PacketDataChange         = "Modification"
	PacketWakeupSet          = "Wake-up set"
	PacketWakeupClear        = "Wake-up clear"
	PacketRoomStatus         = "Room Status"
	PacketWakeupEvent        = "Wake-up Event"
	PacketVoiceMailEvent     = "Voice Mail Event"
	PacketDataTransfer       = "Data Transfer"
	PacketCallPacket         = "Call Packet"
	PacketCallPacketExtended = "Call Packet Extended"
	PacketReply              = "Reply"

	PacketLinkStart = "Link Start"
	PacketLinkAlive = "Link Alive"

	PacketLogOutput       = "Log Output"
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

//

type TplLinkAliveIncoming struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ int    `byte:"len:1,equal:0x24"` // $
	Node []byte `byte:"len:4"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplLinkStartOutgoing struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Cmd_  int    `byte:"len:1,equal:0x40"`       // @
	Node_ []byte `byte:"len:4,equal:0x46464646"` // FFFF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplLinkAliveOutgoing struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	Cmd_  int    `byte:"len:1,equal:0x24"`       // $
	Node_ []byte `byte:"len:4,equal:0x46464646"` // FFFF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplReply_4400_5 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x52"` // R
	Extension []byte `byte:"len:5"`            // 5
	Password  []byte `byte:"len:4"`            // 4
	Status    []byte `byte:"len:2"`            //
	LRC_      []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplReply_4400_8 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x52"` // R
	Extension []byte `byte:"len:8"`            // 8
	Password  []byte `byte:"len:4"`            // 4
	Status    []byte `byte:"len:2"`            //
	LRC_      []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplCallPacket_4400_5 struct {
	STX_            []byte `byte:"len:1,equal:0x02"`
	Cmd_            int    `byte:"len:1,equal:0x4a"` // J
	Extension       []byte `byte:"len:5"`            // 5
	CostCenter      []byte `byte:"len:4"`            //
	CallType        []byte `byte:"len:1"`            //
	ActingExtension []byte `byte:"len:5"`            // 5
	Date            []byte `byte:"len:11"`           //
	Duration        []byte `byte:"len:5"`            // mmmss
	Cost            []byte `byte:"len:8"`            // float ?
	TrunkGroup      []byte `byte:"len:4"`            //
	CallingNumber   []byte `byte:"len:20"`           //
	LRC_            []byte `byte:"len:2"`            //
	ETX_            []byte `byte:"len:1,equal:0x03"`
}

type TplCallPacket_4400_8 struct {
	STX_            []byte `byte:"len:1,equal:0x02"`
	Cmd_            int    `byte:"len:1,equal:0x4a"` // J
	Extension       []byte `byte:"len:8"`            // 8
	CostCenter      []byte `byte:"len:4"`            //
	CallType        []byte `byte:"len:1"`            //
	ActingExtension []byte `byte:"len:8"`            // 8
	Date            []byte `byte:"len:11"`           //
	Duration        []byte `byte:"len:5"`            // mmmss
	Cost            []byte `byte:"len:8"`            // float ?
	TrunkGroup      []byte `byte:"len:4"`            //
	CallingNumber   []byte `byte:"len:20"`           //
	LRC_            []byte `byte:"len:2"`            //
	ETX_            []byte `byte:"len:1,equal:0x03"`
}

type TplCallPacket_4400_8_Extended struct {
	STX_                   []byte `byte:"len:1,equal:0x02"`
	Cmd_                   int    `byte:"len:1,equal:0x4b"` // K
	Extension              []byte `byte:"len:8"`            // 8
	ChargedUserNode        []byte `byte:"len:6"`            //
	ChargedUserSubAddr     []byte `byte:"len:20"`           //
	CostCenterNumber       []byte `byte:"len:4"`            //
	CostCenterName         []byte `byte:"len:10"`           //
	TransferringParty      []byte `byte:"len:8"`            //
	TransferringPartyNode  []byte `byte:"len:6"`            //
	ActingExtension        []byte `byte:"len:8"`            //
	ActingExtensionNode    []byte `byte:"len:6"`            //
	CallType               []byte `byte:"len:1"`            //
	Date                   []byte `byte:"len:15"`           //
	Duration               []byte `byte:"len:6"`            //
	Cost                   []byte `byte:"len:8"`            // float ?
	TrunkNumber            []byte `byte:"len:4"`            //
	TrunkGroupNumber       []byte `byte:"len:4"`            //
	TrunkGroupNode         []byte `byte:"len:6"`            //
	CallingNumber          []byte `byte:"len:30"`           //
	CallUse                int    `byte:"len:1"`            //
	AccessCode             []byte `byte:"len:12"`           //
	UserToUser             []byte `byte:"len:5"`            //
	ExternalFacilities     []byte `byte:"len:15"`           //
	InternalFacilities     []byte `byte:"len:20"`           //
	Carrier                []byte `byte:"len:2"`            //
	InitialDialedNumber    []byte `byte:"len:30"`           //
	WaitingDuration        []byte `byte:"len:4"`            //
	EffectiveDuration      []byte `byte:"len:6"`            //
	RedirectedCalIndicator int    `byte:"len:1"`            //
	LRC_                   []byte `byte:"len:2"`            //
	ETX_                   []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatus_4400_5 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x43"` // C
	Extension []byte `byte:"len:5"`            // 5
	Code      []byte `byte:"len:4"`            //
	Status    int    `byte:"len:1"`            //
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatus_4400_8 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x43"` // C
	Extension []byte `byte:"len:8"`            // 8
	Code      []byte `byte:"len:4"`            //
	Status    int    `byte:"len:1"`            //
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupEvent_4400_5 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Cmd_           int    `byte:"len:1,equal:0x50"` // P
	Extension      []byte `byte:"len:5"`            // 5
	Type_          int    `byte:"len:1,equal:0x57"` // W
	Code           []byte `byte:"len:4"`            //
	ActivationDate []byte `byte:"len:10"`           // ddmmyyhhmm
	Originator     []byte `byte:"len:5"`            //
	WakeupTime     []byte `byte:"len:5"`            // hhmm
	LRC            []byte `byte:"len:2"`            //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupEvent_4400_8 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Cmd_           int    `byte:"len:1,equal:0x50"` // P
	Extension      []byte `byte:"len:8"`            // 8
	Type_          int    `byte:"len:1,equal:0x57"` // W
	Code           []byte `byte:"len:4"`            //
	ActivationDate []byte `byte:"len:10"`           // ddmmyyhhmm
	Originator     []byte `byte:"len:5"`            //
	WakeupTime     []byte `byte:"len:5"`            // hhmm
	LRC            []byte `byte:"len:2"`            //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplDataTransfer_4400_5 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x50"` // P
	Extension []byte `byte:"len:5"`            // 5
	Code      []byte `byte:"len:5"`            //
	Data      []byte `byte:"len:20"`           //
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplDataTransfer_4400_8 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x50"` // P
	Extension []byte `byte:"len:8"`            // 8
	Code      []byte `byte:"len:5"`            //
	Data      []byte `byte:"len:20"`           //
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplVoiceMailEvent_4400_5 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x46"` // F
	Extension []byte `byte:"len:5"`            // 5
	Code      int    `byte:"len:1"`            // M incoming, P outgoing
	Status    int    `byte:"len:1"`            //
	NotUsed_  []byte `byte:"len:5"`            // 5
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplVoiceMailEvent_4400_8 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x46"` // F
	Extension []byte `byte:"len:8"`            // 8
	Code      int    `byte:"len:1"`            // M incoming, P outgoing
	Status    int    `byte:"len:1"`            //
	NotUsed_  []byte `byte:"len:5"`            // 5
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplCheckIn_4400_5 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Cmd_           int    `byte:"len:1,equal:0x41"` // A
	Extension      []byte `byte:"len:5"`            // 5
	Occupation     int    `byte:"len:1"`            //
	GuestName      []byte `byte:"len:20"`           //
	GuestLanguage  int    `byte:"len:1"`            //
	VIPState       int    `byte:"len:1"`            //
	GroupName      []byte `byte:"len:3"`            //
	Password       []byte `byte:"len:4"`            // 4
	DOD            []byte `byte:"len:2"`            //
	DepositAmount  []byte `byte:"len:9"`            //
	MessageWaiting int    `byte:"len:1"`            //
	WakeupTime     []byte `byte:"len:5"`            // hhmm
	DND            int    `byte:"len:1"`            //
	LRC            []byte `byte:"len:2"`            //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplCheckIn_4400_8 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Cmd_           int    `byte:"len:1,equal:0x41"` // A
	Extension      []byte `byte:"len:8"`            // 8
	Occupation     int    `byte:"len:1"`            //
	GuestName      []byte `byte:"len:20"`           //
	GuestLanguage  int    `byte:"len:1"`            //
	VIPState       int    `byte:"len:1"`            //
	GroupName      []byte `byte:"len:3"`            //
	Password       []byte `byte:"len:4"`            // 4
	DOD            []byte `byte:"len:2"`            //
	DepositAmount  []byte `byte:"len:9"`            //
	MessageWaiting int    `byte:"len:1"`            //
	WakeupTime     []byte `byte:"len:5"`            // hhmm
	DND            int    `byte:"len:1"`            //
	LRC            []byte `byte:"len:2"`            //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplCheckOut_4400_5 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x44"` // D
	Extension []byte `byte:"len:5"`            // 5
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplCheckOut_4400_8 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      int    `byte:"len:1,equal:0x44"` // D
	Extension []byte `byte:"len:8"`            // 8
	LRC       []byte `byte:"len:2"`            //
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplDataChange_4400_5 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Cmd_           int    `byte:"len:1,equal:0x4d"` // M
	Extension      []byte `byte:"len:5"`            // 5
	Occupation     int    `byte:"len:1"`            //
	GuestName      []byte `byte:"len:20"`           //
	GuestLanguage  int    `byte:"len:1"`            //
	VIPState       int    `byte:"len:1"`            //
	GroupName      []byte `byte:"len:3"`            //
	Password       []byte `byte:"len:4"`            // 4
	DOD            []byte `byte:"len:2"`            //
	DepositAmount  []byte `byte:"len:9"`            //
	MessageWaiting int    `byte:"len:1"`            //
	WakeupTime     []byte `byte:"len:5"`            // hhmm
	DND            int    `byte:"len:1"`            //
	LRC            []byte `byte:"len:2"`            //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplDataChange_4400_8 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Cmd_           int    `byte:"len:1,equal:0x4d"` // M
	Extension      []byte `byte:"len:8"`            // 8
	Occupation     int    `byte:"len:1"`            //
	GuestName      []byte `byte:"len:20"`           //
	GuestLanguage  int    `byte:"len:1"`            //
	VIPState       int    `byte:"len:1"`            //
	GroupName      []byte `byte:"len:3"`            //
	Password       []byte `byte:"len:4"`            // 4
	DOD            []byte `byte:"len:2"`            //
	DepositAmount  []byte `byte:"len:9"`            //
	MessageWaiting int    `byte:"len:1"`            //
	WakeupTime     []byte `byte:"len:5"`            // hhmm
	DND            int    `byte:"len:1"`            //
	LRC            []byte `byte:"len:2"`            //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

// filter log output

/*
0000000000  32 39 2F 30 36 2F 32 30 20 31 31 3A 31 32 20 20  29/06/20 11:12
0000000010  20 4D 4F 44 49 46 49 43 41 54 49 4F 4E 20 52 4F   MODIFICATION RO
0000000020  4F 4D 20 38 31 32 39 20 20 20 20 20 43 4F 44 45  OM 8129     CODE
0000000030  20 31 20 20 20 20 52 4F 4F 4D 20 53 54 41 54 55   1    ROOM STATU
0000000040  53 20 35 0D 0A
*/

type TplLogOutputIncoming1 struct {
	Day_    []byte `byte:"len:2"`
	Sep1_   []byte `byte:"len:1,equal:0x2f"`
	Month_  []byte `byte:"len:2"`
	Sep2_   []byte `byte:"len:1,equal:0x2f"`
	Year_   []byte `byte:"len:2"`
	Sep3_   []byte `byte:"len:1,equal:0x20"`
	Hour_   []byte `byte:"len:2"`
	Sep4_   []byte `byte:"len:1,equal:0x3a"`
	Minute_ []byte `byte:"len:2"`
	Sep5_   []byte `byte:"len:3,equal:0x202020"`
	Data_   []byte `byte:"len:*"`
	End1_   []byte `byte:"len:1,equal:0x0d"`
	End2_   []byte `byte:"len:1,equal:0x0a"`
}

/*
0000000000  30 35 2F 30 38 2F 32 31 20 31 34 3A 32 31 20 20  05/08/21 14:21
0000000010  3D 3D 20 43 48 45 43 4B 2D 4F 55 54 20 31 32 31  == CHECK-OUT 121
0000000020  38 20 20 20 20 43 3D 30 2E 30 30 20 20 20 20 20  8    C=0.00
0000000030  41 3D 30 2E 30 30 20 20 20 20 20 0D 0A           A=0.00     ..
*/

type TplLogOutputIncoming2 struct {
	Day_    []byte `byte:"len:2"`
	Sep1_   []byte `byte:"len:1,equal:0x2f"`
	Month_  []byte `byte:"len:2"`
	Sep2_   []byte `byte:"len:1,equal:0x2f"`
	Year_   []byte `byte:"len:2"`
	Sep3_   []byte `byte:"len:1,equal:0x20"`
	Hour_   []byte `byte:"len:2"`
	Sep4_   []byte `byte:"len:1,equal:0x3a"`
	Minute_ []byte `byte:"len:2"`
	Sep5_   []byte `byte:"len:3,equal:0x20203d"`
	Data_   []byte `byte:"len:*"`
	End1_   []byte `byte:"len:1,equal:0x0d"`
	End2_   []byte `byte:"len:1,equal:0x0a"`
}
