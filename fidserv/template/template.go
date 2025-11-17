package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketLinkStart       = "Link Start"
	PacketLinkAlive       = "Link Alive"
	PacketLinkEnd         = "Link End"
	PacketLinkDescription = "Link Description"
	PacketLinkRecord      = "Link Record"

	PacketResyncRequest = "Database Resync Request"
	PacketResyncStart   = "Database Resync Start"
	PacketResyncEnd     = "Database Resync End"

	PacketCheckIn      = "Check In"
	PacketCheckInSwap  = "Check In (Swap)"
	PacketCheckOut     = "Check Out"
	PacketCheckOutSwap = "Check Out (Swap)"
	PacketDataChange   = "Data Change"

	PacketRoomData = "Room Equipment Status"

	PacketNightAuditStart = "Night Audit Start"
	PacketNightAuditEnd   = "Night Audit End"

	PacketKeyRequest = "Key Request"
	PacketKeyDelete  = "Key Delete"
	PacketKeyChange  = "Key Change"
	PacketKeyRead    = "Key Read"
	PacketKeyAnswer  = "Key Answer"

	PacketWakeupRequest = "Wakeup Request"
	PacketWakeupClear   = "Wakeup Clear"
	PacketWakeupAnswer  = "Wakeup Answer"

	PacketPostingSimple  = "Posting Simple"
	PacketPostingRequest = "Posting Request"
	PacketPostingAnswer  = "Posting Answer"
	PacketPostingList    = "Posting List"

	PacketGuestMessageOnline  = "Guest Message Online"
	PacketGuestMessageRequest = "Guest Message Request"
	PacketGuestMessageText    = "Guest Message Text"
	PacketGuestMessageDelete  = "Guest Message Delete"

	PacketGuestBillRequest = "Guest Bill Request"
	PacketGuestBillItem    = "Guest Bill Item"
	PacketGuestBillBalance = "Guest Bill Balance"

	PacketRemoteCheckOut    = "Remote Check Out"
	PacketGuestCheckDetails = "Guest Check Details"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"
)

type TplLinkStart struct {
	STX_ uint8  `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4c53"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplLinkAlive struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4c41"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplLinkEnd struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4c45"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplLinkDescription struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4c44"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplLinkRecord struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4c52"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplDBRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4452"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplDBStart struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4453"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplDBEnd struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4445"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCheckIn struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4749"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCheckOut struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x474F"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplDataChange struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4743"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplRoomData struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5245"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5752"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupClear struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5743"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplWakeupAnswer struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5741"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplNightAuditStart struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4e53"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplNightAuditEnd struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4e45"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplKeyRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4b52"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplKeyDelete struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4b44"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplKeyChange struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4b4d"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplKeyRead struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4b5a"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplKeyAnswer struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4b41"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCheckInSwap struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x4749"`
	Data []byte `byte:"len:*"`
	SF_  []byte `byte:"len:3,equal:0x53467C"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplCheckOutSwap struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x474F"`
	Data []byte `byte:"len:*"`
	SF_  []byte `byte:"len:3,equal:0x53467C"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplPostingSimple struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5053"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplPostingRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5052"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplPostingAnswer struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5041"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplPostingList struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x504c"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageOnline struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x584c"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x584d"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageText struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5854"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestMessageDelete struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5844"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestBillRequest struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5852"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestBillItem struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5849"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestBillBalance struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5842"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplRemoteCheckOut struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x5843"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
}

type TplGuestCheckDetails struct {
	STX_ []byte `byte:"len:1,equal:0x02"`
	Cmd_ uint8  `byte:"len:2,equal:0x434B"`
	Data []byte `byte:"len:*"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
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
