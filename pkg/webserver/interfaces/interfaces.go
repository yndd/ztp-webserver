package interfaces

import (
	"net/http"

	storageif "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	"github.com/yndd/ztp-webserver/pkg/structs"
)

type WebserverOperations interface {
	Run(port int, storageFolder string)
}

type WebserverSetup interface {
	AddHandler(up *structs.UrlParams, handler func(http.ResponseWriter, *http.Request))
	GetStorage() storageif.Storage
	GetIndex() storageif.Index
	ResponseFromIndex(rw http.ResponseWriter, r *http.Request)
}
