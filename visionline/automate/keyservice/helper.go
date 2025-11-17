package keyservice

import (
	"sort"
	"strings"

	"github.com/weareplanet/ifcv5-drivers/visionline/template"

	"github.com/spf13/cast"
)

func (p *Plugin) trimTrack(track string) string {
	for len(track) > 2 {
		padding := track[len(track)-2:]
		padding = strings.ToLower(padding)
		if padding == "ff" {
			track = track[:len(track)-2]
		} else {
			break
		}
	}
	return track
}

func (p *Plugin) normalizeRoom(data string) []string {

	var rooms []string
	tmp := strings.Split(data, ",")

	for i := range tmp {
		room := tmp[i]
		index := strings.Index(room, "-")
		if index == -1 {
			rooms = append(rooms, room)
		} else if index > 0 && index+1 < len(room) {
			from := cast.ToUint(room[:index])
			to := cast.ToUint(room[index+1:])
			for counter := from; counter <= to; counter++ {
				rooms = append(rooms, cast.ToString(counter))
			}
		}
	}

	return rooms

}

func (p *Plugin) reconstructKeyOptions(data string, station uint64) string {

	ko := []int{}

	if len(data) == 0 {
		return ""
	}

	var keyOption string
	automate := p.automate
	dispatcher := automate.Dispatcher()

	// exact match -> lookup "Garage,spa"
	keyOption = dispatcher.GetPMSMapping(template.PacketReadKey, station, "CR", data)
	if len(keyOption) == 0 || keyOption == data {
		title := strings.Title(data)
		keyOption = dispatcher.GetPMSMapping(template.PacketReadKey, station, "CR", title)
		if len(keyOption) == 0 || keyOption == title {
			keyOption = ""
		}

	}

	// part match -> lookup "Garage", "spa"
	if len(keyOption) == 0 {

		aps := strings.Split(data, ",")

		for i := range aps {
			ap := strings.Trim(aps[i], " ")
			keyOption = dispatcher.GetPMSMapping(template.PacketReadKey, station, "CR", ap)
			if len(keyOption) == 0 || keyOption == ap {
				ap = strings.ToUpper(ap)
				keyOption = dispatcher.GetPMSMapping(template.PacketReadKey, station, "CR", ap)
			}
			if len(keyOption) == 0 || keyOption == ap {
				ap = strings.ToLower(ap)
				keyOption = dispatcher.GetPMSMapping(template.PacketReadKey, station, "CR", ap)
			}
			if len(keyOption) == 0 || keyOption == ap {
				ap = strings.Title(ap)
				keyOption = dispatcher.GetPMSMapping(template.PacketReadKey, station, "CR", ap)
			}

			count := cast.ToInt(keyOption)
			if count > 0 {
				ko = append(ko, count)
			}

		}

	} else {

		count := cast.ToInt(keyOption)
		if count > 0 {
			ko = append(ko, count)
		}

	}

	sort.Slice(ko, func(i, j int) bool {
		return ko[i] < ko[j]
	})

	if len(ko) == 0 {
		return ""
	}

	max := ko[len(ko)-1]
	if max > 128 {
		max = 128
	}

	ap := []byte(strings.Repeat("0", max))

	for i := range ko {
		if ko[i] <= len(ap) {
			ap[ko[i]-1] = '1'
		}
	}

	keyOption = strings.TrimRight(string(ap), "0")
	return keyOption

}
