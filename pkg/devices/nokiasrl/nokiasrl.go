package nokiasrl

import (
	"net/http"

	"github.com/yndd/ztp-webserver/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/webserver"
)

const (
	vendor = "nokia"
	model  = "srlinux"
)

var nokiasrl *NokiaSRL

type NokiaSRL struct{}

func (srl *NokiaSRL) handleSoftware(rw http.ResponseWriter, r *http.Request) {

}

func (srl *NokiaSRL) handleScript(rw http.ResponseWriter, r *http.Request) {

}

func (srl *NokiaSRL) handleConfig(rw http.ResponseWriter, r *http.Request) {

}

func NewNokiaSRL() *NokiaSRL {
	if nokiasrl == nil {
		nokiasrl = &NokiaSRL{}
	}
	return nokiasrl
}

func init() {
	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	webserver.GetWebserverSetup().AddHandler(upSoftware, NewNokiaSRL().handleSoftware)

	upScript := structs.NewUrlParams(vendor, model, structs.Script)
	webserver.GetWebserverSetup().AddHandler(upScript, NewNokiaSRL().handleScript)

	upConfig := structs.NewUrlParams(vendor, model, structs.Config)
	webserver.GetWebserverSetup().AddHandler(upConfig, NewNokiaSRL().handleConfig)
}
