package interfaces

import (
	"net/url"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

type Index interface {
	DeduceRelativeFilePath(*url.URL, structs.ContentTypes) (string, error)
	LoadBackend(base string) error
}
