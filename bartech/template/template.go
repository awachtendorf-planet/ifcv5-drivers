package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketCheckIn                 = "Check In"
	PacketCheckInSimple           = "Check In (simple)"
	PacketCheckInWithoutHappyHour = "Check In (without happy hour)"
	PacketCheckOut                = "Check Out"
	PacketLockBar                 = "Lock Bar"
	PacketUnlockBar               = "Unlock Bar"
	PacketEndOfDay                = "End Of Day"
	PacketSyncRequest             = "Sync Request"
	PacketRoomStatus              = "Room Status"
	PacketInvoice                 = "Invoice"

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

// some fucking "bei Opera machen wir das aber so"

type TplCheckinSimple struct {
	STX_      []byte `byte:"len:1,equal:0x02"` //
	Cmd       []byte `byte:"len:2"`            // 32 = unlocked, 35 = locked
	Extension []byte `byte:"len:*"`            // 4 or 6 stellig
	ETX_      []byte `byte:"len:1,equal:0x03"` //
}

type TplCheckinWithoutHappyHour struct {
	STX_      []byte `byte:"len:1,equal:0x02"` //
	Cmd       []byte `byte:"len:2"`            // 32 = unlocked, 35 = locked
	Extension []byte `byte:"len:*"`            // 4 or 6 stellig
	GuestName []byte `byte:"len:32"`           //
	Password  []byte `byte:"len:6"`            // numeric
	Opt1      []byte `byte:"len:1"`            // Y/N Safe used
	Opt2      []byte `byte:"len:1"`            // Y/N
	Opt3      []byte `byte:"len:1"`            // Y/N
	Opt4      []byte `byte:"len:1"`            // Y/N
	CODate    []byte `byte:"len:8"`            // mmddccyy
	ETX_      []byte `byte:"len:1,equal:0x03"` //
}

// official protocol
type TplCheckin struct {
	STX_      []byte `byte:"len:1,equal:0x02"` //
	Cmd       []byte `byte:"len:2"`            // 32 = unlocked, 35 = locked
	Extension []byte `byte:"len:*"`            // 4 or 6 stellig
	HappyHour []byte `byte:"len:2"`            // 00 = no, 01 = free, 02 = first comsumption free, 03 = first day free
	GuestName []byte `byte:"len:30"`           //
	Password  []byte `byte:"len:6"`            // numeric
	Opt1      []byte `byte:"len:1"`            // Y/N Safe used
	Opt2      []byte `byte:"len:1"`            // Y/N
	Opt3      []byte `byte:"len:1"`            // Y/N
	Opt4      []byte `byte:"len:1"`            // Y/N
	CODate    []byte `byte:"len:8"`            // mmddccyy
	ETX_      []byte `byte:"len:1,equal:0x03"` //
}

type TplCheckout struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3331"` // 31
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplLockBar struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3333"` // 33
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplUnlockBar struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	Cmd_      []byte `byte:"len:2,equal:0x3334"` // 34
	Extension []byte `byte:"len:*"`
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplEndOfDay struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ []byte `byte:"len:2,equal:0x3336"` // 36
	Opt  []byte `byte:"len:*"`              // filled with "0"
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

// incomming

