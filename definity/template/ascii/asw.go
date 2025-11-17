package ascii

/*
process code
o	1 switch is perform the function
i	2 confirmation, no action was taken because RSN already occupied
i	3 confirmation, DID in name field if available
*/

// outgoing

type CheckIn struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x26"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x434b49"` // CKI
	PROC_         []byte `byte:"len:1,equal:0x01"`     // process code 1
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	COVERAGE_PATH []byte `byte:"len:4"`                // cover path (can also take null values)
	DISPLAY_NAME  []byte `byte:"len:*"`                // 15/30
	VM_PASSWORD   []byte `byte:"len:4"`                // leave blank for default
	VM_LANGUAGE   []byte `byte:"len:2"`                //
	REQ_DID       []byte `byte:"len:1"`                // y to request DID, otherwise NULL
	NULL_         []byte `byte:"len:3,equal:0x202020"` // NULL, NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type CheckInReply5 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x26"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x434b49"` // CKI
	PROC          []byte `byte:"len:1"`                // process code 1
	RSN           []byte `byte:"len:5"`                // room station number (5/7)
	COVERAGE_PATH []byte `byte:"len:4"`                // cover path (can also take null values)
	DISPLAY_NAME  []byte `byte:"len:15"`               // 15/30
	VM_PASSWORD   []byte `byte:"len:4"`                // leave blank for default
	VM_LANGUAGE   []byte `byte:"len:2"`                //
	REQ_DID       []byte `byte:"len:1"`                // y to request DID, otherwise NULL
	NULL_         []byte `byte:"len:3"`                // NULL, NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type CheckInReply7 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x26"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x434b49"` // CKI
	PROC          []byte `byte:"len:1"`                // process code 1
	RSN           []byte `byte:"len:7"`                // room station number (5/7)
	COVERAGE_PATH []byte `byte:"len:4"`                // cover path (can also take null values)
	DISPLAY_NAME  []byte `byte:"len:30"`               // 15/30
	VM_PASSWORD   []byte `byte:"len:4"`                // leave blank for default
	VM_LANGUAGE   []byte `byte:"len:2"`                //
	REQ_DID       []byte `byte:"len:1"`                // y to request DID, otherwise NULL
	NULL_         []byte `byte:"len:3"`                // NULL, NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

/*
process code
o	1 switch is perform the function
i	2 completed as requested
i	3 no action taken, RSN vacant
i 	4 completed, no action taken because information the same as already stored
*/

// outgoing

type DataChange struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x28"`
	MSGCT         []byte `byte:"len:1"`                  // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x474943"`   // GIC
	PROC_         []byte `byte:"len:1,equal:0x01"`       // process code 1
	RSN           []byte `byte:"len:*"`                  // room station number (5/7)
	COVERAGE_PATH []byte `byte:"len:4"`                  // cover path (can also take null values)
	DISPLAY_NAME  []byte `byte:"len:*"`                  // 15/30
	VM_PASSWORD   []byte `byte:"len:4"`                  // leave blank for default
	VM_LANGUAGE   []byte `byte:"len:2"`                  //
	NULL_         []byte `byte:"len:4,equal:0x20202020"` // NULL, NULL, NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type DataChangeReply5 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x28"`
	MSGCT         []byte `byte:"len:1"`                  // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x474943"`   // GIC
	PROC          []byte `byte:"len:1"`                  // process code
	RSN           []byte `byte:"len:5"`                  // room station number (5/7)
	COVERAGE_PATH []byte `byte:"len:4"`                  // cover path (can also take null values)
	DISPLAY_NAME  []byte `byte:"len:15"`                 // 15/30
	VM_PASSWORD   []byte `byte:"len:4"`                  // leave blank for default
	VM_LANGUAGE   []byte `byte:"len:2"`                  //
	NULL_         []byte `byte:"len:4,equal:0x20202020"` // NULL, NULL, NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type DataChangeReply7 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x28"`
	MSGCT         []byte `byte:"len:1"`                  // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x474943"`   // GIC
	PROC          []byte `byte:"len:1"`                  // process code
	RSN           []byte `byte:"len:7"`                  // room station number (5/7)
	COVERAGE_PATH []byte `byte:"len:4"`                  // cover path (can also take null values)
	DISPLAY_NAME  []byte `byte:"len:30"`                 // 15/30
	VM_PASSWORD   []byte `byte:"len:4"`                  // leave blank for default
	VM_LANGUAGE   []byte `byte:"len:2"`                  //
	NULL_         []byte `byte:"len:4,equal:0x20202020"` // NULL, NULL, NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

/*
process code
o	1 switch is perform the function
i	2 completed, message lamp was off
i	3 completed, message lamp was on
i	4 completed, no action was taken
i 	5 completed, message lamp still on
*/

// outgoing

type CheckOut struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x29"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x434b4f"` // CKO
	PROC_         []byte `byte:"len:1,equal:0x01"`     // process code 1
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type CheckOutReply struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x29"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x434b4f"` // CKO
	PROC          []byte `byte:"len:1"`                // process code
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}
