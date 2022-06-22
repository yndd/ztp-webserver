package storage

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

type Index struct {
	Vendors map[string]*VendorEntry
}

type VendorEntry struct {
	Plattforms map[string]*PlattformEntry
}

type PlattformEntry struct {
	Versions map[string]*VersionEntry
}

type VersionEntry struct {
	File string `json:"file"`
}

func (i *Index) DeduceRelativeFilePath(urlPath *url.URL) (string, error) {
	purl, err := structs.UrlParamsFromUrl(urlPath)
	if err != nil {
		return "", fmt.Errorf("error parsing url %s - %v", urlPath, err)
	}

	vendor := i.GetVendor(purl.GetVendor())
	if vendor == nil {
		return "", fmt.Errorf("vendor '%s' not foudn in index", purl.GetVendor())
	}
	plattform := vendor.GetPlattform(purl.GetModel())
	if plattform == nil {
		return "", fmt.Errorf("plattform '%s' not found under vendor '%s' in index", purl.GetModel(), purl.GetVendor())
	}
	version := plattform.GetVersion(purl.GetVersion())
	if version == nil {
		return "", fmt.Errorf("version '%s' not found under vendor '%s', plattform '%s' in index", purl.GetVersion(), purl.GetModel(), purl.GetVendor())
	}
	return version.File, nil
}

func (i *Index) AddVendor(vendor string) *VendorEntry {
	// check if vendor already exists
	if ve, exists := i.Vendors[vendor]; exists {
		return ve
	}
	// create a new entry and return it
	newve := &VendorEntry{Plattforms: map[string]*PlattformEntry{}}
	i.Vendors[vendor] = newve
	return newve
}

func (i *Index) GetVendor(vendor string) *VendorEntry {
	// return the vendor entry
	if val, exists := i.Vendors[vendor]; exists {
		return val
	}
	return nil
}

func (ve *VendorEntry) AddPlattform(plattform string) *PlattformEntry {
	// check if plattform already exists
	if ve, exists := ve.Plattforms[plattform]; exists {
		return ve
	}
	// create a new entry and return it
	newpf := &PlattformEntry{Versions: map[string]*VersionEntry{}}
	ve.Plattforms[plattform] = newpf
	return newpf
}

func (ve *VendorEntry) GetPlattform(plattform string) *PlattformEntry {
	// return the plattform entry
	if val, exists := ve.Plattforms[plattform]; exists {
		return val
	}
	return nil
}

func (pe *PlattformEntry) AddVersion(version string) *VersionEntry {
	// check if version already exists
	if ve, exists := pe.Versions[version]; exists {
		return ve
	}
	// create a new entry and return it
	newv := &VersionEntry{}
	pe.Versions[version] = newv
	return newv
}

func (pe *PlattformEntry) GetVersion(version string) *VersionEntry {
	// return the plattform entry
	if val, exists := pe.Versions[version]; exists {
		return val
	}
	return nil
}

func (ve *VersionEntry) SetFile(file string) error {
	if ve.File != "" {
		return fmt.Errorf("file %s is already referenced ignoring %s", ve.File, file)
	}
	ve.File = file

	return nil
}

func NewIndex() *Index {
	return &Index{Vendors: map[string]*VendorEntry{}}
}

func (i *Index) LoadBackend(base string) error {
	finfo, err := os.Stat(base)
	if err != nil {
		return fmt.Errorf("error occured when loading base dir. %v", err)
	}
	if !finfo.IsDir() {
		return fmt.Errorf("error, %s is not a folder", base)
	}
	indexfilename := path.Join(base, "index.yaml")
	_, err = os.Stat(indexfilename)
	if err != nil {
		return fmt.Errorf("error occured when loading index. %v", err)
	}
	dat, err := os.ReadFile(indexfilename)
	if err != nil {
		return fmt.Errorf("error reading index file: %s", err)
	}
	err = json.Unmarshal(dat, i)
	if err != nil {
		return fmt.Errorf("error unmarshalling index file: %v", err)
	}
	return nil
}
