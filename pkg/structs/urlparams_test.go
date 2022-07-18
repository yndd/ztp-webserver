package structs

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var fakeserver = "127.0.0.66"
var testdata = []struct {
	in              *UrlParams
	outGetUrlParams string
}{{
	in:              NewUrlParams("Nokia", "SRLinux", Config).SetDeviceId("MyDevice43"),
	outGetUrlParams: "Nokia/SRLinux/config?deviceid=MyDevice43",
}, {
	in:              NewUrlParams("Nokia", "SRLinux", Config).SetDeviceId("MyOtherDevice").SetVersion("v8.7.6"),
	outGetUrlParams: "Nokia/SRLinux/config?deviceid=MyOtherDevice&version=v8.7.6",
}, {
	in:              NewUrlParams("Nokia", "SRLinux", Config).SetDeviceId("MyOtherDevice").SetVersion("v8.7.6").SetFilename("myfile.py"),
	outGetUrlParams: "Nokia/SRLinux/config/myfile.py?deviceid=MyOtherDevice&version=v8.7.6",
},
}

func TestUrlParamsGetUrlRelative(t *testing.T) {
	for _, x := range testdata {
		relativeUrl := x.in.GetUrlRelative()
		if relativeUrl.String() != x.outGetUrlParams {
			t.Errorf("geturlrelative expected %s got %s from %s", x.outGetUrlParams, x.in.GetUrlRelative(), spew.Sprint(x.in))
		}
	}
}

func TestUrlParamsParseURLParams(t *testing.T) {
	for _, x := range testdata {
		z, err := ParseURL(x.outGetUrlParams)
		if err != nil {
			t.Errorf("error parsing %s, (%v)", x.outGetUrlParams, err)
		}
		if !z.Equals(x.in) {
			t.Errorf("error in parsing %s != %s", spew.Sprint(z), spew.Sprint(x.in))
		}
	}
}

func TestNewUrlParamsDeviceName(t *testing.T) {
	model := "SRLinux"
	vendor := "Nokia"
	deviceId := "Device85"
	version := "v6.7"
	contentType := Config

	urlp := NewUrlParamsDeviceName(vendor, model, deviceId, contentType)
	urlp.SetVersion(version)

	if urlp.GetModel() != model {
		t.Errorf("Expected model to be %s, but was %s", model, urlp.model)
	}
	if urlp.GetVendor() != vendor {
		t.Errorf("Expected vendor to be %s, but was %s", vendor, urlp.GetVendor())
	}
	if urlp.GetDeviceName() != deviceId {
		t.Errorf("Expected DeviceId to be %s, but was %s", deviceId, urlp.GetDeviceName())
	}
	if urlp.GetVersion() != version {
		t.Errorf("Expected Version to be %s, but was %s", version, urlp.GetVersion())
	}
	if urlp.GetContentType() != contentType {
		t.Errorf("Expected ContentType to be %s, but was %s", contentType, urlp.GetContentType())
	}
}

func TestUrlParamsFromUrlErrors(t *testing.T) {

	// test for too many path elements
	u, _ := url.Parse(fmt.Sprintf("http://%s/Nokia/SRLinux/config/foo/bla?deviceid=MyOtherDevice&version=v8.7.6", fakeserver))

	_, err := UrlParamsFromUrl(u)
	if err == nil {
		t.Errorf("Expected error to be thrown since too many path elements are present in url %s. ", u.String())
	}

	// test for wrong content type path elment
	u, _ = url.Parse(fmt.Sprintf("http://%s/Nokia/SRLinux/blabla?deviceid=MyOtherDevice&version=v8.7.6", fakeserver))

	_, err = UrlParamsFromUrl(u)
	if err == nil {
		t.Errorf("Expected error to be thrown since content type is not a registered type in url %s. ", u.String())
	}
}
