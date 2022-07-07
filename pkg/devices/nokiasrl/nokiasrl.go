package nokiasrl

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/utils"
	"github.com/yndd/ztp-webserver/pkg/webserver"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"
)

const (
	vendor = "nokia"
	model  = "srlinux"
)

//go:embed files/base_config.json.tmpl
var nokiaScriptTemplate string

//go:embed files/provisioning_01.py.tmpl
var nokiaConfigTemplate string

var nokiasrl *NokiaSRL

type NokiaSRL struct {
	webserver webserverIf.WebserverSetup
}

// handleSoftware handles the delivery of the software image to the client
func (srl *NokiaSRL) handleSoftware(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	// deliver the file registered in the index
	srl.webserver.ResponseFromIndex(rw, r, structs.Software)
}

// handleMd5HashFile handles the delivery of md5hash files to the client
func (srl *NokiaSRL) handleMd5HashFile(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)
	// deliver the file registered in the index
	srl.webserver.ResponseFromIndex(rw, r, structs.Md5HashFile)
}

// handleScript handles the generation of node specific ztp configuration scripts
func (srl *NokiaSRL) handleScript(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)

	// parse URL parameter to figure out the Node Name (node identifier)
	reqParams, err := structs.UrlParamsFromUrl(r.URL)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusBadRequest, err, rw)
		return
	}

	// chech that the device name is provided
	deviceId := reqParams.GetDeviceName()
	if deviceId == "" {
		utils.HandleErrorCodeLog(http.StatusBadRequest, fmt.Errorf("error: no device name provided"), rw)
		return
	}

	// retrieve the Topology node data from k8s
	nodeInformation, err := srl.webserver.GetDeviceInformationByName(deviceId)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}

	// generate URL for the software image file
	upSoftware := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Software).SetVersion(nodeInformation.ExpectedSWVersion).GetUrlRelative()
	// generate URL for the md5 hash file
	upHash := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Md5HashFile).SetVersion(nodeInformation.ExpectedSWVersion).GetUrlRelative()
	// generate URL for the node configuration
	upConfig := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Config).GetUrlRelative()

	// add hostname/ip, port and schema to url
	wss := webserver.GetWebserverSetup()
	wss.EnrichUrl(upSoftware)
	wss.EnrichUrl(upHash)
	wss.EnrichUrl(upConfig)

	// load the srl ztp script template
	t, err := template.New("script").Parse(nokiaScriptTemplate)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}
	// init buffer for the template result
	var specificScript = &bytes.Buffer{}

	// generate the node specific script from the tempalte
	err = t.Execute(specificScript, struct {
		ImageUrl  string
		Md5Url    string
		ConfigUrl string
	}{
		ImageUrl:  upSoftware.String(),
		Md5Url:    upHash.String(),
		ConfigUrl: upConfig.String(),
	})
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}
	// finally send the data to the client
	_, err = rw.Write(specificScript.Bytes())
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}
}

// handleConfig handles the generation of base srl configs
func (srl *NokiaSRL) handleConfig(rw http.ResponseWriter, r *http.Request) {
	log.Debugf("handling call on %s", r.URL)

	reqParams, err := structs.UrlParamsFromUrl(r.URL)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusBadRequest, err, rw)
		return
	}

	deviceId := reqParams.GetDeviceName()
	if deviceId == "" {
		utils.HandleErrorCodeLog(http.StatusBadRequest, fmt.Errorf("error: no device name provided"), rw)
		return
	}

	var templateFuncs = template.FuncMap{
		"join": strings.Join,
		"jsonstringify": func(sarr []string) []string {
			result := []string{}
			for _, s := range sarr {
				result = append(result, fmt.Sprintf("\"%s\"", s))
			}
			return result
		},
	}

	t, err := template.New("config").Funcs(templateFuncs).Parse(nokiaConfigTemplate)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}

	// init buffer for the template result
	var specificConfig = &bytes.Buffer{}

	nodeInformation, err := srl.webserver.GetDeviceInformationByName(reqParams.GetDeviceName())
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}

	// generate the node specific script from the tempalte
	err = t.Execute(specificConfig, struct {
		Cidr       string
		Netmask    int
		GatewayIp  string
		DnsServers []string
		NtpServers []string
	}{
		Cidr:       nodeInformation.CIDR,
		GatewayIp:  nodeInformation.Gateway,
		DnsServers: []string{"1.1.1.1", "5.5.5.5", "1.1.1.2", "5.5.5.6"},
		NtpServers: []string{"1.1.1.2", "5.5.5.6", "1.1.1.2", "5.5.5.6"},
	})
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}
	// finally send the data to the client
	_, err = rw.Write(specificConfig.Bytes())
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}

}

// NewNokiaSRL konstructor for the NokiaSRL device endpoint
func NewNokiaSRL(w webserverIf.WebserverSetup) *NokiaSRL {
	if nokiasrl == nil {
		nokiasrl = &NokiaSRL{webserver: w}
	}
	return nokiasrl
}

func init() {
	// register the endpoints this Device model supports in the webserver
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
