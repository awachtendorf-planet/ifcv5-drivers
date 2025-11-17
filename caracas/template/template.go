package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"
	PacketSyn = "SYN"
	PacketEot = "EOT"

	PacketBind     = "Bind"
	PacketDatetime = "Set Datetime"

	PacketCheckIn            = "Check In"
	PacketDataChangeAdvanced = "Advanced DataChange"
	PacketCheckOut           = "Check Out"
	PacketUpdateCheckInName  = "Update Check In Name"

	PacketNoPost = "No Posting"

	PacketMessageLight            = "Message Light"
	PacketMessageLightAlternative = "Message Light Alternative"

	PacketWakeupOrder         = "Wakeup Order"
	PacketWakeupOrderIncoming = "Wakeup Order Incoming"
	PacketWakeupResult        = "Wakeup Result"

	PacketServiceCosts = "Service Costs"

	PacketCallCharge            = "Call Charge"
	PacketServiceCharge         = "Service Charge"
	PacketMinibarCharge         = "Minibar Charge Simple"
	PacketMinibarChargeAdvanced = "Minibar Charge Advanced"
	PacketMinibarChargeComplete = "Minibar Charge Complete"

	PacketVoicemail          = "Voicemail"
	PacketVoicemailAdvanced  = "Voicemail Advanced"
	PacketRoomstatus         = "Roomstatus"
	PacketRoomstatusAdvanced = "Roomstatus Advanced"
	PacketPresence           = "Presence"
	PacketDND                = "Do Not Disturb"
	PacketClassOfService     = "Class of Service"

	PacketDBSyncRequest = "DB Sync Request"
	PacketServerStatus  = "Server Status"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"
)

type TplEOT struct {
	EOT_ []byte `byte:"len:1,equal:0x04"`
}

type TplACK struct {
	ACK_ []byte `byte:"len:1,equal:0x06"`
}

type TplNAK struct {
	NAK_ []byte `byte:"len:1,equal:0x15"`
}

type TplENQ struct {
	ENQ_ []byte `byte:"len:1,equal:0x05"`
}

