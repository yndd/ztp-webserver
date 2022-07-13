package interfaces

import (
	"io/fs"
	"net/url"
)

type Index interface {
	DeduceRelativeFilePath(*url.URL) (string, error)
	LoadBackend(fs.FS) error
}
