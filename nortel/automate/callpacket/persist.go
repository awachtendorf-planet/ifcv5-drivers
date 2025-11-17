package callpacket

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/weareplanet/ifcv5-main/ifc/generic/dispatcher"
	"github.com/weareplanet/ifcv5-main/log"
	"github.com/weareplanet/ifcv5-main/utils"
)

const (
	pendingFile = "pending.records.json"
)

type records struct {
	Records []PbxRecord `json:"records"`
}

func (p *Plugin) loadPendingRecords() {

	var pending records

	fileName := dispatcher.ConfigPath + "/" + pendingFile
	if !utils.Exists(fileName) {
		return
	}

	data, err := ioutil.ReadFile(fileName)
	if err == nil {

		if err = json.Unmarshal(data, &pending); err == nil {

			p.recordsGuard.Lock()

			for i := range pending.Records {
				record := pending.Records[i]
				idx := p.getIndex(record.Station, record.LineNumber)
				p.records[idx] = record

			}

			p.recordsGuard.Unlock()

		}
	}

	if err == nil {
		log.Debug().Msgf("%s load pending records, count: %d", p.GetName(), len(pending.Records))
	} else {
		log.Error().Msgf("%s load pending records failed, err=%s", p.GetName(), err)
	}

	os.Remove(fileName)
}

func (p *Plugin) storePendingRecords() {

	var pending records

	p.recordsGuard.Lock()
	for i := range p.records {
		pending.Records = append(pending.Records, p.records[i])
	}
	p.recordsGuard.Unlock()

	if len(pending.Records) == 0 {
		return
	}

	data, err := json.Marshal(pending)
	if err == nil {
		fileName := dispatcher.ConfigPath + "/" + pendingFile
		err = ioutil.WriteFile(fileName, data, 0644)
	}

	if err == nil {
		log.Debug().Msgf("%s store pending records, count: %d", p.GetName(), len(pending.Records))
	} else {
		log.Error().Msgf("%s store pending records failed, err=%s", p.GetName(), err)
	}
}
