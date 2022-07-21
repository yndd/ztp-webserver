package storage

import (
	"fmt"
	iofs "io/fs"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/yndd/ztp-webserver/pkg/utils"
)

type FolderStorageImpl struct {
	filesystem iofs.FS
}

func NewFolderStorage() *FolderStorageImpl {
	return &FolderStorageImpl{}
}

func (fs *FolderStorageImpl) LoadBackend(filesys iofs.FS) error {
	fs.filesystem = filesys
	return nil
}

// DeliverFile delivers the content of filePath to the call via the http.RespinseWriter
func (fs *FolderStorageImpl) Handle(w http.ResponseWriter, filePath string) {

	// read the file and deliver it finally
	data, err := fs.GetFileContent(filePath)
	if err != nil {
		// show error page if failed to read file
		utils.HandleErrorCodeLog(500, fmt.Errorf("unable to retrieve file: %v", err), w)
		return
	} else {
		if _, err = w.Write(data); err != nil {
			utils.HandleErrorCodeLog(500, fmt.Errorf("unable to deliver file: %v", err), w)
			return
		}
		log.Infof("%s delivered successfully", filePath)
	}
}

//GetFileContent retrieve the data from the given file
func (fs *FolderStorageImpl) GetFileContent(filePath string) ([]byte, error) {
	return iofs.ReadFile(fs.filesystem, filePath)
}
