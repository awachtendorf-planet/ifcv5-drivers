package sonifi

// outgoing

// type TplCheckinPacket struct {
// 	STX_        []byte `byte:"len:1,equal:0x02"`
// 	Transaction []byte `byte:"len:4"`
// 	Sequence    int    `byte:"len:4"`
// 	Command_    []byte `byte:"len:4,equal:0x43484b43"` // CHKC
// 	RoomNumber  []byte `byte:"len:6"`
// 	ETX_        []byte `byte:"len:1,equal:0x03"`
// }

// type TplCheckinEnhancedPacket struct {
// 	STX_           []byte `byte:"len:1,equal:0x02"`
// 	Transaction    []byte `byte:"len:4"`
// 	Sequence       int    `byte:"len:4"`
// 	Command_       []byte `byte:"len:4,equal:0x43484b43"` // CHKC
// 	RoomNumber     []byte `byte:"len:6"`
// 	GroupID        []byte `byte:"len:16"`
// 	AffinityNumber []byte `byte:"len:16"`
// 	ETX_           []byte `byte:"len:1,equal:0x03"`
// }

// type TplCheckinEnhancedPacket struct {
// 	STX_           []byte `byte:"len:1,equal:0x02"`
// 	Transaction    []byte `byte:"len:4"`
// 	Sequence       int    `byte:"len:4"`
// 	Command_       []byte `byte:"len:4,equal:0x43484b49"` // CHKI
// 	RoomNumber     []byte `byte:"len:6"`
// 	GroupID        []byte `byte:"len:16"`
// 	AffinityNumber []byte `byte:"len:16"`
// 	ETX_           []byte `byte:"len:1,equal:0x03"`
// }

// type TplNameSonifiPacket struct {
// 	STX_          []byte `byte:"len:1,equal:0x02"`
// 	Transaction   []byte `byte:"len:4"`
// 	Sequence      int    `byte:"len:4"`
// 	Command_      []byte `byte:"len:4,equal:0x4e414d45"` // NAME
// 	RoomNumber    []byte `byte:"len:6"`                  //
// 	AccountNumber []byte `byte:"len:6"`                  // Reservation
// 	GuestName     []byte `byte:"len:20"`                 //
// 	Unused_       []byte `byte:"len:1, equal:0x20"`      // Unused -> bei otrum ist das MessageWaiting
// 	FolioOption_  []byte `byte:"len:1, eq1ual:0x30"`     // optional -> deshalb kompatible to otrum
// 	ETX_          []byte `byte:"len:1,equal:0x03"`
// }

type TplNameHeloPacket struct {
	STX_              []byte `byte:"len:1,equal:0x02"`
	Transaction       []byte `byte:"len:4"`
	Sequence          int    `byte:"len:4"`
	Command_          []byte `byte:"len:4,equal:0x48454c4f"` // HELO
	ID                []byte `byte:"len:8"`                  // PMS ID (string)
	Version           []byte `byte:"len:4"`                  // PMS Version (int)
	STAT              []byte `byte:"len:1"`                  // Y = STAT/INFO, N = LOOK/NAME
	FUDP              []byte `byte:"len:1"`                  // Y = FUPD/FCLR, N = DISP
	ROST              []byte `byte:"len:1"`                  // Y / N
	MSGS              []byte `byte:"len:1"`                  // Y / N
	DGMH              []byte `byte:"len:1"`                  // Y / N
	MSGD              []byte `byte:"len:1"`                  // Y / N
	PayTVChannel      []byte `byte:"len:1"`                  // Y / N , set to N
	POST              []byte `byte:"len:1"`                  // Y / N
	HSKP              []byte `byte:"len:1"`                  // Y / N
	VariableMsgLength []byte `byte:"len:1"`                  // Y / N, set to N
	LocalFolioNumber  []byte `byte:"len:1"`                  // Y / N, set to Y
	PairedFolioNumber []byte `byte:"len:1"`                  // Y / N, set to ?
	DecimalPosition   []byte `byte:"len:1"`                  // 2
	CharacterSet      []byte `byte:"len:8"`                  // ISO88591
	FutureSpace       []byte `byte:"len:32"`                 //
	ETX_              []byte `byte:"len:1,equal:0x03"`
}

// incoming

type TplPostChargeEnhancedPacket struct {
	STX_        []byte `byte:"len:1,equal:0x02"`
	Transaction []byte `byte:"len:4"`
	Sequence    int    `byte:"len:4"`
	Command_    []byte `byte:"len:4,equal:0x504f5354"` // POST
	RoomNumber  []byte `byte:"len:6"`
	RevenueCode []byte `byte:"len:2"`
	Description []byte `byte:"len:12"`
	PurchaseID  []byte `byte:"len:10"`
	Amount      []byte `byte:"len:7"`
	ETX_        []byte `byte:"len:1,equal:0x03"`
}

type TplCheckoutEnhancedPacket struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	Transaction   []byte `byte:"len:4"`
	Sequence      int    `byte:"len:4"`
	Command_      []byte `byte:"len:4,equal:0x58434b4f"` // XCKO
	RoomNumber    []byte `byte:"len:6"`                  //
	AccountNumber []byte `byte:"len:6"`                  // Reservation
	BalanceAmount []byte `byte:"len:8"`                  //
	FolioFlag     []byte `byte:"len:1"`                  //
	ETX_          []byte `byte:"len:1,equal:0x03"`
}
