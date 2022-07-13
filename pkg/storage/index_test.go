package storage

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIndexStruct(t *testing.T) {
	i := NewIndex()
	NokiaSrlv001 := i.AddVendor("Nokia").AddPlattform("SRLinux").AddVersion("v0.0.1")
	NokiaSrlv001.SetFile("Nokia/SRLinux/v0.0.1_version.txt")
	NokiaSrlv001.SetMd5HashFile("Nokia/SRLinux/v0.0.1_version_hash.txt")
	NokiaSrlv101 := i.AddVendor("Nokia").AddPlattform("SRLinux").AddVersion("v1.0.1")
	NokiaSrlv101.SetFile("Nokia/SRLinux/v1.0.1_version.txt")
	NokiaSrlv101.SetMd5HashFile("Nokia/SRLinux/v1.0.1_version_hash.txt")
	dummydummyv001 := i.AddVendor("dummy").AddPlattform("dummy").AddVersion("v0.0.1")
	dummydummyv001.SetFile("dummy/dummy/v0.0.1_version.txt")
	dummydummyv001.SetMd5HashFile("dummy/dummy/v0.0.1_version_hash.txt")

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))

}

func TestDeduceRelativeFilePath(t *testing.T) {

}