type TplSYN struct {
	SYN_ []byte `byte:"len:1,equal:0x16"`
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

// Outgoing

type CheckInPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3031"`                  // 01
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Permissions_       []byte `byte:"len:2,equal:0x3030"`                  // 00 = vendor default
	CheckIn_           []byte `byte:"len:1,equal:0x31"`                    // 1 = CheckIn
	Language           []byte `byte:"len:2"`                               // 01-09
	GuestName          []byte `byte:"len:*"`                               // max len. 40
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type AdvancedDataChangePacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3039"`                  // 09
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Permissions_       []byte `byte:"len:2,equal:0x3030"`                  // 00 = vendor default
	Order              []byte `byte:"len:1"`                               // 1 = CheckIn, 2 Sharer CheckIn, 3 Move to Empty Room, 4 Move to NOT Empty Room, 5 Checkout Room not Empty, 6 Checkout Room Empty
	Language           []byte `byte:"len:2"`                               // 01-09
	GuestName          []byte `byte:"len:40"`                              //
	SyncInd            []byte `byte:"len:1"`                               // 0 new CheckIn, 1 DBResync
	ReservationID      []byte `byte:"len:10"`                              // left-justified filled with blanks
	Group              []byte `byte:"len:10"`                              // left-justified filled with blanks
	OldRoomNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
	VIP                []byte `byte:"len:3"`                               // left-justified filled with blanks
	GroupName          []byte `byte:"len:40"`                              // left-justified filled with blanks
	CreditLimit        []byte `byte:"len:10"`                              // Leading zeroes
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type CheckOutPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3031"`                  // 01
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Permissions_       []byte `byte:"len:2,equal:0x3030"`                  // 00 = vendor default
	CheckOut_          []byte `byte:"len:1,equal:0x32"`                    // 2 = CheckOut
	Language           []byte `byte:"len:2"`                               // 01-09
	GuestName          []byte `byte:"len:*"`                               // max len. 40
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type MessageLightPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3032"`                  // 02
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	MessageLightStatus []byte `byte:"len:1"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type AlternativeMessageLightPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3330"`                  // 30
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1q
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Mode               []byte `byte:"len:1"`                               // 1 set in Caracas, 2 set in Voicemail System, 3 set in Carcas and VoiceMail
	MessageLightStatus []byte `byte:"len:1"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type WakeupOrderPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3033"`                  // 03
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	DateTime           []byte `byte:"len:10"`                              // ddmmyyhhmm
	Order              []byte `byte:"len:1"`                               // 0 remove order, 1 set order, 2 remove all
	Mode               []byte `byte:"len:1"`                               // 0 one shot, 1 repeat daily
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

// type SetDisplayPacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3034"`                  // 04
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	Status             []byte `byte:"len:1"`                               // 0 silent, 1 with sound
// 	Text               []byte `byte:"len:*"`                               // max len. 40
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// type RequestWakeupOrderPacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3035"`                  // 05
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	BlockWakeups       []byte `byte:"len:1"`                               // 0 set wakeup allowed, 1 set wakeup blocked
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// type EnableDNDPacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3037"`                  // 07
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	DND                []byte `byte:"len:1"`                               // 0 DND blocked, 1 DND enabled
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// type InterruptCallPacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3038"`                  // 08
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	BlockSubsToo       []byte `byte:"len:1"`                               // 0 Extension, 1 Extension + Subextensions
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// type BlockMessageTypesPacket struct { // Christine - evtl NoPost
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3043"`                  // 0C
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension_         []byte `byte:"len:6,equal:0x303030303030"`          // always 000000
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	MessageType        []byte `byte:"len:2"`                               //
// 	Status             []byte `byte:"len:1"`                               // 0 Disable, Enable
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

type DNDBRKPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3231"`                  // 21
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Status             []byte `byte:"len:1"`                               // 2 international, 3 national, 4 local, 5 blocked, 0 cancel DND, 1 DND
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type UpdateCheckInNamePacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3233"`                  // 23
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	GuestName          []byte `byte:"len:40"`                              // left-justified filled with blanks
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

// type ConfigureMinibarArticlePacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3234"`                  // 24
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	ArticleNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	Status             []byte `byte:"len:1"`                               // 0 remove article, 1 add article, 2 remove article file
// 	ArticleName        []byte `byte:"len:20"`                              // left-justified filled with blanks
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// Incoming

type CallChargePacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3131"`                  // 11
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	SystemID           []byte `byte:"len:2"`                               // 50,54,55,70,74,75 8818; 20 Hicom 150E/200; 30 Hicom 300; 81 Voicemail System
	Filler_            []byte `byte:"len:2,equal:0x3130"`                  // always 10
	CallDate           []byte `byte:"len:6"`                               // ddmmyy
	CallStartTime      []byte `byte:"len:6"`                               // hhmmss
	CallEndTime        []byte `byte:"len:6"`                               // hhmmss
	PBXSubDeviceID     []byte `byte:"len:1"`                               // always 0 for Hicom
	PBXServiceID       []byte `byte:"len:1"`                               // always 0 for Hicom
	GroupID            []byte `byte:"len:2"`                               // GroupID/Extension Type 00-99
	CallExtension      []byte `byte:"len:6"`                               // left-justified filled with blanks
	ExchangeLineNumber []byte `byte:"len:4"`                               // left-justified filled with blanks
	Units              []byte `byte:"len:5"`                               // Units/Pulses
	TargetNumber       []byte `byte:"len:16"`                              // left-justified filled with blanks
	TargetNumberID     []byte `byte:"len:1"`                               //
	PBXCallType        []byte `byte:"len:1"`                               // always 0 for Hicom
	CodeNumber         []byte `byte:"len:12"`                              // if used, otherwise always 000000000000
	CallTypeHicom      []byte `byte:"len:1"`                               // always 0 for Hicom
	CallTotalAmount    []byte `byte:"len:10"`                              // leading zeroes
	Costs              []byte `byte:"len:10"`                              // leading zeroes (for instance costs per unit)
	CallType           []byte `byte:"len:1"`                               // National, International, Local, Mobile
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

// type IncomingDNDPacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3232"`                  // 22
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	DND                []byte `byte:"len:1"`                               // 0 DND OFF, 1 DND ON
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

type VoiceMessagePacket struct { // Christine
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3132"`                  // 12
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Status             []byte `byte:"len:1"`                               // 0 Voicemail tapped, 1 new Voicemail
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type AdvancedVoiceMessagePacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3239"`                  // 29
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	AmountNewMessages  []byte `byte:"len:2"`                               // sets Message Light if != 0
	AmountOldMessages  []byte `byte:"len:2"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type RoomstatusPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3133"`                  // 13
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Roomstatus         []byte `byte:"len:2"`                               // 01 Clean, 1 Dirty, 03-09...
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type AdvancedRoomstatusPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3139"`                  // 19
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Roomstatus         []byte `byte:"len:2"`                               // 01 Clean, 1 Dirty, 03-09...
	UserID             []byte `byte:"len:10"`                              // left-justified filled with blanks
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type IncomingWakeupOrderPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3134"`                  // 14
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Time               []byte `byte:"len:4"`                               // hhmm
	Status             []byte `byte:"len:1"`                               // 0 remove order, 1 set order
	Mode               []byte `byte:"len:1"`                               // 0 one shot, 1 repeat daily
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type WakeupAttemptPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3137"`                  // 17
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	DateTime           []byte `byte:"len:10"`                              // ddmmyyhhmm
	Status             []byte `byte:"len:1"`                               // 0 not successful, 1 successful
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type MinibarPostingPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3135"`                  // 15
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	ArticleNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
	Status             []byte `byte:"len:1"`                               // 0 booked, 1 cancelled (Caracas just uses '1')
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type AdvancedMinibarPostingPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3138"`                  // 18
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	ArticleNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
	Status             []byte `byte:"len:1"`                               // 0 booked, 1 cancelled (Caracas just uses '1')
	UserID             []byte `byte:"len:10"`                              // left-justified filled with blanks
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type CompleteMinibarPostingPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3332"`                  // 32
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	ArticleNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
	ArticleName        []byte `byte:"len:30"`                              // left-justified filled with blanks
	ArticleAmount      []byte `byte:"len:2"`                               // leading zero
	ArticleCosts       []byte `byte:"len:10"`                              // leading zeroes EEEEEEE,CC
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

// type AdditionalChargingsPacket struct { // Sauna, Pool and so on
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3136"`                  // 16
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks (Guest)
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	IssuerExtension    []byte `byte:"len:6"`                               // left-justified filled with blanks (Invoice Issuer)
// 	ArticleNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// type PresencePacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3230"`                  // 20
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	Status             []byte `byte:"len:1"`                               // 0 sign out, 1 sign in
// 	UserID             []byte `byte:"len:10"`                              // left-justified filled with blanks
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

type DBResyncRequestPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3235"`                  // 25
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension_         []byte `byte:"len:6,equal:0x303030303030"`          // always 000000
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Status             []byte `byte:"len:1"`                               // 1 CI, 2 Guestnames, 3 Articles, 4 Wakeups, 5 DND, 6 Messages, 0 All
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

// type ServerStatusPacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3236"`                  // 26
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension_         []byte `byte:"len:6,equal:0x303030303030"`          // always 000000
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	SystemID           []byte `byte:"len:2"`                               // 21 Hicom 150E/200; 31 Hicom 300; 51 or 71 reserved; 22 Hicom 150E Office; 61 Caracas Server; 81 Voicemail; 91 Callstar Horizon
// 	ErrorCode          []byte `byte:"len:3"`                               // ann (a application, nn error code)
// 	LogMessage         []byte `byte:"len:60"`                              // System Log Message
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// type ServiceChargePacket struct {
// 	STX_               []byte `byte:"len:1,equal:0x02"`
// 	Cmd_               []byte `byte:"len:2,equal:0x3331"`                  // 31
// 	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
// 	Extension          []byte `byte:"len:6"`                               // left-justified filled with blanks (Guest)
// 	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
// 	IssuerExtension    []byte `byte:"len:6"`                               // left-justified filled with blanks (Invoice Issuer)
// 	ArticleNumber      []byte `byte:"len:6"`                               // left-justified filled with blanks
// 	ArticleCosts       []byte `byte:"len:10"`                              // leading zeroes
// 	ETX_               []byte `byte:"len:1,equal:0x03"`
// }

// Both

type BindPacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3041"`                  // 0A
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension_         []byte `byte:"len:6,equal:0x303030303030"`          // always 000000
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	Status             []byte `byte:"len:1"`                               // 0 not ready, 1 ready
	ETX_               []byte `byte:"len:1,equal:0x03"`
}

type DatetimePacket struct {
	STX_               []byte `byte:"len:1,equal:0x02"`
	Cmd_               []byte `byte:"len:2,equal:0x3042"`                  // 0B
	SubSystemID_       []byte `byte:"len:1,equal:0x31"`                    // always 1
	Extension_         []byte `byte:"len:6,equal:0x303030303030"`          // always 000000
	ApplicationHandle_ []byte `byte:"len:10,equal:0x30303030303030303030"` // always 0000000000
	SystemID_          []byte `byte:"len:2,equal:0x3731"`                  // 21 Hicom150E/200, 31 Hicom 300, 51 or 71 reserved
	DateTime           []byte `byte:"len:12"`                              // ddmmyyhhmmss
	DayOfWeek          []byte `byte:"len:1"`                               // 1 Monday, 2 Tuesday...
	ETX_               []byte `byte:"len:1,equal:0x03"`
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
