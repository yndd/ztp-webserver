package structs

import "fmt"

type VersionEntry struct {
	Files Files `json:"files"`
}

func NewVersionEntry() *VersionEntry {
	return &VersionEntry{
		Files: Files{},
	}
}

func (ve *VersionEntry) SetFile(ct ContentTypes, fe *FileEntry) error {
	if ve.Files[ct] != nil {
		return fmt.Errorf("file %s is already referenced ignoring %s", ve.Files[ct], fe.String())
	}
	ve.Files[ct] = fe
	return nil
}

func (ve *VersionEntry) GetFile(ct ContentTypes) *FileEntry {
	if ve.Files[ct] == nil {
		return nil
	}
	return ve.Files[ct]
}
