package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketHeartbeat      = "Heartbeat"
	PacketHeartbeatReply = "Heartbeat Reply"

	PacketSyncRequest = "Sync Request"
	PacketSyncStart   = "Sync Start"
	PacketSyncEnd     = "Sync End"

	PacketLinkEndRequest   = "Link End Request"
	PacketLinkEndConfirmed = "Link End Confirmed"

	PacketCheckIn                    = "Check In"
	PacketCheckOut                   = "Check Out"
	PacketDataChange                 = "Data Change"
	PacketSetRestriction             = "Set Restriction"
	PacketSwitchMessageLamp          = "Switch Message Lamp"
	PacketMessageLampOn              = "Message Lamp On"
	PacketMessageLampOff             = "Message Lamp Off"
	PacketHouseKeeperRoom            = "House Keeper Status From Room"
	PacketHouseKeeperStation         = "House Keeper Status From Station"
	PacketHouseKeeperRoomAccepted    = "House Keeper Status From Room Accepted"
	PacketHouseKeeperStationAccepted = "House Keeper Status From Station Accepted"
	PacketHouseKeeperRoomRejected    = "House Keeper Status From Room Rejected"
	PacketHouseKeeperStationRejected = "House Keeper Status From Station Rejected"

	PacketRoomDataImageSwap = "Room Data Image"

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
