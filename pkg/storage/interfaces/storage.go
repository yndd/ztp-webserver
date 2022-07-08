package interfaces

import (
	"net/http"
)

type Storage interface {
	Handle(w http.ResponseWriter, filePath string)
	LoadBackend(base string) error
}