type TplSyncBasic4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3134"` // 14
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
	OVR6_     []byte `byte:"len:28"`             // ovr pos 26-53
	OVR7_     []byte `byte:"len:1"`              // ovr pos 54
	Ticket    []byte `byte:"len:4"`              // line number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplSyncAdvanced4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3135"` // 15
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:8"`              // mmddccyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 21
	Time      []byte `byte:"len:6"`              // hhmmss
	OVR5_     []byte `byte:"len:42"`             // ovr pos 28-69
	Ticket    []byte `byte:"len:5"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplSyncBasic6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3134"` // 14
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
	OVR6_     []byte `byte:"len:28"`             // ovr pos 26-53
	OVR7_     []byte `byte:"len:1"`              // ovr pos 54
	Ticket    []byte `byte:"len:4"`              // line number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplSyncAdvanced6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3135"` // 15
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:8"`              // mmddccyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 21
	Time      []byte `byte:"len:6"`              // hhmmss
	OVR5_     []byte `byte:"len:42"`             // ovr pos 28-69
	Ticket    []byte `byte:"len:5"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusBasic4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3939"` // 99
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
	Status    []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:27"`             // ovr pos 28-54
	Ticket    []byte `byte:"len:4"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusAdvanced4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3938"` // 98
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:8"`              // mmddccyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 21
	Time      []byte `byte:"len:6"`              // hhmmss
	OVR5_     []byte `byte:"len:1"`              // ovr pos 28
	Status    []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:39"`             // ovr pos 31-69
	Ticket    []byte `byte:"len:5"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusBasic6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3939"` // 99
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
	Status    []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:27"`             // ovr pos 28-54
	Ticket    []byte `byte:"len:4"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusAdvanced6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3938"` // 98
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:8"`              // mmddccyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 21
	Time      []byte `byte:"len:6"`              // hhmmss
	OVR5_     []byte `byte:"len:1"`              // ovr pos 28
	Status    []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:39"`             // ovr pos 31-69
	Ticket    []byte `byte:"len:5"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplInvoiceBasic4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3030"` // 00
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
	TaxID     []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:2"`              // ovr pos 28,29
	Quantity  []byte `byte:"len:1"`              //
	OVR7_     []byte `byte:"len:1"`              // ovr pos 31
	Text      []byte `byte:"len:15"`             //
	OVR8_     []byte `byte:"len:1"`              // ovr pos 47
	Charge    []byte `byte:"len:6"`              //
	OVR9_     []byte `byte:"len:1"`              // ovr pos 54
	Ticket    []byte `byte:"len:4"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplInvoiceAdvanced4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3031"` // 01
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:8"`              // mmddccyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 21
	Time      []byte `byte:"len:6"`              // hhmmss
	OVR5_     []byte `byte:"len:1"`              // ovr pos 28
	TaxID     []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:1"`              // ovr pos 31
	Text      []byte `byte:"len:15"`             //
	OVR7_     []byte `byte:"len:1"`              // ovr pos 47
	Charge    []byte `byte:"len:10"`             //
	OVR8_     []byte `byte:"len:1"`              // ovr pos 58
	TaxAmount []byte `byte:"len:10"`             //
	OVR9_     []byte `byte:"len:1"`              // ovr pos 69
	Ticket    []byte `byte:"len:5"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplInvoiceBasic6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3030"` // 00
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
	TaxID     []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:2"`              // ovr pos 28,29
	Quantity  []byte `byte:"len:1"`              //
	OVR7_     []byte `byte:"len:1"`              // ovr pos 31
	Text      []byte `byte:"len:15"`             //
	OVR8_     []byte `byte:"len:1"`              // ovr pos 47
	Charge    []byte `byte:"len:6"`              //
	OVR9_     []byte `byte:"len:1"`              // ovr pos 54
	Ticket    []byte `byte:"len:4"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplInvoiceAdvanced6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3031"` // 01
	OVR2_     []byte `byte:"len:1"`              // ovr pos 5
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1"`              // ovr pos 12
	Date      []byte `byte:"len:8"`              // mmddccyy
	OVR4_     []byte `byte:"len:1"`              // ovr pos 21
	Time      []byte `byte:"len:6"`              // hhmmss
	OVR5_     []byte `byte:"len:1"`              // ovr pos 28
	TaxID     []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:1"`              // ovr pos 31
	Text      []byte `byte:"len:15"`             //
	OVR7_     []byte `byte:"len:1"`              // ovr pos 47
	Charge    []byte `byte:"len:10"`             //
	OVR8_     []byte `byte:"len:1"`              // ovr pos 58
	TaxAmount []byte `byte:"len:10"`             //
	OVR9_     []byte `byte:"len:1"`              // ovr pos 69
	Ticket    []byte `byte:"len:5"`              // ticket number
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

