package interfaces

import (
	"net/http"
	"net/url"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

type Storage interface {
	Handle(w http.ResponseWriter, filePath string)
	LoadBackend(base string) error
}

type Index interface {
	DeduceRelativeFilePath(*url.URL, structs.ContentTypes) (string, error)
	LoadBackend(base string) error
}
