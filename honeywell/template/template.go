package template

const (
	PacketAck = "ACK"
	PacketNak = "NAK"

	PacketCheckIn  = "Check In"
	PacketCheckOut = "Check Out"

	PacketGarbage  = "Garbage"
	PacketResponse = "EMC Response"
)

type TplACK struct {
	ACK_ []byte `byte:"len:1,equal:0x06"`
}

type TplNAK struct {
	NAK_ []byte `byte:"len:1,equal:0x15"`
}

type TplGarbage_ACK struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x06"`
}

type TplGarbage_NAK struct {
	Data_ []byte `byte:"len:*"`
	STX_  []byte `byte:"len:1,equal:0x15"`
}

// HONEYWELL

// HoneywellGarbage_CR overread until CR
type HoneywellGarbage_CR struct {
	Data_ []byte `byte:"len:*"`
	CR_   []byte `byte:"len:1,equal:0x0d"`
}

// HoneywellGarbage_CI overread until CI, rewind 2
type HoneywellGarbage_CI struct {
	Data_ []byte `byte:"len:*"`
	CI_   []byte `byte:"len:2,equal:0x4349"`
}

// HoneywellGarbage_CO overread until CO, rewind 2
type HoneywellGarbage_CO struct {
	Data_ []byte `byte:"len:*"`
	CO_   []byte `byte:"len:2,equal:0x434f"`
}

// HoneywellGarbage_Unknown overread if more then 4 bytes
type HoneywellGarbage_Unknown struct {
	Data_ []byte `byte:"len:4"`
}

// HoneywellCheckin ...
type HoneywellCheckin struct {
	CI_  []byte `byte:"len:2,equal:0x4349"` // CI
	Room []byte `byte:"len:4"`              // 1...9999
	CR_  []byte `byte:"len:1,equal:0x0d"`   // CR
}

// HoneywellCheckout ...
type HoneywellCheckout struct {
	CO_  []byte `byte:"len:2,equal:0x434f"` // CO
	Room []byte `byte:"len:4"`              // 1...9999
	CR_  []byte `byte:"len:1,equal:0x0d"`   // CR
}

// HoneywellResponse ...
type HoneywellResponse struct {
	C_      []byte `byte:"len:1,equal:0x43"` // C
	Request []byte `byte:"len:1"`            // 49 I, 4F O
	Answer  []byte `byte:"len:1"`            // 06 ACK, 15 NAK
	CR_     []byte `byte:"len:1,equal:0x0d"` // CR
}

// ALERTON

// AlertonCheckin ...
type AlertonCheckin struct {
	CI_  []byte `byte:"len:1,equal:0x45"` // E
	Room []byte `byte:"len:4"`            // 1...9999
}

// AlertonCheckout ...
type AlertonCheckout struct {
	CO_  []byte `byte:"len:1,equal:0x56"` // V
	Room []byte `byte:"len:4"`            // 1...9999
}
