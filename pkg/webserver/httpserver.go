package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/storage"
	storageIf "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	"github.com/yndd/ztp-webserver/pkg/structs"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"
)

var webserver *WebserverImpl

type WebserverImpl struct {
	mux     *http.ServeMux
	storage storageIf.Storage
	index   storageIf.Index
}

func (ws *WebserverImpl) Run(port int, storageFolder string) {
	log.Infof("starting webserver on port %d", port)
	ws.index.LoadBackend(storageFolder)
	ws.storage.LoadBackend(storageFolder)
	ws.mux.Handle("/storage", http.StripPrefix("/storage", http.FileServer(http.Dir(storageFolder))))
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), ws.mux)
}

func (ws *WebserverImpl) AddHandler(up *structs.UrlParams, handler func(http.ResponseWriter, *http.Request)) {
	pattern := up.GetUrlRelative()
	path := "/" + pattern.Path
	log.Infof("handler added for pattern %s", path)
	ws.mux.HandleFunc(path, handler)
}

func (ws *WebserverImpl) GetStorage() storageIf.Storage {
	return ws.storage
}

func (ws *WebserverImpl) GetIndex() storageIf.Index {
	return ws.index
}

func (ws *WebserverImpl) ResponseFromIndex(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)

	relativeFileToBeDelivered, err := ws.index.DeduceRelativeFilePath(r.URL)
	if err != nil {
		status := http.StatusBadRequest
		rw.WriteHeader(status)
		rw.Write([]byte(fmt.Sprintf("%d - %v", status, err)))
		return
	}

	ws.storage.Handle(rw, relativeFileToBeDelivered)
}

// GetWebserverOperations return the webserver operations interface
func GetWebserverOperations() webserverIf.WebserverOperations {
	return newWebserver()
}

// GetWebserverSetup return the webserver setup interface
func GetWebserverSetup() webserverIf.WebserverSetup {
	return newWebserver()
}

// newWebserver constructor for the singleton webserver
func newWebserver() *WebserverImpl {
	if webserver == nil {
		webserver = &WebserverImpl{
			mux:     http.NewServeMux(),
			storage: storage.NewFolderStorage(),
			index:   storage.NewIndex(),
		}
	}
	return webserver
}
