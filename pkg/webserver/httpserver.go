package webserver

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	"github.com/yndd/ztp-dhcp/pkg/backend/k8s"
	dhcpstructs "github.com/yndd/ztp-dhcp/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/storage"
	storageIf "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	"github.com/yndd/ztp-webserver/pkg/structs"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"
)

var webserver *WebserverImpl

type WebserverImpl struct {
	mux        *http.ServeMux
	storage    storageIf.Storage
	index      storageIf.Index
	k8sBackend backend.ZtpBackend
}

func (ws *WebserverImpl) Run(port int, storageFolder string) {

	err := ws.index.LoadBackend(storageFolder)
	if err != nil {
		log.Fatalf("error loading index backend: %v", err)
	}
	err = ws.storage.LoadBackend(storageFolder)
	if err != nil {
		log.Fatalf("error loading storage backend: %v", err)
	}
	//ws.mux.Handle("/storage", http.StripPrefix("/storage", http.FileServer(http.Dir(storageFolder))))
	log.Infof("starting webserver on port %d", port)
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

func (ws *WebserverImpl) EnrichUrl(theUrl *url.URL) error {
	wsi, err := ws.k8sBackend.GetWebserverInformation()
	if err != nil {
		return fmt.Errorf("error retrieving webserver information: %v", err)
	}
	theUrl.Host = wsi.IpFqdn

	return nil
}

func (ws *WebserverImpl) ResponseFromIndex(rw http.ResponseWriter, r *http.Request, ct structs.ContentTypes) {
	relativeFileToBeDelivered, err := ws.index.DeduceRelativeFilePath(r.URL, ct)
	if err != nil {
		log.Errorf("error deducing relative file path: %v", err)
		status := http.StatusBadRequest
		rw.WriteHeader(status)
		rw.Write([]byte(fmt.Sprintf("%d - %v", status, err)))
		return
	}

	ws.storage.Handle(rw, relativeFileToBeDelivered)
}

func (ws *WebserverImpl) GetDeviceInformationByName(deviceId string) (*dhcpstructs.DeviceInformation, error) {
	// Forward the call to the k8sBackend
	return ws.k8sBackend.GetDeviceInformationByName(deviceId)
}

// GetWebserverOperations return the webserver operations interface
func GetWebserverOperations() webserverIf.WebserverOperations {
	return newWebserver()
}

func (ws *WebserverImpl) SetKubeConfig(kubeconfig string) {
	ws.k8sBackend = k8s.NewZtpK8sBackend(kubeconfig)
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
