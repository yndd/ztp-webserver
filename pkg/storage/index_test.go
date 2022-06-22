package storage

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestIndexStruct(t *testing.T) {

	i := NewIndex()

	i.AddVendor("nokia").AddPlattform("srlinux").AddVersion("v0.0.1").SetFile("nokia/srlinux/v0.0.1_version.txt")
	i.AddVendor("nokia").AddPlattform("srlinux").AddVersion("v1.0.1").SetFile("nokia/srlinux/v1.0.1_version.txt")
	i.AddVendor("dummy").AddPlattform("dummy").AddVersion("v0.0.1").SetFile("dummy/dummy/v0.0.1_version.txt")

	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))

}
