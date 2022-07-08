package structs

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

var fakeserver = "127.0.0.66"
var testdata = []struct {
	in              *UrlParams
	outGetUrlParams string
}{{
	in:              NewUrlParams("Nokia", "SRLinux", Config).SetDeviceId("MyDevice43"),
	outGetUrlParams: fmt.Sprintf("http://%s/Nokia/SRLinux/config?deviceid=MyDevice43", fakeserver),
}, {
	in:              NewUrlParams("Nokia", "SRLinux", Config).SetDeviceId("MyOtherDevice").SetVersion("v8.7.6"),
	outGetUrlParams: fmt.Sprintf("http://%s/Nokia/SRLinux/config?deviceid=MyOtherDevice&version=v8.7.6", fakeserver),
},
}

func TestUrlParamsGetUrlRelative(t *testing.T) {
	for _, x := range testdata {
		relativeUrl := x.in.GetUrlRelative()
		relativeUrl.Host = fakeserver
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
