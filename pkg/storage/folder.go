package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/webserver"
)

type FolderStorageImpl struct {
	base string
}

func NewFolderStorage(base string) (*FolderStorageImpl, error) {
	fileinfo, err := os.Stat(base)
	if err != nil {
		return nil, fmt.Errorf("error initializing FolderStorage. %v ", err)
	}
	if !fileinfo.IsDir() {
		return nil, fmt.Errorf("error initializing FolderStorage. %s is not a file", base)
	}
	return &FolderStorageImpl{base: base}, nil
}

// DeliverFile delivers the content of filePath to the call via the http.RespinseWriter
func (sf *FolderStorageImpl) Handle(w http.ResponseWriter, filePath string) {

	// lets check that only files within the base directory are being delivered.
	relbasepath, err := filepath.Rel(sf.base, filePath)
	if err != nil {
		log.Errorf("error deducing relative path: %v", err)
	}
	// if the path starts with .. we would be leaving the base directory -> Error
	if strings.HasPrefix(relbasepath, "..") {
		log.Errorf("requested path is outside of the base directory")
		webserver.HandleErrorCode(500, "requested path is outside of the base directory", w)
		return
	}

	// read the file and deliver it finally
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// show error page if failed to read file
		webserver.HandleErrorCode(500, "Unable to retrieve file", w)
		return
	} else {
		if _, err = w.Write(data); err != nil {
			log.Error("error when delivering file: %v", err)
			return
		}
		log.Info("%s delivered successfully", filePath)
	}
}
