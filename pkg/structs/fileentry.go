package structs

import (
	"fmt"
)

type FileEntry struct {
	ReferenceType FileReferenceType `json:"type"`
	Reference     string            `json:"ref"`
}

// String returns a string representation of the FileEntry
func (fe *FileEntry) String() string {
	return fmt.Sprintf("[%s] %s", string(fe.ReferenceType), fe.Reference)
}

// NewFileEntry constructs a new FileEntry
func NewFileEntry(rt FileReferenceType, ref string) *FileEntry {
	return &FileEntry{
		ReferenceType: rt,
		Reference:     ref,
	}
}

func (fe *FileEntry) GetReference() string {
	return fe.Reference
}

func (fe *FileEntry) GetReferenceType() FileReferenceType {
	return fe.ReferenceType
}
