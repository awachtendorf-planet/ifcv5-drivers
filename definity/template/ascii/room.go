package ascii

/*
process code
o	1 informational purpose only
i	2 response to process code 1
o	3 database update
i	4 response to process code 3
*/

/*
controlled restiction level code
00 no restiction
01 outward restriction
02 station to station restriction
03 outward and station to station restriction
04 total restriction
05 termination restiction, denies all calls to the room
06 outward and termination restiction
07 station to station and termination restiction
*/

// outgoing

type RoomDataImageSwap struct {
	STX_            []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_   []byte `byte:"len:1,equal:0x27"`
	MSGCT           []byte `byte:"len:1"`                  // message counter x(f)
	FRAME_          []byte `byte:"len:3,equal:0x524d49"`   // RMI
	PROC_           []byte `byte:"len:1,equal:0x03"`       // process code 3
	RSN             []byte `byte:"len:*"`                  // room station number (5/7)
	OCCUPIED        []byte `byte:"len:1"`                  // 0=vacant, 1=occupied
	MESSAGE_WAITING []byte `byte:"len:1"`                  // message waiting 0=off, 1=on
	RESTRICT_LEVEL  []byte `byte:"len:2"`                  // controlled restriction level code
	COVERAGE_PATH   []byte `byte:"len:4"`                  // cover path (can also take null values)
	DISPLAY_NAME    []byte `byte:"len:*"`                  // 15/30
	VM_PASSWORD     []byte `byte:"len:4"`                  // leave blank for default
	VM_LANGUAGE     []byte `byte:"len:2"`                  //
	NULL_           []byte `byte:"len:4,equal:0x20202020"` // NULL, NULL, NULL, NULL
	ETX_            []byte `byte:"len:1,equal:0x03"`
}

// incoming

type RoomDataImageSwapReply5 struct {
	STX_            []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_   []byte `byte:"len:1,equal:0x27"`
	MSGCT           []byte `byte:"len:1"`                // message counter x(f)
	FRAME_          []byte `byte:"len:3,equal:0x524d49"` // RMI
	PROC_           []byte `byte:"len:1,equal:0x04"`     // process code 4
	RSN             []byte `byte:"len:5"`                // room station number (5/7)
	OCCUPIED        []byte `byte:"len:1"`                // 0=vacant, 1=occupied
	MESSAGE_WAITING []byte `byte:"len:1"`                // message waiting 0=off, 1=on
	RESTRICT_LEVEL  []byte `byte:"len:2"`                // controlled restriction level code
	COVERAGE_PATH   []byte `byte:"len:4"`                // cover path (can also take null values)
	DISPLAY_NAME    []byte `byte:"len:15"`               // 15/30
	VM_PASSWORD     []byte `byte:"len:4"`                // leave blank for default
	VM_LANGUAGE     []byte `byte:"len:2"`                //
	NULL_           []byte `byte:"len:4"`
	ETX_            []byte `byte:"len:1,equal:0x03"`
}

type RoomDataImageSwapReply7 struct {
	STX_            []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_   []byte `byte:"len:1,equal:0x27"`
	MSGCT           []byte `byte:"len:1"`                // message counter x(f)
	FRAME_          []byte `byte:"len:3,equal:0x524d49"` // RMI
	PROC_           []byte `byte:"len:1,equal:0x04"`     // process code 4
	RSN             []byte `byte:"len:7"`                // room station number (5/7)
	OCCUPIED        []byte `byte:"len:1"`                // 0=vacant, 1=occupied
	MESSAGE_WAITING []byte `byte:"len:1"`                // message waiting 0=off, 1=on
	RESTRICT_LEVEL  []byte `byte:"len:2"`                // controlled restriction level code
	COVERAGE_PATH   []byte `byte:"len:4"`                // cover path (can also take null values)
	DISPLAY_NAME    []byte `byte:"len:30"`               // 15/30
	VM_PASSWORD     []byte `byte:"len:4"`                // leave blank for default
	VM_LANGUAGE     []byte `byte:"len:2"`                //
	NULL_           []byte `byte:"len:4"`
	ETX_            []byte `byte:"len:1,equal:0x03"`
}

/*
process code
o	1 turn on message lamp
o	2 turn off message lamp
i	3 message lamp has been turned on
i	4 message lamp has been turned off
i 	5 message lamp was already on or is still on, reply to 1 or 2 ?
*/

