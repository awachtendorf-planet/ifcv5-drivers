package ascii

/*
process code
o	F heart beat
i	0 heart beat reply
i	1 initiate db
i	2 initiate db
o	3 start db
o	4 end db
i	5 release requested
o	6 release confirmed
*/

// type TplStatusInquiry struct {
// 	STX_          []byte `byte:"len:1,equal:0x02"`
// 	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
// 	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
// 	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
// 	PROC          []byte `byte:"len:1"`                // process code
// 	RR            []byte `byte:"len:1"`                // release reason
// 	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
// 	ETX_          []byte `byte:"len:1,equal:0x03"`
// }

// outgoing

type Heartbeat struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x46"`     // process code F
	RR_           []byte `byte:"len:1,equal:0x46"`     // release reason F
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type LinkEndReply struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x36"`     // process code 6
	RR_           []byte `byte:"len:1,equal:0x46"`     // release reason F
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type SyncStart struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x33"`     // process code 3
	RR_           []byte `byte:"len:1,equal:0x46"`     // release reason F
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type SyncEnd struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x34"`     // process code 4
	RR_           []byte `byte:"len:1,equal:0x46"`     // release reason F
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type HeartbeatReply struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x30"`     // process code 0
	RR            []byte `byte:"len:1"`                // release reason
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type SyncRequest1 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x31"`     // process code 1
	RR            []byte `byte:"len:1"`                // release reason
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type SyncRequest2 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x32"`     // process code 2
	RR            []byte `byte:"len:1"`                // release reason
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type LinkEndRequest struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x71"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x535453"` // STS
	PROC_         []byte `byte:"len:1,equal:0x35"`     // process code 5
	RR            []byte `byte:"len:1"`                // release reason
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}
