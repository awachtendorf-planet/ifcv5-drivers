package ht24

var (
	answerText = map[string]string{
		"ES": "Syntax error",
		"NC": "The encoder does not answer. Failure in the communications, switched off, do not exist...",
		"NF": "Data files in the PC interface are damaged or not found",
		"OV": "The encoder has not already accomplished the previous task",
		"EP": "Error in magnetic track, card inserted wrongly or without magnetic stripe",
		"EF": "The card has been encoded by another system or the magnetic strip may be damaged",
		"EN": "The card has been encoded with a too low magnetic level due to dust in the reader magnetic head or low quality card",
		"ET": "Stuck card. The card does not have the required physical size",
		"TD": "Unknown card",
		"ED": "Timeout error. Operation cancelled",
		"EA": "Requested copy corresponds to a room checked-out",
		"OS": "Room out of service",
		"EO": "The requested card is being recorded by other station",
		"EE": "Task canceled",

		"E121": "Difference between written and read media data",
		"E129": "Writeerror / Deleteerror (Media not in sendrange)",
		"E130": "Illegal application",
		"E131": "Illegal company",
		"E132": "Timeout. Time between read and write or read and delete or read and levlection too long",
		"E133": "Datadevice full. Level can not be created or too much data for level creation",
		"E134": "Illegal writedata. Chip-ID, zeroinfo, application- or companynumber have been changed",
		"E135": "Level not existing",
		"E136": "Timeout. Time between read and delete too long",
		"E137": "Wrong password",
		"E138": "Unknown command",
		"E139": "Error at writing configuration data",
		"E141": "Error at reading media. No media can be accessed",
		"E142": "Unknown or illegal parameter",
		"E143": "Change setup parameters illegal or illegal value",
		"E144": "Illegal parameter at change setup parameters",
		"E145": "Access denial at access to a level. No write- or readpermissions",
		"E146": "Error at writing flash memory",
		"E147": "Wrong or illegal level chosen",
		"E148": "No read permissions for chosen level",
		"E150": "Error 150 in internal process, please contact vendor",
		"E151": "Error 151 in internal process, please contact vendor",
		"E152": "Error 152 in internal process, please contact vendor",
		"E153": "New date not accepted (because smaller than last transmitted)",
		"E154": "Illegal levelnumber (less 1 or greater 6) at write process",
		"E155": "Leveldata are expected",
		"E156": "Too much leveldata at level creation",
		"E157": "Error at write to media (maybe defect media)",
		"E158": "Illegal command. Maybe new init command with correct password necessary! Origin e.g. power failure in keydetector",
		"E159": "No write access for actual level",
		"E160": "Error with write directory (at level creation)",
		"E161": "No permission to delete level",
		"E162": "No permission to delete level",
		"E163": "Error with write directory (at level deletion)",
		"E164": "Too much data at write level",
		"E165": "Illegal time at W-command transmitted",
		"E166": "Illegal date at W-command transmitted",
		"E168": "Level with application / company combination already existing",
	}
)

// GetAnswerText returns a clear text for fias
func (d *Dispatcher) GetAnswerText(answerStatus string) string {
	text, exist := answerText[answerStatus]
	if !exist {
		return answerStatus
	}
	return text
}
