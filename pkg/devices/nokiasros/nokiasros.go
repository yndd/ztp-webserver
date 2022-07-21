package nokiasros

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
	model  = "SROS"
)

//go:embed files/config.cfg.tmpl
var nokiaConfigTemplate string

//go:embed files/provision.txt.tmpl
var nokiaProvisionTemplate string

// nokiasros stores the singleton instance of the NokiaSROS instance
var nokiasros *NokiaSROS

type NokiaSROS struct {
	webserver webserverIf.WebserverSetupper
}

// handleScript handles the generation of node specific ztp configuration scripts
func (sros *NokiaSROS) handleScript(rw http.ResponseWriter, r *http.Request) {
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
	nodeInformation, err := sros.webserver.GetDeviceInformationByName(deviceId)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}

	softwareFilename := fmt.Sprintf("srlinux-%s.bin", nodeInformation.ExpectedSWVersion)

	// generate URL for the software image file
	upSoftware := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Software).SetVersion(nodeInformation.ExpectedSWVersion).SetFilename(softwareFilename).GetUrlRelative()
	// generate URL for the md5 hash file
	upHash := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Md5HashFile).SetVersion(nodeInformation.ExpectedSWVersion).SetFilename(softwareFilename + ".md5").GetUrlRelative()
	// generate URL for the node configuration
	upConfig := structs.NewUrlParamsDeviceName(vendor, model, deviceId, structs.Config).SetFilename("device.cfg").GetUrlRelative()

	// add hostname/ip, port and schema to url
	wss := webserver.GetWebserverSetup()
	wss.EnrichUrl(upSoftware)
	wss.EnrichUrl(upHash)
	wss.EnrichUrl(upConfig)

	// load the provisioning ztp script template
	t, err := template.New("script").Funcs(getTemplatingFunctions()).Parse(nokiaProvisionTemplate)
	if err != nil {
		utils.HandleErrorCodeLog(http.StatusInternalServerError, err, rw)
		return
	}
	// init buffer for the template result
	var specificScript = &bytes.Buffer{}

	// generate the node specific script from the tempalte
	err = t.Execute(specificScript, struct {
		CpmTimUrl     string
		IomTimUrl     string
		KernelTimUrl  string
		SupportTimUrl string
		BothTimUrl    string
	}{})
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
func (sros *NokiaSROS) handleConfig(rw http.ResponseWriter, r *http.Request) {
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

	nodeInformation, err := sros.webserver.GetDeviceInformationByName(reqParams.GetDeviceName())
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
		DnsServers: []string{"1.1.1.1", "8.8.8.8", "1.1.1.2"},
		NtpServers: []string{"1.1.1.2", "5.5.5.6", "1.1.1.2"},
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
func (sros *NokiaSROS) SetWebserverSetupper(webserver webserverIf.WebserverSetupper) {
	sros.webserver = webserver

	upSoftware := structs.NewUrlParams(vendor, model, structs.Software)
	// delegate software retrieval to the webserver to respond straight from index
	webserver.AddHandler(upSoftware, sros.webserver.ResponseFromIndex)

	upmd5hash := structs.NewUrlParams(vendor, model, structs.Md5HashFile)
	// delegate md5hash retrieval to the webserver to respond straight from index
	webserver.AddHandler(upmd5hash, sros.webserver.ResponseFromIndex)

	upScript := structs.NewUrlParams(vendor, model, structs.Script)
	webserver.AddHandler(upScript, sros.handleScript)

	upConfig := structs.NewUrlParams(vendor, model, structs.Config)
	webserver.AddHandler(upConfig, sros.handleConfig)

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
func GetNokiaSROS() *NokiaSROS {
	if nokiasros == nil {
		nokiasros = &NokiaSROS{}
	}
	return nokiasros
}

func init() {
	// create a new NokiaSRL instance
	newsrl := GetNokiaSROS()
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
