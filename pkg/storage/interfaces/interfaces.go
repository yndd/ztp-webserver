package interfaces

import (
	"net/http"
	"net/url"
)

type Storage interface {
	Handle(w http.ResponseWriter, filePath string)
	LoadBackend(base string) error
}

type Index interface {
	DeduceRelativeFilePath(urlPath *url.URL) (string, error)
	LoadBackend(base string) error
}
