package interfaces

import (
	"io/fs"
	"net/http"
)

type Storage interface {
	Handle(w http.ResponseWriter, filePath string)
	LoadBackend(fs.FS) error
}
