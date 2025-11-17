package otrum

type TplVerifyPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x56455220"` // VER
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplErrorPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x45525220"` // ERR
	ErrorCode   []byte `byte:"len:2"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplBalancePacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x42414c20"` // BAL
	Amount      []byte `byte:"len:8"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplCheckinPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x43484b49"` // CHKI
	RoomNumber  []byte `byte:"len:6"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplCheckoutPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x43484b4f"` // CHKO
	RoomNumber  []byte `byte:"len:6"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplDisplayRequest struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x44495350"` // DISP
	RoomNumber    []byte `byte:"len:6"`                  //
	AccountNumber []byte `byte:"len:6"`                  // Reservation
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplRoomStatusPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x48534b50"` // HSKP
	RoomNumber  []byte `byte:"len:6"`
	StatusCode  []byte `byte:"len:2"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplInfoPacket struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Transaction    []byte `byte:"len:4"`
	Sequence       int    `byte:"len:4"`
	Command_       []byte `byte:"len:4,equal:0x494e464f"` // INFO
	RoomNumber     []byte `byte:"len:6"`                  //
	AccountNumber  []byte `byte:"len:6"`                  // Reservation
	GuestName      []byte `byte:"len:20"`                 //
	MessageWaiting []byte `byte:"len:1"`                  // Y = unread messages, N = no unread messages
	GroupName      []byte `byte:"len:5"`                  //
	BillView       []byte `byte:"len:1"`                  // Y = allowed, N = not allowed, " " = default
	Checkout       []byte `byte:"len:1"`                  // Y = express cechkout allowed, N = not allowed, " " = default
	Language       []byte `byte:"len:1"`                  //
	Welcome        []byte `byte:"len:1"`                  // Y = welcome picture, N = dont show, " " = default
	Rights         []byte `byte:"len:3"`                  //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

/*
type TplInfoPacket2 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Transaction    []byte `byte:"len:4"`
	Sequence       int    `byte:"len:4"`
	Command_       []byte `byte:"len:4,equal:0x494e464f"` // INFO
	RoomNumber     []byte `byte:"len:6"`                  //
	AccountNumber  []byte `byte:"len:6"`                  // Reservation
	GuestName      []byte `byte:"len:20"`                 //
	MessageWaiting []byte `byte:"len:1"`                  // Y = unread messages, N = no unread messages
	GroupName      []byte `byte:"len:5"`                  //
	BillView       []byte `byte:"len:1"`                  // Y = allowed, N = not allowed, " " = default
	Checkout       []byte `byte:"len:1"`                  // Y = express cechkout allowed, N = not allowed, " " = default
	Language       []byte `byte:"len:1"`                  //
	Welcome        []byte `byte:"len:1"`                  // Y = welcome picture, N = dont show, " " = default
	Rights         []byte `byte:"len:3"`                  //
	ProfileID      []byte `byte:"len:3"`                  //
	ETX_           []byte `byte:"len:1,equal:0x03"`
}
*/

type TplInitPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x494e4954"` // INIT
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplItemPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x4954454d"` // ITEM
	Date        []byte `byte:"len:4"`                  // MMDD
	Description []byte `byte:"len:12"`                 //
	Indicator   []byte `byte:"len:2"`                  // "  " = Charge, "CR" = Payment
	Amount      []byte `byte:"len:7"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplLookupRequest struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x4c4f4f4b"` // LOOK
	RoomNumber  []byte `byte:"len:6"`                  //
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageCallerPacket struct {
	STX_              []byte `byte:"len:1,equal:0x02"`
	Transaction       []byte `byte:"len:4"`
	Sequence          int    `byte:"len:4"`
	Command_          []byte `byte:"len:4,equal:0x4d434c52"` // MCLR
	MessageNumber     []byte `byte:"len:6"`                  //
	CallerName        []byte `byte:"len:24"`                 //
	CallerLocation    []byte `byte:"len:24"`                 //
	CallerPhoneNumber []byte `byte:"len:24"`                 //
	ETX_              []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageHeaderPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x4d484452"` // MHDR
	AccountNumber []byte `byte:"len:6"`                  //
	MessageNumber []byte `byte:"len:6"`                  //
	Date          []byte `byte:"len:6"`                  // MMDDYY
	Time          []byte `byte:"len:6"`                  // HHMMSS
	ReceiverName  []byte `byte:"len:24"`                 //
	ByPerson      []byte `byte:"len:1"`                  // Y/N
	ByPhone       []byte `byte:"len:1"`                  // Y/N
	PleaseCall    []byte `byte:"len:1"`                  // Y/N
	Callback      []byte `byte:"len:1"`                  // Y/N
	ReturnedCall  []byte `byte:"len:1"`                  // Y/N
	Urgent        []byte `byte:"len:1"`                  // Y/N
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageReadPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x4d534744"` // MSGD
	AccountNumber []byte `byte:"len:6"`                  //
	MessageNumber []byte `byte:"len:6"`                  //
	Date          []byte `byte:"len:6"`                  // MMDDYY
	Time          []byte `byte:"len:6"`                  // HHMMSS
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageRequestPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x4d534752"` // MSGR
	RoomNumber    []byte `byte:"len:6"`
	AccountNumber []byte `byte:"len:6"` //
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageStatusPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x4d534757"` // MSGW
	RoomNumber  []byte `byte:"len:6"`                  //
	StatusCode  []byte `byte:"len:1"`                  // Y = unread messages, N = no unread messages
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageTextPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x4d545854"` // MTXT
	MessageNumber []byte `byte:"len:6"`                  //
	MessageText   []byte `byte:"len:64"`                 //
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type TplNamePacket struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	Transaction    []byte `byte:"len:4"`
	Sequence       int    `byte:"len:4"`
	Command_       []byte `byte:"len:4,equal:0x4e414d45"` // NAME
	RoomNumber     []byte `byte:"len:6"`                  //
	AccountNumber  []byte `byte:"len:6"`                  // Reservation
	GuestName      []byte `byte:"len:20"`                 //
	MessageWaiting []byte `byte:"len:1"`                  // Y = unread messages, N = no unread messages
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type TplPostChargePacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x504f5354"` // POST
	RoomNumber  []byte `byte:"len:6"`
	RevenueCode []byte `byte:"len:2"`
	Description []byte `byte:"len:12"`
	Amount      []byte `byte:"len:7"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplStatusRequest struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x53544154"` // STAT
	RoomNumber  []byte `byte:"len:6"`                  //
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplStartPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x53545254"` //STRT
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplTestPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x54455354"` //TEST
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupResult struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x574b4445"` // WKDE
	RoomNumber  []byte `byte:"len:6"`                  //
	Time        []byte `byte:"len:4"`                  // HHMM
	Date        []byte `byte:"len:6"`                  //MMDDYY
	StatusCode  []byte `byte:"len:1"`                  // D = success, T = timeout, H = hardware error, TripleGuest U = undeliverable
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupSetPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x574b4f44"` // WKOD
	RoomNumber  []byte `byte:"len:6"`                  //
	Time        []byte `byte:"len:4"`                  // HHMM
	Date        []byte `byte:"len:6"`                  // MMDDYY
	Order_      []byte `byte:"len:1,equal:0x4f"`       // O = order
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupClearPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x574b4f44"` // WKOD
	RoomNumber  []byte `byte:"len:6"`                  //
	Time        []byte `byte:"len:4"`                  // HHMM
	Date        []byte `byte:"len:6"`                  // MMDDYY
	Order_      []byte `byte:"len:1,equal:0x43"`       // C = cancel
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplCheckoutRequest struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x58434b4f"` // XCKO
	RoomNumber    []byte `byte:"len:6"`                  //
	AccountNumber []byte `byte:"len:6"`                  // Reservation
	BalanceAmount []byte `byte:"len:8"`                  //
	ETX_          []byte `byte:"len:1,equal:0x03"`
}
