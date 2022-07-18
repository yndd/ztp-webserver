package nokiasrl

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/deviceregistry"
	"github.com/yndd/ztp-webserver/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/utils"
	"github.com/yndd/ztp-webserver/pkg/webserver"
	webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"
)

const (
	vendor = "Nokia"
	model  = "SRLinux"
)

//go:embed files/base_config.json.tmpl
var nokiaConfigTemplate string

//go:embed files/provisioning_01.py.tmpl
var nokiaScriptTemplate string

// nokiasrl stores the singleton instance of the NokiaSRL instance
var nokiasrl *NokiaSRL

type NokiaSRL struct {
	webserver webserverIf.WebserverSetupper
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
	upSoftware := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Software).SetVersion(nodeInformation.ExpectedSWVersion).SetFilename("srlinux.bin").GetUrlRelative()
	// generate URL for the md5 hash file
	upHash := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Md5HashFile).SetVersion(nodeInformation.ExpectedSWVersion).SetFilename("script.py").GetUrlRelative()
	// generate URL for the node configuration
	upConfig := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Config).SetFilename("srlinux.md5").GetUrlRelative()

	// add hostname/ip, port and schema to url
	wss := webserver.GetWebserverSetup()
	wss.EnrichUrl(upSoftware)
	wss.EnrichUrl(upHash)
	wss.EnrichUrl(upConfig)

	// load the srl ztp script template
	t, err := template.New("script").Funcs(getTemplatingFunctions()).Parse(nokiaScriptTemplate)
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

	// // Add content-disposition header
	// filename := reqParams.GetFilename()
	// if filename != "" {
	// 	filename = "provision.py"
	// }
	// rw.Header().Add("content-disposition", fmt.Sprintf("attachment; filename=%s", filename))

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

	t, err := template.New("config").Funcs(getTemplatingFunctions()).Parse(nokiaConfigTemplate)
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

// SetWebserverSetupper used to inject the dependency of the WebserverSetupper
func (srl *NokiaSRL) SetWebserverSetupper(webserver webserverIf.WebserverSetupper) {
	srl.webserver = webserver

	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	// delegate software retrieval to the webserver to respond straight from index
	webserver.AddHandler(upSoftware, srl.webserver.ResponseFromIndex)

	upmd5hash := structs.NewUrlParams(vendor, model, structs.Md5HashFile)
	// delegate md5hash retrieval to the webserver to respond straight from index
	webserver.AddHandler(upmd5hash, srl.webserver.ResponseFromIndex)

	upScript := structs.NewUrlParams(vendor, model, structs.Script)
	webserver.AddHandler(upScript, srl.handleScript)

	upConfig := structs.NewUrlParams(vendor, model, structs.Config)
	webserver.AddHandler(upConfig, srl.handleConfig)

}

// getTemplatingFunctions returns the function map used in multiple templating
// instances, e.g. Config and Script
func getTemplatingFunctions() template.FuncMap {
	var templateFuncs = template.FuncMap{
		"join":          strings.Join,
		"jsonstringify": jsonStringifyArray,
	}
	return templateFuncs
}

// GetNokiaSRL konstructor for the NokiaSRL device endpoint
func GetNokiaSRL() *NokiaSRL {
	if nokiasrl == nil {
		nokiasrl = &NokiaSRL{}
	}
	return nokiasrl
}

func init() {
	// create a new NokiaSRL instance
	newsrl := GetNokiaSRL()
	// acquire the handle on the deviceregistry
	dr := deviceregistry.GetDeviceRegistry()
	// register the device with the registry
	dr.RegisterDevice(newsrl)
}

// jsonStringifyArray function used in golang templating
// this adds Quotationmarks to the strings in the array
func jsonStringifyArray(sarr []string) []string {
	result := []string{}
	for _, s := range sarr {
		result = append(result, fmt.Sprintf("\"%s\"", s))
	}
	return result
}
