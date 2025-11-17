package datachange

func (p *Plugin) blockCommunication(addr string) bool {

	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	ret := dispatcher.BlockCommunication(addr, name)
	return ret

}

func (p *Plugin) freeCommunication(addr string) {
	automate := p.automate
	name := automate.Name()
	dispatcher := automate.Dispatcher()

	dispatcher.FreeCommunication(addr, name)
}
