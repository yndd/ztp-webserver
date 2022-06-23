package nokiasrl

import (
	"bytes"
	"fmt"
	"net/http"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/webserver"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"
)

const (
	vendor = "nokia"
	model  = "srlinux"
)

var nokiasrl *NokiaSRL

type NokiaSRL struct {
	webserver webserverIf.WebserverSetup
}

func (srl *NokiaSRL) handleSoftware(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	srl.webserver.ResponseFromIndex(rw, r, structs.Software)
}

func (srl *NokiaSRL) handleMd5HashFile(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	srl.webserver.ResponseFromIndex(rw, r, structs.Md5HashFile)
}

func (srl *NokiaSRL) handleScript(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)

	t, err := template.ParseFiles("./pkg/devices/nokiasrl/files/provisioning_01.py.tmpl")
	if err != nil {
		log.Errorf("error deducing relative file path: %v", err)
		status := http.StatusInternalServerError
		rw.WriteHeader(status)
		rw.Write([]byte(fmt.Sprintf("%d - %v", status, err)))
		return
	}
	var specificScript = &bytes.Buffer{}
	err = t.Execute(specificScript, struct {
		ImageUrl  string
		Md5Url    string
		ConfigUrl string
	}{
		ImageUrl:  "imageURL",
		Md5Url:    "md5url",
		ConfigUrl: "configURL",
	})
	if err != nil {
		log.Errorf("error generating provisioning script: %v", err)
		status := http.StatusInternalServerError
		rw.WriteHeader(status)
		rw.Write([]byte(fmt.Sprintf("%d - %v", status, err)))
		return
	}
	_, err = rw.Write(specificScript.Bytes())
	if err != nil {
		log.Errorf("error delivering script content: %v", err)
		status := http.StatusInternalServerError
		rw.WriteHeader(status)
		rw.Write([]byte(fmt.Sprintf("%d - %v", status, err)))
		return
	}

}

func (srl *NokiaSRL) handleConfig(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
}

func NewNokiaSRL(w webserverIf.WebserverSetup) *NokiaSRL {
	if nokiasrl == nil {
		nokiasrl = &NokiaSRL{webserver: w}
	}
	return nokiasrl
}

func init() {

	wsSetup := webserver.GetWebserverSetup()
	nsrl := NewNokiaSRL(wsSetup)

	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	wsSetup.AddHandler(upSoftware, nsrl.handleSoftware)

	upScript := structs.NewUrlParams(vendor, model, structs.Script)
	wsSetup.AddHandler(upScript, nsrl.handleScript)

	upConfig := structs.NewUrlParams(vendor, model, structs.Config)
	wsSetup.AddHandler(upConfig, nsrl.handleConfig)

	upmd5hash := structs.NewUrlParams(vendor, model, structs.Md5HashFile)
	wsSetup.AddHandler(upmd5hash, nsrl.handleMd5HashFile)
}
