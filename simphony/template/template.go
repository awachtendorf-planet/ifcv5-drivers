package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketNonFrontOfficeCharge       = "Non FrontOffice Charge"
	PacketChargePostingSelectedGuest = "Charge Posting Selected Guest"
	PacketPostingInquiry             = "Posting Inquiry"
	PacketCheckFacsimileRequest      = "Check Facsimile Request"

	PacketInquiryResponse        = "Inquiry Response"
	PacketChargePostingAck       = "Charge Posting Acknowledge"
	PacketCheckFacsimileResponse = "Check Facsimile Response"

	PacketPing = "Extra Ping"

	PacketGarbage         = "Garbage"
	PacketGarbageLowLevel = "Garbage Low Level"
	PacketUnknown         = "Unknown framed packet"

	PacketHandshake = "Handshake"
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

type TplGarbage_Framing_2 struct {
	STX_OVR_ []byte `byte:"len:1,equal:0x01"`
	Data_    []byte `byte:"len:*"`
	STX_     []byte `byte:"len:1,equal:0x01"`
}

type TplGarbage_Framing_3 struct {
	Data_ []byte `byte:"len:*"`
	SOH_  []byte `byte:"len:1,equal:0x01"`
}

type Tcp_Ping struct {
	SOH_ []byte `byte:"len:1,equal:0x01"`
	ID   []byte `byte:"len:*"`
	STX_ []byte `byte:"len:1,equal:0x02"`
	ETX_ []byte `byte:"len:1,equal:0x03"`
	EOT_ []byte `byte:"len:1,equal:0x04"`
}

// INCOMING
type TplSimphonyHandshakePacket struct { // POS => PMS
	SOH_     []byte `byte:"len:1,equal:0x01"`
	SourceID []byte `byte:"len:*"`
	STX_     []byte `byte:"len:1,equal:0x02"`
	FS1_     []byte `byte:"len:1,equal:0x1c"`
	PmsText  []byte `byte:"len:*"`
	FS2_     []byte `byte:"len:1,equal:0x1c"`
	FS3_     []byte `byte:"len:1,equal:0x1c"`
	Extra2   []byte `byte:"len:1,equal:0x31"`
	ETX_     []byte `byte:"len:1,equal:0x03"`
	Checksum []byte `byte:"len:*"`
	EOT_     []byte `byte:"len:1,equal:0x04"`
}

