package dummydevice

import (
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

// SetWebserverSetupper used to inject the dependency of the WebserverSetupper
func (dd *DummyDevice) SetWebserverSetupper(webserver webserverIf.WebserverSetupper) {
	dd.webserver = webserver

	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	// register ResponseFromIndex, because its just a file that can come straight from index
	webserver.AddHandler(upSoftware, dd.webserver.ResponseFromIndex)

	upMd5Hash := structs.NewUrlParams(vendor, model, structs.Md5HashFile)
	// register ResponseFromIndex, because its just a file that can come straight from index
	webserver.AddHandler(upMd5Hash, dd.webserver.ResponseFromIndex)
}

// GetDummyDevice singleton creator for the DummyDevice
func GetDummyDevice() *DummyDevice {
	if dummydevice == nil {
		dummydevice = &DummyDevice{}
	}
	return dummydevice
}

func init() {
	// create a new DummyDevice instance
	newdd := GetDummyDevice()
	// acquire the handle on the deviceregistry
	dr := deviceregistry.GetDeviceRegistry()
	// register the device with the registry
	dr.RegisterDevice(newdd)
}