// outgoing

type MessageLamp struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x23"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x4d5347"` // MSG
	PROC          []byte `byte:"len:1"`                // process code
	RSN           []byte `byte:"len:*"`                // room station number (5/7)
	AR_           []byte `byte:"len:1,equal:0x20"`     // add/rem
	TF_           []byte `byte:"len:1,equal:0x20"`     // text/fax
	NULL_         []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

// incoming

type MessageLampOn5 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x23"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x4d5347"` // MSG
	PROC_         []byte `byte:"len:1,equal:0x03"`     // process code 3
	RSN           []byte `byte:"len:5"`                // room station number (5/7)
	AR_           []byte `byte:"len:1"`                // add/rem
	TF_           []byte `byte:"len:1"`                // text/fax
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type MessageLampOn7 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x23"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x4d5347"` // MSG
	PROC_         []byte `byte:"len:1,equal:0x03"`     // process code 3
	RSN           []byte `byte:"len:7"`                // room station number (5/7)
	AR_           []byte `byte:"len:1"`                // add/rem
	TF_           []byte `byte:"len:1"`                // text/fax
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type MessageLampOff5 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x23"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x4d5347"` // MSG
	PROC_         []byte `byte:"len:1,equal:0x04"`     // process code 4
	RSN           []byte `byte:"len:5"`                // room station number (5/7)
	AR_           []byte `byte:"len:1"`                // add/rem
	TF_           []byte `byte:"len:1"`                // text/fax
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

type MessageLampOff7 struct {
	STX_          []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_ []byte `byte:"len:1,equal:0x23"`
	MSGCT         []byte `byte:"len:1"`                // message counter x(f)
	FRAME_        []byte `byte:"len:3,equal:0x4d5347"` // MSG
	PROC_         []byte `byte:"len:1,equal:0x04"`     // process code 4
	RSN           []byte `byte:"len:7"`                // room station number (5/7)
	AR_           []byte `byte:"len:1"`                // add/rem
	TF_           []byte `byte:"len:1"`                // text/fax
	NULL_         []byte `byte:"len:2"`                // NULL, NULL
	ETX_          []byte `byte:"len:1,equal:0x03"`
}

/*
process code
o	1 set restriction
i	2 restriction set
*/

// outgoing

type SetRestriction struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_  []byte `byte:"len:1,equal:0x25"`
	MSGCT          []byte `byte:"len:1"`                // message counter x(f)
	FRAME_         []byte `byte:"len:3,equal:0x435220"` // CR
	PROC_          []byte `byte:"len:1,equal:0x01"`     // process code 1
	RSN            []byte `byte:"len:*"`                // room station number (5/7)
	RESTRICT_LEVEL []byte `byte:"len:2"`                // controlled restriction level code
	NULL_          []byte `byte:"len:2,equal:0x2020"`   // NULL, NULL
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

// incoming

type SetRestriction5 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_  []byte `byte:"len:1,equal:0x25"`
	MSGCT          []byte `byte:"len:1"`                // message counter x(f)
	FRAME_         []byte `byte:"len:3,equal:0x435220"` // CR
	PROC_          []byte `byte:"len:1,equal:0x02"`     // process code 2
	RSN            []byte `byte:"len:5"`                // room station number (5/7)
	RESTRICT_LEVEL []byte `byte:"len:2"`                // controlled restriction level code
	NULL_          []byte `byte:"len:2"`                // NULL, NULL
	ETX_           []byte `byte:"len:1,equal:0x03"`
}

type SetRestriction7 struct {
	STX_           []byte `byte:"len:1,equal:0x02"`
	FEATURE_CODE_  []byte `byte:"len:1,equal:0x25"`
	MSGCT          []byte `byte:"len:1"`                // message counter x(f)
	FRAME_         []byte `byte:"len:3,equal:0x435220"` // CR
	PROC_          []byte `byte:"len:1,equal:0x02"`     // process code 2
	RSN            []byte `byte:"len:7"`                // room station number (5/7)
	RESTRICT_LEVEL []byte `byte:"len:2"`                // controlled restriction level code
	NULL_          []byte `byte:"len:2"`                // NULL, NULL
	ETX_           []byte `byte:"len:1,equal:0x03"`
}
