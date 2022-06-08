package structs

import (
	"fmt"
)

type UrlParams struct {
	vendor      string
	model       string
	contentType ContentTypes
}

func (u *UrlParams) GetUrlRelative() string {
	return fmt.Sprintf("%s/%s/%s", u.vendor, u.model, u.contentType)
}

func NewUrlParams(vendor, model string, ct ContentTypes) *UrlParams {
	return &UrlParams{
		vendor:      vendor,
		model:       model,
		contentType: ct,
	}
}
