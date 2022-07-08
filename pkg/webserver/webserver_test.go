package webserver

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	dhcp_mocks "github.com/yndd/ztp-dhcp/pkg/mocks"
	dhcp_structs "github.com/yndd/ztp-dhcp/pkg/structs"
	"github.com/yndd/ztp-webserver/pkg/structs"
)

func TestEnrichUrl(t *testing.T) {

	testdef := []struct {
		description string
		got_port    int32
		got_ipfqdn  string
		got_proto   string
		want_scheme string
		want_host   string
	}{
		{
			description: "HTTP with port 80",
			got_port:    int32(80),
			got_ipfqdn:  "myfunnywebserverIp",
			got_proto:   "http",
			want_scheme: "http",
			want_host:   "myfunnywebserverIp",
		},
		{
			description: "HTTPS with port 443",
			got_port:    int32(443),
			got_ipfqdn:  "myfunnywebserverIp",
			got_proto:   "https",
			want_scheme: "https",
			want_host:   "myfunnywebserverIp",
		},
		{
			description: "HTTP with port 90",
			got_port:    int32(90),
			got_ipfqdn:  "myfunnywebserverIp",
			got_proto:   "http",
			want_scheme: "http",
			want_host:   "myfunnywebserverIp:90",
		},
		{
			description: "HTTPS with port 90",
			got_port:    int32(90),
			got_ipfqdn:  "myfunnywebserverIp",
			got_proto:   "https",
			want_scheme: "https",
			want_host:   "myfunnywebserverIp:90",
		},
	}

	for _, test := range testdef {

		mockCtrl := gomock.NewController(t)

		k8b := dhcp_mocks.NewMockZtpBackend(mockCtrl)
		k8b.EXPECT().GetWebserverInformation().Return(&dhcp_structs.WebserverInfo{Port: test.got_port, IpFqdn: test.got_ipfqdn, Protocol: test.got_proto}, nil)

		ws := WebserverImpl{
			k8sBackend: k8b,
		}
		url := &url.URL{}
		if err := ws.EnrichUrl(url); err != nil {
			t.Errorf("Call to EnrichUrl failed: %v", err)
		}

		if url.Scheme != test.want_scheme {
			t.Errorf("Expected schema of URL to be '%s' got '%s'", test.want_scheme, url.Scheme)
		}

		if url.Host != test.want_host {
			t.Errorf("Expected host portion of URL to be '%s' got '%s'", test.want_host, url.Host)
		}

		mockCtrl.Finish()
	}
}

func TestAddHandler(t *testing.T) {
	m := http.NewServeMux()

	ws := &WebserverImpl{
		mux: m,
	}

	// create a UrlParams struct used to register the handler
	up := structs.NewUrlParams("Nokia", "SRLinux", structs.Config)

	// register the noop handler with the use of the UrlParams
	ws.AddHandler(up, func(w http.ResponseWriter, r *http.Request) { /* NOOP */ })

	// get an instance of net/url from the UrlParams struct
	url := up.GetUrlRelative()
	// it does not contain any host portion, so we have to add that, otherwise
	// the url.String() function will scrumble the whole url
	url.Host = "somehost"

	// create a request on the bases of the url
	req, err := http.NewRequest("GET", url.String(), strings.NewReader("my request"))
	if err != nil {
		t.Errorf("%v", err)
	}

	// figure out from the mux, what would be the handler.
	// since we have added just one, the patter can only reflect the
	// handler path or is "" the empty string if something went wrong
	_, pattern := m.Handler(req)

	// check the pattern is not "" -> the implicit "Page Not Found" handler
	if pattern == "" {
		t.Errorf("Should match the pattern registered")
	}
}
