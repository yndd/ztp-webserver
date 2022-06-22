package storage

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/utils"
)

type FolderStorageImpl struct {
	base string
}

func NewFolderStorage() *FolderStorageImpl {
	return &FolderStorageImpl{}
}

func (fs *FolderStorageImpl) LoadBackend(base string) error {
	finfo, err := os.Stat(base)
	if err != nil {
		return fmt.Errorf("error occured when loading base dir. %v", err)
	}
	if !finfo.IsDir() {
		return fmt.Errorf("error, %s is not a folder", base)
	}
	fs.base = base
	return nil
}

// DeliverFile delivers the content of filePath to the call via the http.RespinseWriter
func (fs *FolderStorageImpl) Handle(w http.ResponseWriter, filePath string) {

	// lets check that only files within the base directory are being delivered.
	abspath := path.Join(fs.base, filePath)
	relbasepath, err := filepath.Rel(fs.base, abspath)
	if err != nil {
		log.Errorf("error deducing relative path: %v", err)
	}
	// if the path starts with .. we would be leaving the base directory -> Error
	if strings.HasPrefix(relbasepath, "..") {
		log.Errorf("requested path '%s' is outside of the base '%s' directory", filePath, fs.base)
		utils.HandleErrorCode(500, "requested path is outside of the base directory", w)
		return
	}

	// read the file and deliver it finally
	data, err := ioutil.ReadFile(abspath)
	if err != nil {
		// show error page if failed to read file
		log.Errorf("error reading file '%s': %v", filePath, err)
		utils.HandleErrorCode(500, "Unable to retrieve file", w)
		return
	} else {
		if _, err = w.Write(data); err != nil {
			log.Errorf("error when delivering file: %v", err)
			return
		}
		log.Infof("%s delivered successfully", filePath)
	}
}
