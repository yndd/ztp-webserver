package webserver

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-dhcp/pkg/backend"
	dhcpstructs "github.com/yndd/ztp-dhcp/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/deviceregistry"
	"github.com/yndd/ztp-webserver/pkg/storage"
	storageIf "github.com/yndd/ztp-webserver/pkg/storage/interfaces"
	"github.com/yndd/ztp-webserver/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/utils"
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

	storageFs := os.DirFS(storageFolder)

	err := ws.index.LoadBackend(storageFs)
	if err != nil {
		log.Fatalf("error loading index backend: %v", err)
	}
	err = ws.storage.LoadBackend(storageFs)
	if err != nil {
		log.Fatalf("error loading storage backend: %v", err)
	}

	// get a pointer to the DeviceRegistry
	dr := deviceregistry.GetDeviceRegistry()

	// Get all the registered devices
	for _, x := range dr.GetRegistryDevices() {
		// inject the WebserverSetupper Interface reference
		x.SetWebserverSetupper(ws)
	}
	//ws.mux.Handle("/storage", http.StripPrefix("/storage", http.FileServer(http.Dir(storageFolder))))
	log.Infof("starting webserver on port %d", port)
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), ws.mux)
}

func (ws *WebserverImpl) AddHandler(up *structs.UrlParams, handler func(http.ResponseWriter, *http.Request)) {
	pattern := up.GetUrlRelative()
	path := "/" + pattern.Path
	log.Infof("handler added for pattern %s", path)
	ws.mux.HandleFunc(path+"/", handler)
}

func (ws *WebserverImpl) GetStorage() storageIf.Storage {
	return ws.storage
}

func (ws *WebserverImpl) GetIndex() storageIf.Index {
	return ws.index
}

// EnrichUrl add the Host portion to the given URL.
// The information is retrieved from the Kubernetes Service Kind
func (ws *WebserverImpl) EnrichUrl(theUrl *url.URL) error {
	wsi, err := ws.k8sBackend.GetWebserverInformation()
	if err != nil {
		return fmt.Errorf("error retrieving webserver information: %v", err)
	}

	theUrl.Scheme = wsi.Protocol
	port := fmt.Sprintf(":%d", wsi.Port)
	// if default http / https ports are in use,
	// strip the port portion from the URL
	switch wsi.Protocol {
	case "http":
		if wsi.Port == 80 {
			port = ""
		}
	case "https":
		if wsi.Port == 443 {
			port = ""
		}
	}
	// there is no specific port field, the port also goes into the
	// host portion of the url seperated via a colon (:)
	theUrl.Host = fmt.Sprintf("%s%s", wsi.IpFqdn, port)
	return nil
}

func (ws *WebserverImpl) ResponseFromIndex(rw http.ResponseWriter, r *http.Request) {
	fileEntry, err := ws.index.DeduceRelativeFilePath(r.URL)
	if err != nil {
		utils.HandleErrorCodeLog(404, err, rw)
	}
	switch fileEntry.ReferenceType {
	case structs.Filesystem:
		// handle local strorage delivery
		ws.storage.Handle(rw, fileEntry.Reference)
	case structs.HTTPRedirect:
		// Handle redirects
		http.Redirect(rw, r, fileEntry.Reference, http.StatusMovedPermanently)
		return
	}
}

func (ws *WebserverImpl) GetDeviceInformationByName(deviceId string) (*dhcpstructs.DeviceInformation, error) {
	// Forward the call to the k8sBackend
	return ws.k8sBackend.GetDeviceInformationByName(deviceId)
}

// GetWebserverOperations return the webserver operations interface
func GetWebserverOperations() webserverIf.WebserverOperations {
	return newWebserver()
}

func (ws *WebserverImpl) SetBackend(backend backend.ZtpBackend) {
	ws.k8sBackend = backend
}

// GetWebserverSetup return the webserver setup interface
func GetWebserverSetup() webserverIf.WebserverSetupper {
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
