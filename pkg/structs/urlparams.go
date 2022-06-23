package structs

import (
	"fmt"
	"net/url"
	"strings"
)

type UrlParams struct {
	vendor      string
	model       string
	contentType ContentTypes
	// version will not be present in calls for config or scripts
	// but is necessary for software downloads
	version  string
	deviceId string
}

// Equals check value equality of two UrlParam structs
func (u *UrlParams) Equals(x *UrlParams) bool {
	return u.deviceId == x.deviceId &&
		u.vendor == x.vendor &&
		u.model == x.model &&
		u.version == x.version &&
		u.contentType == x.contentType
}

func (u *UrlParams) SetDeviceId(s string) *UrlParams {
	u.deviceId = s
	return u
}

func (u *UrlParams) SetVersion(v string) *UrlParams {
	u.version = v
	return u
}

// GetUrlRelative retrieve the url string from the UrlParams struct
func (u *UrlParams) GetUrlRelative() *url.URL {

	newUrl := &url.URL{}
	newUrl.Path = fmt.Sprintf("%s/%s/%s", url.PathEscape(u.vendor), url.PathEscape(u.model), url.PathEscape(ContentType2String(u.contentType)))
	newUrl.Scheme = "http"

	q := newUrl.Query()
	// add version parameter if set
	if u.version != "" {
		q.Set("version", u.version)
	}
	// add version parameter if set
	if u.deviceId != "" {
		q.Set("deviceid", u.deviceId)
	}
	newUrl.RawQuery = q.Encode()

	return newUrl
}

// NewUrlParams generate a UrlParams struct without the version parameter
func NewUrlParams(vendor, model string, ct ContentTypes) *UrlParams {
	return &UrlParams{
		vendor:      vendor,
		model:       model,
		contentType: ct,
	}
}

// NewUrlParams generate a UrlParams struct without the version parameter
func NewUrlParamsDeviceId(vendor, model string, deviceId string, ct ContentTypes) *UrlParams {
	return &UrlParams{
		vendor:      vendor,
		model:       model,
		contentType: ct,
		deviceId:    deviceId,
	}
}

func UrlParamsFromUrl(u *url.URL) (*UrlParams, error) {
	var result *UrlParams = nil

	path := u.Path
	if u.Path[0] == '/' {
		path = path[1:]
	}

	splitPath := strings.Split(path, "/")
	if len(splitPath) != 3 {
		return nil, fmt.Errorf("malformed url %s. expected 3 element path", u.String())
	}
	contentType, err := String2ContentTypes(splitPath[2])
	if err != nil {
		return nil, fmt.Errorf("error converting %s to content type(%w)", splitPath[2], err)
	}

	result = NewUrlParams(splitPath[0], splitPath[1], *contentType)
	// parse version
	if val, exists := u.Query()["version"]; exists {
		result.version = val[0]
	}
	// parse device ID
	if val, exists := u.Query()["deviceid"]; exists {
		result.deviceId = val[0]
	}
	return result, err
}

// ParseURL parse the given url and return the corresponding UrlParams struct
func ParseURL(url_str string) (*UrlParams, error) {

	var err error = nil

	u, err := url.Parse(url_str)
	if err != nil {
		return nil, fmt.Errorf("error parsing provided url \"%s\", %w", url_str, err)
	}
	return UrlParamsFromUrl(u)
}

// GetVendor getter for the vendor attribute
func (u *UrlParams) GetVendor() string {
	return u.vendor
}

// GetModel getter for the model attribute
func (u *UrlParams) GetModel() string {
	return u.model
}

// GetContentType getter for the contentType attribute
func (u *UrlParams) GetContentType() ContentTypes {
	return u.contentType
}

// GetVersion getter for the version attribute
func (u *UrlParams) GetVersion() string {
	return u.version
}

// GetDeviceId getter for the deviceId attribute
func (u *UrlParams) GetDeviceId() string {
	return u.deviceId
}
