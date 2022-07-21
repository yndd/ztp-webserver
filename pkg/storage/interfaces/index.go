package interfaces

import (
	"io/fs"
	"net/url"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

type Index interface {
	DeduceRelativeFilePath(*url.URL) (*structs.FileEntry, error)
	LoadBackend(fs.FS) error
}
