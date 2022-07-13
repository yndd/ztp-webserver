package interfaces

import (
	"net/http"
	"net/url"

	dhcpstructs "github.com/yndd/ztp-dhcp/pkg/structs"
	storageif "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	"github.com/yndd/ztp-webserver/pkg/structs"
)

// WebserverSetupper interface, implemented by the webserver, that is utilizted by the device
// specific implementation to register handler, access the storage, delagate file requiests to the
// index or retrieve information on the requesting Node from the k8s-apiserver
type WebserverSetupper interface {
	// AddHandler called by the device specific implementations to register their
	// handler methods
	AddHandler(*structs.UrlParams, func(http.ResponseWriter, *http.Request))
	// retrieve the pointer to the implementation of the storage interface
	GetStorage() storageif.Storage
	// retrieve the pointer to the implementation of the index interface
	GetIndex() storageif.Index
	// ResponseFromIndex allows device specific implementations to delegate a
	// request to the index, which will deliver the underlying file
	ResponseFromIndex(http.ResponseWriter, *http.Request)
	// EnrichUrl enriches the given url with the schema, host and port portion
	// this information is retrieved from the k8s service entry of the webserver
	EnrichUrl(*url.URL) error
	// GetDeviceInformationByName retrieve DeviceInformation, from the
	// kubernetes-apiserver, necesarry to generate scripts or figure out which files
	// to deliver
	GetDeviceInformationByName(deviceId string) (*dhcpstructs.DeviceInformation, error)
}
