package structs

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	versionPrefix  = "version"
	deviceIdPrefix = "deviceid"
)

type UrlParams struct {
	vendor      string
	model       string
	contentType ContentTypes
	// version will not be present in calls for config or scripts
	// but is necessary for software downloads
	version    string
	devicename string
	filename   string
}

// Equals check value equality of two UrlParam structs
func (u *UrlParams) Equals(x *UrlParams) bool {
	return u.devicename == x.devicename &&
		u.vendor == x.vendor &&
		u.model == x.model &&
		u.version == x.version &&
		u.contentType == x.contentType &&
		u.filename == x.filename
}

func (u *UrlParams) SetDeviceId(s string) *UrlParams {
	u.devicename = s
	return u
}

func (u *UrlParams) SetVersion(v string) *UrlParams {
	u.version = v
	return u
}

func (u *UrlParams) SetFilename(fn string) *UrlParams {
	u.filename = fn
	return u
}

// GetUrlRelative retrieve the url string from the UrlParams struct
func (u *UrlParams) GetUrlRelative() *url.URL {

	newUrl := &url.URL{}
	newUrl.Path = fmt.Sprintf("%s/%s/%s", url.PathEscape(u.vendor), url.PathEscape(u.model), url.PathEscape(ContentType2String(u.contentType)))

	// add version parameter if set
	if u.version != "" {
		newUrl.Path = fmt.Sprintf("%s/%s", newUrl.Path, url.PathEscape(fmt.Sprintf("%s=%s", versionPrefix, u.version)))
	}
	// add version parameter if set
	if u.devicename != "" {
		newUrl.Path = fmt.Sprintf("%s/%s", newUrl.Path, url.PathEscape(fmt.Sprintf("%s=%s", deviceIdPrefix, u.devicename)))
	}

	if u.filename != "" {
		newUrl.Path = fmt.Sprintf("%s/%s", newUrl.Path, url.PathEscape(u.filename))
	}

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
func NewUrlParamsDeviceName(vendor, model string, deviceId string, ct ContentTypes) *UrlParams {
	return &UrlParams{
		vendor:      vendor,
		model:       model,
		contentType: ct,
		devicename:  deviceId,
	}
}

// UrlParamsFromUrl takes a given URL and converts it into the URLParams object
func UrlParamsFromUrl(u *url.URL) (*UrlParams, error) {
	var result *UrlParams = nil

	path := u.Path
	if u.Path[0] == '/' {
		path = path[1:]
	}

	splitPath := strings.Split(path, "/")
	contentType, err := String2ContentTypes(splitPath[2])
	if err != nil {
		return nil, fmt.Errorf("error converting %s to content type(%w)", splitPath[2], err)
	}

	result = NewUrlParams(splitPath[0], splitPath[1], *contentType)

	for _, x := range splitPath {
		lowerX := strings.ToLower(x)
		switch {
		case strings.HasPrefix(lowerX, versionPrefix+"="):
			result.version = x[len(versionPrefix)+1:]
		case strings.HasPrefix(lowerX, deviceIdPrefix+"="):
			result.devicename = x[len(deviceIdPrefix)+1:]
		}
	}

	if !strings.Contains(splitPath[len(splitPath)-1], "=") {
		result.filename = splitPath[len(splitPath)-1]
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

// GetDeviceName getter for the deviceId attribute
func (u *UrlParams) GetDeviceName() string {
	return u.devicename
}

// GetFilename getter for the deviceId attribute
func (u *UrlParams) GetFilename() string {
	return u.filename
}