// type TplInvoiceBasic struct {
// 	STX_      []byte `byte:"len:1,equal:0x02"`
// 	OVR1_     []byte `byte:"len:2"`
// 	Cmd_      []byte `byte:"len:2,equal:0x3030"` // 00
// 	OVR2_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 5, blank
// 	Extension []byte `byte:"len:*"`              // extension undefined length
// 	OVR3_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 12, blank
// 	Date      []byte `byte:"len:6"`              // ddmmyy
// 	OVR4_     []byte `byte:"len:1"`              // ovr pos 19
// 	Time      []byte `byte:"len:4"`              // hhmm
// 	OVR5_     []byte `byte:"len:2"`              // ovr pos 24,25
// 	TaxID     []byte `byte:"len:2"`              //
// 	OVR6_     []byte `byte:"len:2"`              // ovr pos 28,29
// 	Quantity  []byte `byte:"len:1"`              //
// 	OVR7_     []byte `byte:"len:1"`              // ovr pos 31
// 	Text      []byte `byte:"len:15"`             //
// 	OVR8_     []byte `byte:"len:1"`              // ovr pos 47
// 	Charge    []byte `byte:"len:6"`              //
// 	OVR9_     []byte `byte:"len:1"`              // ovr pos 54
// 	Ticket    []byte `byte:"len:4"`              // ticket number
// 	OVR10_    []byte `byte:"len:*"`              // ovr bartech bullshit
// 	ETX_      []byte `byte:"len:1,equal:0x03"`
// }

type TplInvoiceUndocumented1_4 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3030"` // 00
	OVR2_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 5, blank
	Extension []byte `byte:"len:4"`              //
	OVR3_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 12, blank
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 19, blank
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2,equal:0x2020"` // ovr pos 24,25, blank
	TaxID     []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:2,equal:0x2020"` // ovr pos 28,29, blank
	Quantity  []byte `byte:"len:1"`              //
	OVR7_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 31, blank
	Text      []byte `byte:"len:15"`             //
	OVR8_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 47, blank
	Charge    []byte `byte:"len:6"`              //
	OVR9_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 54, blank
	Ticket    []byte `byte:"len:4"`              // ticket number
	OVR10_    []byte `byte:"len:*"`              // ovr bartech bullshit
	ETX_      []byte `byte:"len:1,equal:0x03"`
}

type TplInvoiceUndocumented1_6 struct {
	STX_      []byte `byte:"len:1,equal:0x02"`
	OVR1_     []byte `byte:"len:2"`
	Cmd_      []byte `byte:"len:2,equal:0x3030"` // 00
	OVR2_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 5, blank
	Extension []byte `byte:"len:6"`              //
	OVR3_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 12, blank
	Date      []byte `byte:"len:6"`              // ddmmyy
	OVR4_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 19, blank
	Time      []byte `byte:"len:4"`              // hhmm
	OVR5_     []byte `byte:"len:2,equal:0x2020"` // ovr pos 24,25, blank
	TaxID     []byte `byte:"len:2"`              //
	OVR6_     []byte `byte:"len:2,equal:0x2020"` // ovr pos 28,29, blank
	Quantity  []byte `byte:"len:1"`              //
	OVR7_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 31, blank
	Text      []byte `byte:"len:15"`             //
	OVR8_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 47, blank
	Charge    []byte `byte:"len:6"`              //
	OVR9_     []byte `byte:"len:1,equal:0x20"`   // ovr pos 54, blank
	Ticket    []byte `byte:"len:4"`              // ticket number
	OVR10_    []byte `byte:"len:*"`              // ovr bartech bullshit
	ETX_      []byte `byte:"len:1,equal:0x03"`
}
