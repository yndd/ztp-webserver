package nokiasrl

import (
	"net/http"

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
	srl.webserver.ResponseFromIndex(rw, r)
}

func (srl *NokiaSRL) handleScript(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
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
}
