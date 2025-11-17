package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckInExtension  = "Check In Extension"
	PacketCheckOutExtension = "Check Out Extension"
	PacketDisplayName       = "Change Display Name"
	PacketRoomStatus        = "Room Status"
	PacketLanguage          = "Language"
	PacketMessageLamp       = "Message Lamp"
	PacketDoNotDisturb      = "Do Not Disturb"
	PacketVipState          = "VIP State"
	PacketClassOfService    = "Class Of Service"
	// PacketSetCCRS           = "CCRS Level"
	// PacketSetECC1           = "ECC1 Level"
	// PacketSetECC2           = "ECC2 Level"

	PacketMinibarItem      = "Minibar Item"
	PacketMinibarTotal     = "Minibar Total"
	PacketVoiceCount       = "Voice Count"
	PacketCallPacket       = "Call Packet"
	PacketCallPacketSingle = "Call Packet Single Line"

	PacketWakeupSet    = "Wakeup Set"
	PacketWakeupClear  = "Wakeup Clear"
	PacketWakeupAnswer = "Wakeup Answer"

	PacketError   = "Error Reply"
	PacketPolling = "Polling"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketTermination     = "Termination Sequence"
	PacketUnknown         = "Unknown Framed Packet"
	PacketUnknownBgd      = "Unknown Packet"
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

type TplCR struct {
	CR_ []byte `byte:"len:1,equal:0x0d"`
}

