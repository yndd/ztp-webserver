package storage

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"time"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

const indexfilename = "index.yaml"

type Index struct {
	Vendors      map[string]*VendorEntry `json:"vendors"`
	filesystem   fs.FS                   `json:"-"`
	indexModTime time.Time               `json:"-"`
	indexSize    int64                   `json:"-"`
}

type VendorEntry struct {
	Platforms map[string]*PlatformEntry `json:"platforms"`
}

type PlatformEntry struct {
	Versions map[string]*VersionEntry `json:"versions"`
}

type VersionEntry struct {
	File        string `json:"file"`
	Md5HashFile string `json:"md5HashFile,omitempty"`
}

func (i *Index) DeduceRelativeFilePath(urlPath *url.URL) (string, error) {
	purl, err := structs.UrlParamsFromUrl(urlPath)
	if err != nil {
		return "", fmt.Errorf("error parsing url %s - %v", urlPath, err)
	}

	// reload index maybe index changed meanwhile
	i.reloadIndex()

	vendor := i.GetVendor(purl.GetVendor())
	if vendor == nil {
		return "", fmt.Errorf("vendor '%s' not found in index", purl.GetVendor())
	}
	plattform := vendor.GetPlattform(purl.GetModel())
	if plattform == nil {
		return "", fmt.Errorf("plattform '%s' not found under vendor '%s' in index", purl.GetModel(), purl.GetVendor())
	}
	version := plattform.GetVersion(purl.GetVersion())
	if version == nil {
		return "", fmt.Errorf("version '%s' not found under vendor '%s', plattform '%s' in index", purl.GetVersion(), purl.GetModel(), purl.GetVendor())
	}

	switch purl.GetContentType() {
	case structs.Software:
		return version.File, nil
	case structs.Md5HashFile:
		return version.Md5HashFile, nil
	}

	return "", fmt.Errorf("content not found in index")
}

func (i *Index) AddVendor(vendor string) *VendorEntry {
	// check if vendor already exists
	if ve, exists := i.Vendors[vendor]; exists {
		return ve
	}
	// create a new entry and return it
	newve := &VendorEntry{Platforms: map[string]*PlatformEntry{}}
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

func (ve *VendorEntry) AddPlattform(plattform string) *PlatformEntry {
	// check if plattform already exists
	if ve, exists := ve.Platforms[plattform]; exists {
		return ve
	}
	// create a new entry and return it
	newpf := &PlatformEntry{Versions: map[string]*VersionEntry{}}
	ve.Platforms[plattform] = newpf
	return newpf
}

func (ve *VendorEntry) GetPlattform(plattform string) *PlatformEntry {
	// return the plattform entry
	if val, exists := ve.Platforms[plattform]; exists {
		return val
	}
	return nil
}

func (pe *PlatformEntry) AddVersion(version string) *VersionEntry {
	// check if version already exists
	if ve, exists := pe.Versions[version]; exists {
		return ve
	}
	// create a new entry and return it
	newv := &VersionEntry{}
	pe.Versions[version] = newv
	return newv
}

func (pe *PlatformEntry) GetVersion(version string) *VersionEntry {
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

func (ve *VersionEntry) SetMd5HashFile(file string) error {
	if ve.Md5HashFile != "" {
		return fmt.Errorf("md5 hash file %s is already referenced ignoring %s", ve.Md5HashFile, file)
	}
	ve.Md5HashFile = file

	return nil
}

func NewIndex() *Index {
	return &Index{Vendors: map[string]*VendorEntry{}}
}

func (i *Index) LoadBackend(filesystem fs.FS) error {
	i.filesystem = filesystem
	return i.reloadIndex()
}

func (i *Index) reloadIndex() error {
	istat, err := fs.Stat(i.filesystem, indexfilename)
	if err != nil {
		return fmt.Errorf("error occured when loading index. %v", err)
	}
	// check if index changed
	if istat.ModTime() == i.indexModTime && istat.Size() == int64(i.indexSize) {
		// nothing changed, so no reaload required, return to caller
		return nil
	}

	// index changed, not the actual time and size
	i.indexModTime = istat.ModTime()
	i.indexSize = istat.Size()

	// read the index file
	dat, err := fs.ReadFile(i.filesystem, indexfilename)
	if err != nil {
		return fmt.Errorf("error reading index file: %s", err)
	}
	// unmarshal it from json
	err = json.Unmarshal(dat, i)
	if err != nil {
		return fmt.Errorf("error unmarshalling index file: %v", err)
	}
	return nil
}
