package template

// https://github.com/bastengao/bytesparser

const (
	PacketAck = "ACK"
	PacketNak = "NAK"
	PacketEnq = "ENQ"

	PacketInquiryRequest      = "Computer Inquire Request Message"
	PacketOutletChargeRequest = "Outlet Charge Posting Request Message"
	PacketCFReqM              = "Check Facsimile Request Message"

	PacketPing = "Extra Ping"

	PacketInquiryResponse      = "Computer Inquire Response Message"
	PacketOutletChargeResponse = "Outlet Charge Posting Response Message"
	PacketCFRspM               = "Check Facsimile Response Message"

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

type TplInquiryRequest16Packet struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	WorkStation           []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	MessageType_          []byte `byte:"len:2,equal:0x2031"`
	MessageRetransmitFlag []byte `byte:"len:1"` /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID             []byte `byte:"len:16"`
	EmployeeNumber        []byte `byte:"len:4"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplInquiryRequest19Packet struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	WorkStation           []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	MessageType_          []byte `byte:"len:2,equal:0x2031"`
	MessageRetransmitFlag []byte `byte:"len:1"` /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID             []byte `byte:"len:19"`
	EmployeeNumber        []byte `byte:"len:4"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplInquiryResponsePacket struct { // PMS => POS
	SOH_                []byte `byte:"len:1,equal:0x01"`
	WorkStation         []byte `byte:"len:*"`
	STX_                []byte `byte:"len:1,equal:0x02"`
	MessageType         []byte `byte:"len:2"`
	InformationMessages []byte `byte:"len:*"`
	ETX_                []byte `byte:"len:1,equal:0x03"`
	Checksum            []byte `byte:"len:*"`
	EOT_                []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting48M19Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}
type TplOutletChargePosting48M19kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting48M16Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting48M16kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting44M19Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting44M19kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting44M16kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting44M16Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting8M19Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting8M19kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting8M16Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting8M16kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting16M19Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	Sales9Total                 []byte `byte:"len:10"` // 3700 = 0
	Sales10Total                []byte `byte:"len:10"` // 3700 = 0
	Sales11Total                []byte `byte:"len:10"` // 3700 = 0
	Sales12Total                []byte `byte:"len:10"` // 3700 = 0
	Sales13Total                []byte `byte:"len:10"` // 3700 = 0
	Sales14Total                []byte `byte:"len:10"` // 3700 = 0
	Sales15Total                []byte `byte:"len:10"` // 3700 = 0
	Sales16Total                []byte `byte:"len:10"` // 3700 = 0
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting16M19kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:19"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	Sales9Total                 []byte `byte:"len:10"` // 3700 = 0
	Sales10Total                []byte `byte:"len:10"` // 3700 = 0
	Sales11Total                []byte `byte:"len:10"` // 3700 = 0
	Sales12Total                []byte `byte:"len:10"` // 3700 = 0
	Sales13Total                []byte `byte:"len:10"` // 3700 = 0
	Sales14Total                []byte `byte:"len:10"` // 3700 = 0
	Sales15Total                []byte `byte:"len:10"` // 3700 = 0
	Sales16Total                []byte `byte:"len:10"` // 3700 = 0
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting16M16Packet struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	Sales9Total                 []byte `byte:"len:10"` // 3700 = 0
	Sales10Total                []byte `byte:"len:10"` // 3700 = 0
	Sales11Total                []byte `byte:"len:10"` // 3700 = 0
	Sales12Total                []byte `byte:"len:10"` // 3700 = 0
	Sales13Total                []byte `byte:"len:10"` // 3700 = 0
	Sales14Total                []byte `byte:"len:10"` // 3700 = 0
	Sales15Total                []byte `byte:"len:10"` // 3700 = 0
	Sales16Total                []byte `byte:"len:10"` // 3700 = 0
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	CheckOpenDate               []byte `byte:"len:8"`
	CheckOpenTime               []byte `byte:"len:8"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargePosting16M16kPacket struct { // POS => PMS
	SOH_                        []byte `byte:"len:1,equal:0x01"`
	WorkStation                 []byte `byte:"len:*"`
	STX_                        []byte `byte:"len:1,equal:0x02"`
	MessageType_                []byte `byte:"len:2,equal:0x2032"`
	MessageRetransmitFlag       []byte `byte:"len:1"`  /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	AccountID                   []byte `byte:"len:16"` // Inquiry
	ExpireDate                  []byte `byte:"len:4"`  // -
	SelectionInfo               []byte `byte:"len:16"`
	SelectionNumber             []byte `byte:"len:1"` // Match from posting list
	TransEmployeeNumber         []byte `byte:"len:4"`
	ChkEmployeeNumber           []byte `byte:"len:4"`
	RevenueCenterNumber         []byte `byte:"len:3"`
	ServingPeriodNumber         []byte `byte:"len:3"`
	GuestCheckNumber            []byte `byte:"len:4"`
	TransactionNumber           []byte `byte:"len:4"`
	NumberOfCovers              []byte `byte:"len:4"`
	CurrentPaymentNumber        []byte `byte:"len:3"`  // zahlungart
	CurrentPaymentAmount        []byte `byte:"len:10"` // totalamount
	Sales1Total                 []byte `byte:"len:10"`
	Sales2Total                 []byte `byte:"len:10"`
	Sales3Total                 []byte `byte:"len:10"`
	Sales4Total                 []byte `byte:"len:10"`
	Sales5Total                 []byte `byte:"len:10"`
	Sales6Total                 []byte `byte:"len:10"`
	Sales7Total                 []byte `byte:"len:10"`
	Sales8Total                 []byte `byte:"len:10"`
	Sales9Total                 []byte `byte:"len:10"` // 3700 = 0
	Sales10Total                []byte `byte:"len:10"` // 3700 = 0
	Sales11Total                []byte `byte:"len:10"` // 3700 = 0
	Sales12Total                []byte `byte:"len:10"` // 3700 = 0
	Sales13Total                []byte `byte:"len:10"` // 3700 = 0
	Sales14Total                []byte `byte:"len:10"` // 3700 = 0
	Sales15Total                []byte `byte:"len:10"` // 3700 = 0
	Sales16Total                []byte `byte:"len:10"` // 3700 = 0
	DiscountTotal               []byte `byte:"len:10"`
	ServiceChargeTotalEntered   []byte `byte:"len:10"` // tip
	ServiceChargeTotalAutomatic []byte `byte:"len:10"` // auslage
	Tax1Total                   []byte `byte:"len:10"` // schalter ignore!
	Tax2Total                   []byte `byte:"len:10"`
	Tax3Total                   []byte `byte:"len:10"`
	Tax4Total                   []byte `byte:"len:10"`
	Tax5Total                   []byte `byte:"len:10"`
	Tax6Total                   []byte `byte:"len:10"`
	Tax7Total                   []byte `byte:"len:10"`
	Tax8Total                   []byte `byte:"len:10"`
	PreviousPaymentTotal        []byte `byte:"len:10"`
	ETX_                        []byte `byte:"len:1,equal:0x03"`
	Checksum                    []byte `byte:"len:*"`
	EOT_                        []byte `byte:"len:1,equal:0x04"`
}

type TplOutletChargeResponseMPacket struct { // PMS => POS
	SOH_               []byte `byte:"len:1,equal:0x01"`
	WorkStation        []byte `byte:"len:*"`
	STX_               []byte `byte:"len:1,equal:0x02"`
	MessageType_       []byte `byte:"len:2,equal:0x2032"`
	AcceptanceDenial   []byte `byte:"len:16"`
	AdditionalMessages []byte `byte:"len:*"`
	ETX_               []byte `byte:"len:1,equal:0x03"`
	Checksum           []byte `byte:"len:*"`
	EOT_               []byte `byte:"len:1,equal:0x04"`
}

type TplCFReqMPacket struct { // POS => PMS
	SOH_                  []byte `byte:"len:1,equal:0x01"`
	ID                    []byte `byte:"len:*"`
	STX_                  []byte `byte:"len:1,equal:0x02"`
	MessageType_          []byte `byte:"len:2,equal:0x2033"`
	MessageRetransmitFlag []byte `byte:"len:1"` /* This field will be set to (space) (one space, Hex 0x20) for all initial Computer Inquire Request messages. This field will be set to R (Hex 0x52) to identify the applicationslevel retransmission of a previously sent Computer Inquire Request message (not a network-level retransmission). */
	RevenueCenterNumber   []byte `byte:"len:3"`
	GuestCheckNumber      []byte `byte:"len:4"`
	CheckOpenDate         []byte `byte:"len:8"`
	CheckOpenTime         []byte `byte:"len:8"`
	FieldSeperator1_      []byte `byte:"len:1,equal:0x1c"`
	CheckFacsimile        []byte `byte:"len:*"` // The actual check facsimile, variable in length. Each line is terminated with LF (0x0A). Each line can have a maximum length of 40 characters; the total size of this field may range up to 32k.
	FieldSeperator2_      []byte `byte:"len:1,equal:0x1c"`
	ETX_                  []byte `byte:"len:1,equal:0x03"`
	Checksum              []byte `byte:"len:*"`
	EOT_                  []byte `byte:"len:1,equal:0x04"`
}

type TplCFRspMPacket struct { // PMS => POS
	SOH_         []byte `byte:"len:1,equal:0x01"`
	WorkStation  []byte `byte:"len:*"`
	STX_         []byte `byte:"len:1,equal:0x02"`
	MessageType_ []byte `byte:"len:2,equal:0x2033"`
	PMSMessage   []byte `byte:"len:16"`
	ETX_         []byte `byte:"len:1,equal:0x03"`
	Checksum     []byte `byte:"len:*"`
	EOT_         []byte `byte:"len:1,equal:0x04"`
}
