package ascii

/*
process code
i	1-6 associated feature access code was dialed from RSN
o	8 pms reject message
o	9 pms accept message
*/

// outgoing

type HouseKeeperRoomAccepted struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x21"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b52"` // HKR
	PROC_         []byte `byte:"len:1,equal:0x09"`     // process code 9
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type HouseKeeperRoomRejected struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x21"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b52"` // HKR
	PROC_         []byte `byte:"len:1,equal:0x08"`     // process code 8
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type HouseKeeperRoom5 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x21"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b52"` // HKR
	PROC          []byte `byte:"len:1"`                // process code
	RSN           []byte `byte:"len:5"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type HouseKeeperRoom7 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x21"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b52"` // HKR
	PROC          []byte `byte:"len:1"`                // process code
	RSN           []byte `byte:"len:5"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

/*
process code
i	1-4 associated feature access code was dialed from station
o	8 pms reject message
o	9 pms accept message
*/

// outgoing

type HouseKeeperStationAccepted struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x22"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b53"` // HKS
	PROC_         []byte `byte:"len:1,equal:0x09"`     // process code 9
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type HouseKeeperStationRejected struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x22"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b53"` // HKS
	PROC_         []byte `byte:"len:1,equal:0x08"`     // process code 8
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type HouseKeeperStation5 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x22"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b53"` // HKS
	PROC          []byte `byte:"len:1"`                // process code 8/9
	RSN           []byte `byte:"len:5"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type HouseKeeperStation7 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x22"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x484b53"` // HKS
	PROC          []byte `byte:"len:1"`                // process code 8/9
	RSN           []byte `byte:"len:7"`                // room station number (5/7)
	DIG           []byte `byte:"len:6"`                // dig
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}
