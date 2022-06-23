package storage

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIndexStruct(t *testing.T) {

	i := NewIndex()

	nokiaSrlv001 := i.AddVendor("nokia").AddPlattform("srlinux").AddVersion("v0.0.1")
	nokiaSrlv001.SetFile("nokia/srlinux/v0.0.1_version.txt")
	nokiaSrlv001.SetMd5HashFile("nokia/srlinux/v0.0.1_version_hash.txt")
	nokiaSrlv101 := i.AddVendor("nokia").AddPlattform("srlinux").AddVersion("v1.0.1")
	nokiaSrlv101.SetFile("nokia/srlinux/v1.0.1_version.txt")
	nokiaSrlv101.SetMd5HashFile("nokia/srlinux/v1.0.1_version_hash.txt")
	dummydummyv001 := i.AddVendor("dummy").AddPlattform("dummy").AddVersion("v0.0.1")
	dummydummyv001.SetFile("dummy/dummy/v0.0.1_version.txt")
	dummydummyv001.SetMd5HashFile("dummy/dummy/v0.0.1_version_hash.txt")

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))

}
