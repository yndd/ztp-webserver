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
	AddHandler(*structs.UrlParams, func(http.ResponseWriter, *http.Request))
	GetStorage() storageif.Storage
	GetIndex() storageif.Index
	ResponseFromIndex(http.ResponseWriter, *http.Request, structs.ContentTypes)
}
