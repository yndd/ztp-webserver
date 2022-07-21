package storage

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"time"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

const indexfilename = "index.json"

type Vendors map[string]structs.Platforms

type Index struct {
	Vendors      Vendors   `json:"vendors"`
	filesystem   fs.FS     `json:"-"`
	indexModTime time.Time `json:"-"`
	indexSize    int64     `json:"-"`
}

func (i *Index) DeduceRelativeFilePath(urlPath *url.URL) (*structs.FileEntry, error) {
	purl, err := structs.UrlParamsFromUrl(urlPath)
	if err != nil {
		return nil, fmt.Errorf("error parsing url %s - %v", urlPath, err)
	}

	// reload index maybe index changed meanwhile
	i.reloadIndex()

	vendor := i.GetVendor(purl.GetVendor())
	if vendor == nil {
		return nil, fmt.Errorf("vendor '%s' not found in index", purl.GetVendor())
	}
	plattform := vendor.GetPlatform(purl.GetModel())
	if plattform == nil {
		return nil, fmt.Errorf("plattform '%s' not found under vendor '%s' in index", purl.GetModel(), purl.GetVendor())
	}
	version := plattform.GetVersion(purl.GetVersion())
	if version == nil {
		return nil, fmt.Errorf("version '%s' not found under vendor '%s', plattform '%s' in index", purl.GetVersion(), purl.GetModel(), purl.GetVendor())
	}

	if filepath, exists := version.Files[purl.GetContentType()]; exists {
		return filepath, nil
	}

	return nil, fmt.Errorf("content not found in index")
}

// AddVendor adds a new vendor to the index if it does not exists
// otherwise return the existing Platforms reference
func (i *Index) AddVendor(vendor string) structs.Platforms {
	// check if vendor already exists
	if ve, exists := i.Vendors[vendor]; exists {
		return ve
	}
	// create a new entry and return it
	newp := structs.Platforms{}
	i.Vendors[vendor] = newp
	return newp
}

func (i *Index) GetVendor(vendor string) structs.Platforms {
	// return the vendor entry
	if val, exists := i.Vendors[vendor]; exists {
		return val
	}
	return nil
}

func NewIndex() *Index {
	return &Index{Vendors: Vendors{}}
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
	err = json.Unmarshal(dat, &i.Vendors)
	if err != nil {
		return fmt.Errorf("error unmarshalling index file: %v", err)
	}

	return nil
}
