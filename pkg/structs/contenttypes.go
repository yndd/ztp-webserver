package structs

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ContentTypes string

const (
	Config         ContentTypes = "config"
	Software       ContentTypes = "software"
	Script         ContentTypes = "script"
	Md5HashFile    ContentTypes = "md5hashfile"
	KernelTim      ContentTypes = "kernel.tim"
	BothTim        ContentTypes = "both.tim"
	IomTim         ContentTypes = "iom.tim"
	SupportTim     ContentTypes = "support.tim"
	CpmTim         ContentTypes = "cpm.tim"
	IsaTim         ContentTypes = "isa.tim"
	HypervisorsTim ContentTypes = "hypervisors.tim"
	BootLdr        ContentTypes = "boot.ldr"
)

var AllContentTypes = []ContentTypes{Config, Software, Script, Md5HashFile, KernelTim, BothTim, IomTim, SupportTim, CpmTim, IsaTim, HypervisorsTim, BootLdr}

// String2ContentTypes converts a regular string into the corresponding ContentType enum
func String2ContentTypes(s string) (*ContentTypes, error) {
	for _, ct := range AllContentTypes {
		if strings.ToLower(s) == string(ct) {
			return &ct, nil
		}
	}
	return nil, fmt.Errorf("no equivalent for %s found in ContentTypes", s)
}

func ContentType2String(ct ContentTypes) string {
	return string(ct)
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (ct *ContentTypes) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	result, err := String2ContentTypes(j)
	if err != nil {
		return err
	}
	*ct = *result
	return nil
}
