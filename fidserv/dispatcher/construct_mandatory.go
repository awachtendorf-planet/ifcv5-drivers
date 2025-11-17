package fidserv

var (
	outgoingMandatory = map[string][]string{

		"DE": {"DA", "TI"},
		"DR": {"DA", "TI"},
		"DS": {"DA", "TI"},

		"GI": {"G#", "RN", "GS"},
		"GO": {"G#", "RN", "GS"},
		"GC": {"G#", "RN", "GS"},

		"KA": {"AS", "CT", "KC", "WS"},
		"KD": {"KC", "RN", "WS"},
		"KM": {"G#", "KC", "RN", "RO", "WS"},
		"KR": {"KC", "KT", "RN", "WS"},
		"KZ": {"AS", "CT", "KC", "RN", "WS"},

		"NE": {"DA", "TI"},
		"NS": {"DA", "TI"},

		"RE": {"RN"},

		//"PA": {"AS", "CT", "DA", "P#", "RN", "TI", "WS"},
		"PA": {"AS", "CT", "DA", "P#", "TI", "WS"}, // AS not OK, there is no RN
		"PL": {"G#", "GN", "P#", "RN", "WS"},

		"CK": {"AS", "C#", "CT", "DA", "TI", "SO"},

		"WA": {"AS", "RN", "DA", "TI"},
		"WC": {"RN", "DA", "TI"},
		"WR": {"RN", "DA", "TI"},

		"XB": {"BA", "G#", "RN"},
		"XC": {"AS", "BA", "CT", "G#", "RN"},
		"XD": {"G#", "MI", "RN"},
		"XI": {"BD", "BI", "DC", "G#", "RN", "F#", "FD"},
		"XL": {"G#", "MI", "MT", "RN"},
		//"XT": {"G#", "MI", "MT", "RN"},
		"XT": {"G#", "RN"}, // without MI/MT signals that no guest message exist
	}
)

// MandatoryField returns true if a  record/field is mantatory
func (d *Dispatcher) MandatoryField(station uint64, recordName string, key string) bool {
	if fields, exist := outgoingMandatory[recordName]; exist {
		for i := range fields {
			if fields[i] == key {
				return true
			}
		}
	}
	return false
}
