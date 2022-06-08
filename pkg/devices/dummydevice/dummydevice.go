package dummydevice

import (
	"net/http"

	"github.com/yndd/ztp-webserver/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/webserver"
)

const (
	vendor = "dummy"
	model  = "dummy"
)

var dummydevice *DummyDevice

type DummyDevice struct{}

func (dd *DummyDevice) handleSoftware(rw http.ResponseWriter, r *http.Request) {

}

func NewDummyDevice() *DummyDevice {
	if dummydevice == nil {
		dummydevice = &DummyDevice{}
	}
	return dummydevice
}

func init() {
	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	webserver.GetWebserverSetup().AddHandler(upSoftware, NewDummyDevice().handleSoftware)
}
