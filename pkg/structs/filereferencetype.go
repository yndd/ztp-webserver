package structs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type FileReferenceType string

const (
	Filesystem   FileReferenceType = "filesystem"
	HTTPRedirect FileReferenceType = "httpredirect"
)

var AllFileReferenceTypes = []FileReferenceType{Filesystem, HTTPRedirect}

// String2FileReferenceType converts a regular string into the corresponding FileReferenceType enum
func String2FileReferenceType(s string) (*FileReferenceType, error) {
	for _, ct := range AllFileReferenceTypes {
		if strings.ToLower(s) == string(ct) {
			return &ct, nil
		}
	}
	return nil, fmt.Errorf("no equivalent for %s found in FileReferenceType", s)
}

func FileReferenceType2String(ct ContentTypes) string {
	return string(ct)
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (frt *FileReferenceType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	result, err := String2FileReferenceType(j)
	if err != nil {
		return err
	}
	*frt = *result
	return nil
}
