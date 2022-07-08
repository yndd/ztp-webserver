package interfaces

import (
	"net/http"
	"net/url"

	dhcpstructs "github.com/yndd/ztp-dhcp/pkg/structs"
	storageif "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	"github.com/yndd/ztp-webserver/pkg/structs"
)

type WebserverSetupper interface {
	AddHandler(*structs.UrlParams, func(http.ResponseWriter, *http.Request))
	GetStorage() storageif.Storage
	GetIndex() storageif.Index
	ResponseFromIndex(http.ResponseWriter, *http.Request, structs.ContentTypes)
	EnrichUrl(*url.URL) error
	GetDeviceInformationByName(deviceId string) (*dhcpstructs.DeviceInformation, error)
}