type TplCRSequence struct {
	CR_ []byte `byte:"len:7,equal:0x0a000000000000"`
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

/* ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

// incoming polling/error

type TplPolling struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:2,equal:0x5354"` // ST
	SEP_  []byte `byte:"len:1,equal:0x20"`   //
	Data_ []byte `byte:"len:2,equal:0x504f"` // PO
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplTest struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:2,equal:0x4953"`     // IS
	SEP_  []byte `byte:"len:1,equal:0x20"`       //
	Data_ []byte `byte:"len:4,equal:0x54455354"` // TEST
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplErrorMnemonic struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:8,equal:0x4d4e454d4f4e4943"` // MNEMONIC
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorInput struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:11,equal:0x494e505554204552524f52"` // INPUT ERROR
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorTryAgain struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:9,equal:0x54525920414741494e"` // TRY AGAIN
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorNameBig struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:8,equal:0x4e414d4520424947"` // NAME BIG
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorDuplicate struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:9,equal:0x4455504c4943415445"` // DUPLICATE
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorNoData struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:13,equal:0x4e4f204441544120464f554e44"` // NO DATA FOUND
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorNoSetCPNDData struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:16,equal:0x4e4f205345542043504e442044415441"` // NO SET CPND DATA
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorNoCustCPNDData struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:17,equal:0x4e4f20435553542043504e442044415441"` // NO CUST CPND DATA
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplErrorNoCPNDMemory struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Msg  []byte `byte:"len:14,equal:0x4e4f2043504e44204d454d4f5259"` // NO CPND MEMORY
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

/* ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

// outgoing

type TplCheckinExtension struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext  []byte `byte:"len:*"`                      //
	ACT_ []byte `byte:"len:6,equal:0x20434820494e"` // _CH_IN
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCheckoutExtension struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext  []byte `byte:"len:*"`                      //
	ACT_ []byte `byte:"len:6,equal:0x204348204f55"` // _CH_OU
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplSetDisplayName struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:6,equal:0x534520435020"` // SE_CP_
	Ext  []byte `byte:"len:*"`                      // extension
	SPC_ []byte `byte:"len:1,equal:0x20"`           // _
	Name []byte `byte:"len:*"`                      // display name
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplSetRoomStatus struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext   []byte `byte:"len:*"`                      // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`           // _
	State []byte `byte:"len:2"`                      // room status RE/PR/CL/PA/FA/SK/NS
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplSetLanguage struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext  []byte `byte:"len:*"`                      //
	ACT_ []byte `byte:"len:4,equal:0x204c4120"`     // _LA_
	Lang []byte `byte:"len:*"`                      // language code 3 or FR
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplSetMessageLamp struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext   []byte `byte:"len:*"`                      //
	ACT_  []byte `byte:"len:4,equal:0x204d5720"`     // _MW_
	State []byte `byte:"len:2"`                      // lamp ON/OF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplSetDoNotDisturb struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext   []byte `byte:"len:*"`                      //
	ACT_  []byte `byte:"len:4,equal:0x20444e20"`     // _DN_
	State []byte `byte:"len:2"`                      // dnd ON/OF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplSetVipState struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext   []byte `byte:"len:*"`                      //
	ACT_  []byte `byte:"len:4,equal:0x20564920"`     // _VI_
	State []byte `byte:"len:2"`                      // vip ON/OF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

// type TplSetCCRS struct {
// 	STX_  []byte `byte:"len:1,equal:0x02"`
// 	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
// 	Ext   []byte `byte:"len:*"`                      //
// 	ACT_  []byte `byte:"len:4,equal:0x20434f20"`     // _CO_
// 	State []byte `byte:"len:2"`                      // level ON/OF
// 	ETX_  []byte `byte:"len:1,equal:0x03"`
// }

// type TplSetECC1 struct {
// 	STX_  []byte `byte:"len:1,equal:0x02"`
// 	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
// 	Ext   []byte `byte:"len:*"`                      //
// 	ACT_  []byte `byte:"len:4,equal:0x20453120"`     // _E1_
// 	State []byte `byte:"len:2"`                      // level ON/OF
// 	ETX_  []byte `byte:"len:1,equal:0x03"`
// }

// type TplSetECC2 struct {
// 	STX_  []byte `byte:"len:1,equal:0x02"`
// 	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
// 	Ext   []byte `byte:"len:*"`                      //
// 	ACT_  []byte `byte:"len:4,equal:0x20453220"`     // _E2_
// 	State []byte `byte:"len:2"`                      // level ON/OF
// 	ETX_  []byte `byte:"len:1,equal:0x03"`
// }

type TplSetClassOfService struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:6,equal:0x534520535420"` // SE_ST_
	Ext   []byte `byte:"len:*"`                      //
	SPC1_ []byte `byte:"len:1,equal:0x20"`           // _
	Level []byte `byte:"len:2"`                      // level
	SPC2_ []byte `byte:"len:1,equal:0x20"`           // _
	State []byte `byte:"len:2"`                      // state
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplSetWakeup struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:6,equal:0x534520574120"` // SE_WA_
	Ext  []byte `byte:"len:*"`                      //
	ACT_ []byte `byte:"len:4,equal:0x20544920"`     // _TI_
	Time []byte `byte:"len:4"`                      // hhmm
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplClearWakeup struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:6,equal:0x534520574120"` // SE_WA_
	Ext  []byte `byte:"len:*"`                      //
	ACT_ []byte `byte:"len:4,equal:0x20544920"`     // _TI_
	CLR_ []byte `byte:"len:2,equal:0x4f46"`         // OF
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

/* ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

// incoming

type TplRoomStatus struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x535420"` // ST_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	State []byte `byte:"len:2"`                // room status RE/PR/CL/PA/FA/SK/NS
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusAdvanced struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x535420"` // ST_
	Ext   []byte `byte:"len:*"`                // extension
	SEP1  []byte `byte:"len:1,equal:0x20"`     // _
	State []byte `byte:"len:2"`                // room status RE/PR/CL/PA/FA/SK/NS
	SEP2  []byte `byte:"len:1,equal:0x20"`     // _
	MI_   []byte `byte:"len:3,equal:0x4d4920"` // MI_
	Maid  []byte `byte:"len:*"`                // maid
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplVoiceCount struct {
	STX_   []byte `byte:"len:1,equal:0x02"`
	CMD_   []byte `byte:"len:3,equal:0x495320"` // IS_
	ACT_   []byte `byte:"len:3,equal:0x564320"` // VC_
	Ext    []byte `byte:"len:*"`                // extension
	SEP1_  []byte `byte:"len:1,equal:0x20"`     // _
	Unread []byte `byte:"len:*"`                // unread
	SEP2_  []byte `byte:"len:1,equal:0x20"`     // _
	Read   []byte `byte:"len:*"`                // read
	ETX_   []byte `byte:"len:1,equal:0x03"`
}

type TplMinibarItem struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	CMD_        []byte `byte:"len:3,equal:0x4d4220"` // MB_
	Ext         []byte `byte:"len:*"`                // extension
	SEP1_       []byte `byte:"len:1,equal:0x20"`     // _
	ItemCode    []byte `byte:"len:2"`                // item code
	SEP2_       []byte `byte:"len:2,equal:0x2022"`   // _"
	Description []byte `byte:"len:*"`                // description
	SEP3_       []byte `byte:"len:2,equal:0x2220"`   // "_
	Quantity    []byte `byte:"len:2"`                // quantity
	SEP4_       []byte `byte:"len:1,equal:0x20"`     // _
	UnitPrice   []byte `byte:"len:10"`               // unit price
	SEP5_       []byte `byte:"len:1,equal:0x20"`     // _
	SubTotal    []byte `byte:"len:10"`               // sub total
	SEP6_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax1        []byte `byte:"len:10"`               // tax 1
	SEP7_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax2        []byte `byte:"len:10"`               // tax 2
	SEP8_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax3        []byte `byte:"len:10"`               // tax 3
	SEP9_       []byte `byte:"len:1,equal:0x20"`     // _
	Total       []byte `byte:"len:10"`               // total
	SEP10_      []byte `byte:"len:1,equal:0x20"`     // _
	Date        []byte `byte:"len:8"`                // yyyymmdd
	Time        []byte `byte:"len:4"`                // hhmm
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplMinibarTotal struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	CMD_     []byte `byte:"len:3,equal:0x4d4220"` // MB_
	Ext      []byte `byte:"len:*"`                // extension
	SEP1_    []byte `byte:"len:1,equal:0x20"`     // _
	SubTotal []byte `byte:"len:10"`               // sub total
	SEP6_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax1     []byte `byte:"len:10"`               // tax 1
	SEP7_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax2     []byte `byte:"len:10"`               // tax 2
	SEP8_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax3     []byte `byte:"len:10"`               // tax 3
	SEP9_    []byte `byte:"len:1,equal:0x20"`     // _
	Total    []byte `byte:"len:10"`               // total
	SEP10_   []byte `byte:"len:1,equal:0x20"`     // _
	Date     []byte `byte:"len:8"`                // yyyymmdd
	Time     []byte `byte:"len:4"`                // hhmm
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupSet4 struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext  []byte `byte:"len:*"`                // extension
	SEP_ []byte `byte:"len:1,equal:0x20"`     // _
	ACT_ []byte `byte:"len:3,equal:0x544920"` // TI_
	Time []byte `byte:"len:4"`                // hhmm
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupSet3 struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext  []byte `byte:"len:*"`                // extension
	SEP_ []byte `byte:"len:1,equal:0x20"`     // _
	ACT_ []byte `byte:"len:3,equal:0x544920"` // TI_
	Time []byte `byte:"len:3"`                // hmm
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupClear struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	State []byte `byte:"len:2,equal:0x4f46"`   // OF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupAnswer struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	State []byte `byte:"len:2"`                // State
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatus_FixedLength struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x535420"` // ST_
	Ext   []byte `byte:"len:7"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	State []byte `byte:"len:2"`                // room status RE/PR/CL/PA/FA/SK/NS
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusAdvanced_FixedLength struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x535420"` // ST_
	Ext   []byte `byte:"len:7"`                // extension
	SEP1  []byte `byte:"len:1,equal:0x20"`     // _
	State []byte `byte:"len:2"`                // room status RE/PR/CL/PA/FA/SK/NS
	SEP2  []byte `byte:"len:1,equal:0x20"`     // _
	MI_   []byte `byte:"len:3,equal:0x4d4920"` // MI_
	Maid  []byte `byte:"len:*"`                // maid
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplVoiceCount_FixedLength struct {
	STX_   []byte `byte:"len:1,equal:0x02"`
	CMD_   []byte `byte:"len:3,equal:0x495320"` // IS_
	ACT_   []byte `byte:"len:3,equal:0x564320"` // VC_
	Ext    []byte `byte:"len:7"`                // extension
	SEP1_  []byte `byte:"len:1,equal:0x20"`     // _
	Unread []byte `byte:"len:*"`                // unread
	SEP2_  []byte `byte:"len:1,equal:0x20"`     // _
	Read   []byte `byte:"len:*"`                // read
	ETX_   []byte `byte:"len:1,equal:0x03"`
}

type TplMinibarItem_FixedLength struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	CMD_        []byte `byte:"len:3,equal:0x4d4220"` // MB_
	Ext         []byte `byte:"len:7"`                // extension
	SEP1_       []byte `byte:"len:1,equal:0x20"`     // _
	ItemCode    []byte `byte:"len:2"`                // item code
	SEP2_       []byte `byte:"len:2,equal:0x2022"`   // _"
	Description []byte `byte:"len:*"`                // description
	SEP3_       []byte `byte:"len:2,equal:0x2220"`   // "_
	Quantity    []byte `byte:"len:2"`                // quantity
	SEP4_       []byte `byte:"len:1,equal:0x20"`     // _
	UnitPrice   []byte `byte:"len:10"`               // unit price
	SEP5_       []byte `byte:"len:1,equal:0x20"`     // _
	SubTotal    []byte `byte:"len:10"`               // sub total
	SEP6_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax1        []byte `byte:"len:10"`               // tax 1
	SEP7_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax2        []byte `byte:"len:10"`               // tax 2
	SEP8_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax3        []byte `byte:"len:10"`               // tax 3
	SEP9_       []byte `byte:"len:1,equal:0x20"`     // _
	Total       []byte `byte:"len:10"`               // total
	SEP10_      []byte `byte:"len:1,equal:0x20"`     // _
	Date        []byte `byte:"len:8"`                // yyyymmdd
	Time        []byte `byte:"len:4"`                // hhmm
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplMinibarTotal_FixedLength struct {
	STX_     []byte `byte:"len:1,equal:0x02"`
	CMD_     []byte `byte:"len:3,equal:0x4d4220"` // MB_
	Ext      []byte `byte:"len:7"`                // extension
	SEP1_    []byte `byte:"len:1,equal:0x20"`     // _
	SubTotal []byte `byte:"len:10"`               // sub total
	SEP6_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax1     []byte `byte:"len:10"`               // tax 1
	SEP7_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax2     []byte `byte:"len:10"`               // tax 2
	SEP8_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax3     []byte `byte:"len:10"`               // tax 3
	SEP9_    []byte `byte:"len:1,equal:0x20"`     // _
	Total    []byte `byte:"len:10"`               // total
	SEP10_   []byte `byte:"len:1,equal:0x20"`     // _
	Date     []byte `byte:"len:8"`                // yyyymmdd
	Time     []byte `byte:"len:4"`                // hhmm
	ETX_     []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupSet4_FixedLength struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext  []byte `byte:"len:7"`                // extension
	SEP_ []byte `byte:"len:1,equal:0x20"`     // _
	ACT_ []byte `byte:"len:3,equal:0x544920"` // TI_
	Time []byte `byte:"len:4"`                // hhmm
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupSet3_FixedLength struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	CMD_ []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext  []byte `byte:"len:7"`                // extension
	SEP_ []byte `byte:"len:1,equal:0x20"`     // _
	ACT_ []byte `byte:"len:3,equal:0x544920"` // TI_
	Time []byte `byte:"len:3"`                // hmm
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupClear_FixedLength struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:7"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	State []byte `byte:"len:2,equal:0x4f46"`   // OF
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupAnswer_FixedLength struct {
	STX_  []byte `byte:"len:1,equal:0x02"`
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:7"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	State []byte `byte:"len:2"`                // State
	ETX_  []byte `byte:"len:1,equal:0x03"`
}

/* ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

// incoming call packets

type TplCallPacketN struct {
	RecordType     []byte `byte:"len:1,equal:0x4e"`   // N
	SEP1_          []byte `byte:"len:1,equal:0x20"`   // _
	RecordNumber   []byte `byte:"len:3"`              //
	SEP2_          []byte `byte:"len:1,equal:0x20"`   // _
	CustomerNumber []byte `byte:"len:2"`              //
	SEP3_          []byte `byte:"len:1,equal:0x20"`   // _
	Originator     []byte `byte:"len:*"`              // extension
	SEP4_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload1       []byte `byte:"len:*"`              // payload line one
	END1_          []byte `byte:"len:2,equal:0x0d0a"` // end line one
	NEXT_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload2       []byte `byte:"len:*"`              // payload line two
	END2_          []byte `byte:"len:2,equal:0x0d0a"` // end line two
}

type TplCallPacketS struct {
	RecordType     []byte `byte:"len:1,equal:0x53"`   // S
	SEP1_          []byte `byte:"len:1,equal:0x20"`   // _
	RecordNumber   []byte `byte:"len:3"`              //
	SEP2_          []byte `byte:"len:1,equal:0x20"`   // _
	CustomerNumber []byte `byte:"len:2"`              //
	SEP3_          []byte `byte:"len:1,equal:0x20"`   // _
	Originator     []byte `byte:"len:*"`              // extension
	SEP4_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload1       []byte `byte:"len:*"`              // payload line one
	END1_          []byte `byte:"len:2,equal:0x0d0a"` // end line one
	NEXT_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload2       []byte `byte:"len:*"`              // payload line two
	END2_          []byte `byte:"len:2,equal:0x0d0a"` // end line two
}

type TplCallPacketE struct {
	RecordType     []byte `byte:"len:1,equal:0x45"`   // E
	SEP1_          []byte `byte:"len:1,equal:0x20"`   // _
	RecordNumber   []byte `byte:"len:3"`              //
	SEP2_          []byte `byte:"len:1,equal:0x20"`   // _
	CustomerNumber []byte `byte:"len:2"`              //
	SEP3_          []byte `byte:"len:1,equal:0x20"`   // _
	Originator     []byte `byte:"len:*"`              // extension
	SEP4_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload1       []byte `byte:"len:*"`              // payload line one
	END1_          []byte `byte:"len:2,equal:0x0d0a"` // end line one
	NEXT_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload2       []byte `byte:"len:*"`              // payload line two
	END2_          []byte `byte:"len:2,equal:0x0d0a"` // end line two
}

type TplCallPacketX struct {
	RecordType     []byte `byte:"len:1,equal:0x59"`   // X
	SEP1_          []byte `byte:"len:1,equal:0x20"`   // _
	RecordNumber   []byte `byte:"len:3"`              //
	SEP2_          []byte `byte:"len:1,equal:0x20"`   // _
	CustomerNumber []byte `byte:"len:2"`              //
	SEP3_          []byte `byte:"len:1,equal:0x20"`   // _
	Originator     []byte `byte:"len:*"`              // extension
	SEP4_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload1       []byte `byte:"len:*"`              // payload line one
	END1_          []byte `byte:"len:2,equal:0x0d0a"` // end line one
	NEXT_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload2       []byte `byte:"len:*"`              // payload line two
	END2_          []byte `byte:"len:2,equal:0x0d0a"` // end line two
}

type TplCallPacketL struct {
	RecordType     []byte `byte:"len:1,equal:0x4c"`   // L
	SEP1_          []byte `byte:"len:1,equal:0x20"`   // _
	RecordNumber   []byte `byte:"len:3"`              //
	SEP2_          []byte `byte:"len:1,equal:0x20"`   // _
	CustomerNumber []byte `byte:"len:2"`              //
	SEP3_          []byte `byte:"len:1,equal:0x20"`   // _
	Originator     []byte `byte:"len:*"`              // extension
	SEP4_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload1       []byte `byte:"len:*"`              // payload line one
	END1_          []byte `byte:"len:2,equal:0x0d0a"` // end line one
	NEXT_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload2       []byte `byte:"len:*"`              // payload line two
	END2_          []byte `byte:"len:2,equal:0x0d0a"` // end line two
}

// kann auch als Minibar misbraucht werden
// denn müssen die ersten drei ziffern der gewählten rufnummer mit einem pattern matchen, zb. 896
// die 4te stelle der gewählten rufnummer unterscheidet dann ob
//  - total amount (format unklar) 896_1_850 -> 8,50
//  - oder count und artikel nummer (format unklar) -> 896_2_345 -> 3*45 /34*5/1*345 man weiß es halt nicht

type TplCallPacketNSingle struct {
	RecordType     []byte `byte:"len:1,equal:0x4e"`   // N
	SEP1_          []byte `byte:"len:1,equal:0x20"`   // _
	RecordNumber   []byte `byte:"len:3"`              //
	SEP2_          []byte `byte:"len:1,equal:0x20"`   // _
	CustomerNumber []byte `byte:"len:2"`              //
	SEP3_          []byte `byte:"len:1,equal:0x20"`   // _
	Originator     []byte `byte:"len:*"`              // extension
	SEP4_          []byte `byte:"len:1,equal:0x20"`   // _
	Payload1       []byte `byte:"len:*"`              // payload line one
	END1_          []byte `byte:"len:2,equal:0x0d0a"` // end line one
}

type TplCallPacketOvr struct {
	OVR_ []byte `byte:"len:*"`
	CR_  []byte `byte:"len:2,equal:0x0d0a"`
}

/* ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- */

// incomig background terminal

type TplBgdGarbage_Overread_1a struct {
	Data_ []byte `byte:"len:*"`
	CMD_  []byte `byte:"len:2,equal:0x5354"` // ST
}

type TplBgdGarbage_Overread_1b struct {
	Data_ []byte `byte:"len:*"`
	CMD_  []byte `byte:"len:2,equal:0x4d42"` // MB
}

type TplBgdGarbage_Overread_1c struct {
	Data_ []byte `byte:"len:*"`
	CMD_  []byte `byte:"len:2,equal:0x5741"` // WA
}

type TplBgdGarbage_Overread_1d struct {
	Data_ []byte `byte:"len:*"`
	CMD_  []byte `byte:"len:2,equal:0x4953"` // IS
}

type TplBgdGarbage_Overread_2 struct {
	Data_ []byte `byte:"len:*"`
	CMD_  []byte `byte:"len:2,equal:0x0d0a"` // CR LF
}

type TplBgdErrorMnemonic struct {
	Msg  []byte `byte:"len:8,equal:0x4d4e454d4f4e4943"` // MNEMONIC
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`             // CR LF
}

type TplBgdErrorInput struct {
	Msg  []byte `byte:"len:11,equal:0x494e505554204552524f52"` // INPUT ERROR
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`                    // CR LF
}

type TplBgdErrorTryAgain struct {
	Msg  []byte `byte:"len:9,equal:0x54525920414741494e"` // TRY AGAIN
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`               // CR LF
}

type TplBgdErrorNameBig struct {
	Msg  []byte `byte:"len:8,equal:0x4e414d4520424947"` // NAME BIG
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`             // CR LF
}

type TplBgdErrorDuplicate struct {
	Msg  []byte `byte:"len:9,equal:0x4455504c4943415445"` // DUPLICATE
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`               // CR LF
}

type TplBgdErrorNoData struct {
	Msg  []byte `byte:"len:13,equal:0x4e4f204441544120464f554e44"` // NO DATA FOUND
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`                        // CR LF
}

type TplBgdErrorNoSetCPNDData struct {
	Msg  []byte `byte:"len:16,equal:0x4e4f205345542043504e442044415441"` // NO SET CPND DATA
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`                              // CR LF
}

type TplBgdErrorNoCustCPNDData struct {
	Msg  []byte `byte:"len:17,equal:0x4e4f20435553542043504e442044415441"` // NO CUST CPND DATA
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`                                // CR LF
}

type TplBgdErrorNoCPNDMemory struct {
	Msg  []byte `byte:"len:14,equal:0x4e4f2043504e44204d454d4f5259"` // NO CPND MEMORY
	CMD_ []byte `byte:"len:2,equal:0x0d0a"`                          // CR LF
}

type TplBgdRoomStatus struct {
	CMD_  []byte `byte:"len:3,equal:0x535420"` // ST_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	State []byte `byte:"len:2"`                // room status RE/PR/CL/PA/FA/SK/NS
	CRLF_ []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdRoomStatusAdvanced struct {
	CMD_  []byte `byte:"len:3,equal:0x535420"` // ST_
	Ext   []byte `byte:"len:*"`                // extension
	SEP1  []byte `byte:"len:1,equal:0x20"`     // _
	State []byte `byte:"len:2"`                // room status RE/PR/CL/PA/FA/SK/NS
	SEP2  []byte `byte:"len:1,equal:0x20"`     // _
	MI_   []byte `byte:"len:3,equal:0x4d4920"` // MI_
	Maid  []byte `byte:"len:*"`                // maid
	CRLF_ []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdVoiceCount struct {
	CMD_   []byte `byte:"len:3,equal:0x495320"` // IS_
	ACT_   []byte `byte:"len:3,equal:0x564320"` // VC_
	Ext    []byte `byte:"len:*"`                // extension
	SEP1_  []byte `byte:"len:1,equal:0x20"`     // _
	Unread []byte `byte:"len:*"`                // unread
	SEP2_  []byte `byte:"len:1,equal:0x20"`     // _
	Read   []byte `byte:"len:*"`                // read
	CRLF_  []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdMinibarItem struct {
	CMD_        []byte `byte:"len:3,equal:0x4d4220"` // MB_
	Ext         []byte `byte:"len:*"`                // extension
	SEP1_       []byte `byte:"len:1,equal:0x20"`     // _
	ItemCode    []byte `byte:"len:2"`                // item code
	SEP2_       []byte `byte:"len:2,equal:0x2022"`   // _"
	Description []byte `byte:"len:*"`                // description
	SEP3_       []byte `byte:"len:2,equal:0x2220"`   // "_
	Quantity    []byte `byte:"len:2"`                // quantity
	SEP4_       []byte `byte:"len:1,equal:0x20"`     // _
	UnitPrice   []byte `byte:"len:10"`               // unit price
	SEP5_       []byte `byte:"len:1,equal:0x20"`     // _
	SubTotal    []byte `byte:"len:10"`               // sub total
	SEP6_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax1        []byte `byte:"len:10"`               // tax 1
	SEP7_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax2        []byte `byte:"len:10"`               // tax 2
	SEP8_       []byte `byte:"len:1,equal:0x20"`     // _
	Tax3        []byte `byte:"len:10"`               // tax 3
	SEP9_       []byte `byte:"len:1,equal:0x20"`     // _
	Total       []byte `byte:"len:10"`               // total
	SEP10_      []byte `byte:"len:1,equal:0x20"`     // _
	Date        []byte `byte:"len:8"`                // yyyymmdd
	Time        []byte `byte:"len:4"`                // hhmm
	CRLF_       []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdMinibarTotal struct {
	CMD_     []byte `byte:"len:3,equal:0x4d4220"` // MB_
	Ext      []byte `byte:"len:*"`                // extension
	SEP1_    []byte `byte:"len:1,equal:0x20"`     // _
	SubTotal []byte `byte:"len:10"`               // sub total
	SEP6_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax1     []byte `byte:"len:10"`               // tax 1
	SEP7_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax2     []byte `byte:"len:10"`               // tax 2
	SEP8_    []byte `byte:"len:1,equal:0x20"`     // _
	Tax3     []byte `byte:"len:10"`               // tax 3
	SEP9_    []byte `byte:"len:1,equal:0x20"`     // _
	Total    []byte `byte:"len:10"`               // total
	SEP10_   []byte `byte:"len:1,equal:0x20"`     // _
	Date     []byte `byte:"len:8"`                // yyyymmdd
	Time     []byte `byte:"len:4"`                // hhmm
	CRLF_    []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdWakeupSet4 struct {
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	Time  []byte `byte:"len:4"`                // hhmm
	CRLF_ []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdWakeupSet3 struct {
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	Time  []byte `byte:"len:3"`                // hmm
	CRLF_ []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdWakeupClear struct {
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	State []byte `byte:"len:2,equal:0x4f46"`   // OF
	CRLF_ []byte `byte:"len:2,equal:0x0d0a"`
}

type TplBgdWakeupAnswer struct {
	CMD_  []byte `byte:"len:3,equal:0x574120"` // WA_
	Ext   []byte `byte:"len:*"`                // extension
	SEP_  []byte `byte:"len:1,equal:0x20"`     // _
	ACT_  []byte `byte:"len:3,equal:0x544920"` // TI_
	State []byte `byte:"len:2"`                // State
	CRLF_ []byte `byte:"len:2,equal:0x0d0a"`
}
