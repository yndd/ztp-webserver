package dummydevice

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/deviceregistry"
	"github.com/yndd/ztp-webserver/pkg/structs"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"
)

const (
	vendor = "Dummy"
	model  = "Dummy"
)

var dummydevice *DummyDevice

type DummyDevice struct {
	webserver webserverIf.WebserverSetupper
}

func (dd *DummyDevice) handleSoftware(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	dd.webserver.ResponseFromIndex(rw, r, structs.Software)
}

func (dd *DummyDevice) handleMd5HashFile(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	dd.webserver.ResponseFromIndex(rw, r, structs.Md5HashFile)
}

// SetWebserverSetupper used to inject the dependency of the WebserverSetupper
func (dd *DummyDevice) SetWebserverSetupper(webserver webserverIf.WebserverSetupper) {
	dd.webserver = webserver

	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	webserver.AddHandler(upSoftware, dd.handleSoftware)

	upMd5Hash := structs.NewUrlParams(vendor, model, structs.Md5HashFile)
	webserver.AddHandler(upMd5Hash, dd.handleMd5HashFile)
}

// NewDummyDevice singleton creator for the DummyDevice
func NewDummyDevice() *DummyDevice {
	if dummydevice == nil {
		dummydevice = &DummyDevice{}
	}
	return dummydevice
}

func init() {
	// create a new DummyDevice instance
	newdd := NewDummyDevice()
	// acquire the handle on the deviceregistry
	dr := deviceregistry.GetDeviceRegistry()
	// register the device with the registry
	dr.RegisterDevice(newdd)
}
