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
	dd.webserver.ResponseFromIndex(rw, r, structs.Software)
}

func (dd *DummyDevice) handleMd5HashFile(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	dd.webserver.ResponseFromIndex(rw, r, structs.Md5HashFile)
}

func NewDummyDevice(w webserverIf.WebserverSetup) *DummyDevice {
	if dummydevice == nil {
		dummydevice = &DummyDevice{webserver: w}
	}
	return dummydevice
}

func init() {
	wsSetup := webserver.GetWebserverSetup()
	newdd := NewDummyDevice(wsSetup)

	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	wsSetup.AddHandler(upSoftware, newdd.handleSoftware)

	upMd5Hash := structs.NewUrlParams(vendor, model, structs.Md5HashFile)
	wsSetup.AddHandler(upMd5Hash, newdd.handleMd5HashFile)
}