type TplSimphonyChargePostingNonFrontOfficeChargePacket struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                                 /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID             []byte `byte:"len:11,equal:0x4946435f4348475f505354"` // IFC_CHG_PST
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	Space1_               []byte `byte:"len:1,equal:0x20"`
	FS3_                  []byte `byte:"len:1,equal:0x1c"`
	Space2_               []byte `byte:"len:1,equal:0x20"`
	FS4_                  []byte `byte:"len:1,equal:0x1c"`
	PaymentType           []byte `byte:"len:*"`
	FS6_                  []byte `byte:"len:1,equal:0x1c"`
	TenderAmount          []byte `byte:"len:*"`
	FS7_                  []byte `byte:"len:1,equal:0x1c"`
	NumSalesItemizer      []byte `byte:"len:*"`
	FS8_                  []byte `byte:"len:1,equal:0x1c"`
	Data                  []byte `byte:"len:*"` // starts and ends without FS!
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplSimphonyChargePostingSelectedGuestPacket struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                                 /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID_            []byte `byte:"len:11,equal:0x4946435f4348475f505354"` // IFC_CHG_PST
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	GuestID               []byte `byte:"len:*"`
	FS3_                  []byte `byte:"len:1,equal:0x1c"`
	GuestName             []byte `byte:"len:*"`
	FS4_                  []byte `byte:"len:1,equal:0x1c"`
	SelectionInitiator_   []byte `byte:"len:1,equal:0x53"`
	SelectionNumber       []byte `byte:"len:*"`
	FS5_                  []byte `byte:"len:1,equal:0x1c"`
	PaymentType           []byte `byte:"len:*"`
	FS6_                  []byte `byte:"len:1,equal:0x1c"`
	TenderAmount          []byte `byte:"len:*"`
	FS7_                  []byte `byte:"len:1,equal:0x1c"`
	NumSalesItemizer      []byte `byte:"len:*"`
	FS8_                  []byte `byte:"len:1,equal:0x1c"`
	Data                  []byte `byte:"len:*"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplSimphonyPostingInquiryPacket struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                                 /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID_            []byte `byte:"len:11,equal:0x4946435f4348475f505354"` // IFC_CHG_PST
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	GuestID               []byte `byte:"len:*"`
	FS3_                  []byte `byte:"len:1,equal:0x1c"`
	PaymentType           []byte `byte:"len:*"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplSimphonyCheckFascimileRequestPacket struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                                              /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID_            []byte `byte:"len:*,equal:0x464F5F46495F46494E414C5F54454E444552"` // FO_FI_FINAL_TENDER
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	Date_                 []byte `byte:"len:*"`
	FS3_                  []byte `byte:"len:1,equal:0x1c"`
	Time_                 []byte `byte:"len:*"`
	FS4_                  []byte `byte:"len:1,equal:0x1c"`
	RevenueCenterNumber_  []byte `byte:"len:*"`
	FS5_                  []byte `byte:"len:1,equal:0x1c"`
	GuestCheckNumber_     []byte `byte:"len:*"`
	FS6_                  []byte `byte:"len:1,equal:0x1c"`
	Data                  []byte `byte:"len:*"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

// OUTGOING
// TODO: header übernehmen und zurücksenden
// type TplSimphonyInquiryResponsePacket struct { // PMS => POS
// 	SOH_                  []byte `byte:"len:1,equal:0x01"`
// 	SourceID              []byte `byte:"len:*"`
// 	STX_                  []byte `byte:"len:1,equal:0x02"`
// 	FS1_                  []byte `byte:"len:1,equal:0x1c"`
// 	SequenceNumber        []byte `byte:"len:2"`
// 	MessageRetransmitFlag []byte `byte:"len:1"`                                 /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
// 	MessageID_            []byte `byte:"len:11,equal:0x4946435f4753545f53454c"` // IFC_GST_SEL
// 	FS2_                  []byte `byte:"len:1,equal:0x1c"`
// 	GuestID               []byte `byte:"len:*"`
// 	FS3_                  []byte `byte:"len:1,equal:0x1c"`
// 	GuestName             []byte `byte:"len:*"`
// 	FS4_                  []byte `byte:"len:1,equal:0x1c"`
// 	ListSize              []byte `byte:"len:*"`
// 	FS5_                  []byte `byte:"len:1,equal:0x1c"`
// 	GuestList             []byte `byte:"len:*"`
// 	ETX_                  []byte `byte:"len:1,equal:0x03"`
// 	Checksum              []byte `byte:"len:*"`
// 	EOT_                  []byte `byte:"len:1,equal:0x04"`
// }

type TplSimphonyInquiryResponsePacket struct { // PMS => POS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                                 /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID_            []byte `byte:"len:11,equal:0x4946435f4753545f53454c"` // IFC_GST_SEL
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	GuestID               []byte `byte:"len:*"`
	FS3_                  []byte `byte:"len:1,equal:0x1c"`
	ChargePosting_        []byte `byte:"len:11,equal:0x4946435f4348475f505354"` // IFC_CHG_PST
	FS4_                  []byte `byte:"len:1,equal:0x1c"`
	ListSize              []byte `byte:"len:*"`
	FS5_                  []byte `byte:"len:1,equal:0x1c"`
	GuestList             []byte `byte:"len:*"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplSimphonyChargePostingAckPacket struct { // PMS => POS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                                 /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID_            []byte `byte:"len:11,equal:0x4946435f4348475f505354"` // IFC_CHG_PST
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	MessageStatus         []byte `byte:"len:1"`
	FS3_                  []byte `byte:"len:1,equal:0x1c"`
	Status                []byte `byte:"len:1"`
	FS4_                  []byte `byte:"len:1,equal:0x1c"`
	Message               []byte `byte:"len:30"`
	FS5_                  []byte `byte:"len:1,equal:0x1c"`
	Tndttl_               []byte `byte:"len:12,equal:0x202020202020202020202020"` // zb 1
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplSimphonyCheckFascimileResponsePacket struct { // IFC => POS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	SourceID              []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	FS1_                  []byte `byte:"len:1,equal:0x1c"`
	SequenceNumber        []byte `byte:"len:2"`
	MessageRetransmitFlag []byte `byte:"len:1"`                               /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	MessageID_            []byte `byte:"len:10,equal:0x464F5F38375F444F4E45"` // FO_87_DONE
	FS2_                  []byte `byte:"len:1,equal:0x1c"`
	Text                  []byte `byte:"len:70"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}
