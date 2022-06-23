package structs

import (
	"fmt"
	"strings"
)

type ContentTypes string

const (
	Config      ContentTypes = "config"
	Software    ContentTypes = "software"
	Script      ContentTypes = "script"
	Md5HashFile ContentTypes = "md5hashfile"
)

// String2ContentTypes converts a regular string into the corresponding ContentType enum
func String2ContentTypes(s string) (*ContentTypes, error) {
	var result ContentTypes
	var err = fmt.Errorf("no equivalent for %s found in ContentTypes", s)
	switch strings.ToLower(s) {
	case "config":
		result = Config
		err = nil
	case "software":
		result = Software
		err = nil
	case "script":
		result = Script
		err = nil
	case "md5hashfile":
		result = Md5HashFile
		err = nil
	}
	return &result, err
}

func ContentType2String(ct ContentTypes) string {
	return string(ct)
}
