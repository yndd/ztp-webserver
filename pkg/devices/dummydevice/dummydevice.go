package dummydevice

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/structs"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"

	"github.com/yndd/ztp-webserver/pkg/webserver"
)

const (
	vendor = "dummy"
	model  = "dummy"
)

var dummydevice *DummyDevice

type DummyDevice struct {
	webserver webserverIf.WebserverSetup
}

func (dd *DummyDevice) handleSoftware(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	dd.webserver.ResponseFromIndex(rw, r)
}

func NewDummyDevice(w webserverIf.WebserverSetup) *DummyDevice {
	if dummydevice == nil {
		dummydevice = &DummyDevice{webserver: w}
	}
	return dummydevice
}

func init() {
	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	wsSetup := webserver.GetWebserverSetup()
	wsSetup.AddHandler(upSoftware, NewDummyDevice(wsSetup).handleSoftware)
}
