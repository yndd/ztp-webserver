package webserver

import (
	"fmt"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/structs"
)

var webserver *WebserverImpl

type WebserverImpl struct {
	mux *http.ServeMux
}

type WebserverOperations interface {
	Run(ip string, port int)
}

type WebserverSetup interface {
	AddHandler(up *structs.UrlParams, handler func(http.ResponseWriter, *http.Request))
}

func (ws *WebserverImpl) Run(ip string, port int) {
	log.Infof("starting webserver on %s:%d", ip, port)
	http.ListenAndServe(ip+":"+strconv.Itoa(port), ws.mux)
}

func (ws *WebserverImpl) AddHandler(up *structs.UrlParams, handler func(http.ResponseWriter, *http.Request)) {
	pattern := up.GetUrlRelative()
	log.Infof("handler added for pattern %s", pattern)
	ws.mux.HandleFunc(pattern, handler)
}

// GetWebserverOperations return the webserver operations interface
func GetWebserverOperations() WebserverOperations {
	return newWebserver()
}

// GetWebserverSetup return the webserver setup interface
func GetWebserverSetup() WebserverSetup {
	return newWebserver()
}

// newWebserver constructor for the singleton webserver
func newWebserver() *WebserverImpl {
	if webserver == nil {
		webserver = &WebserverImpl{
			mux: http.NewServeMux(),
		}
	}
	return webserver
}

// HandleErrorCode Generates error page
func HandleErrorCode(errorCode int, description string, w http.ResponseWriter) {
	w.WriteHeader(errorCode)                    // set HTTP status code (example 404, 500)
	w.Header().Set("Content-Type", "text/html") // clarify return type (MIME)
	w.Write([]byte(fmt.Sprintf(
		"<html><body><h1>Error %d</h1><p>%s</p></body></html>",
		errorCode,
		description)))
}
