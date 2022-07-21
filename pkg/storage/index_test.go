package storage

import (
	"encoding/json"
	"fmt"
	"testing"
	"testing/fstest"

	"github.com/yndd/ztp-webserver/pkg/structs"
)

func TestIndexStruct(t *testing.T) {
	i := NewIndex()
	NokiaSrlv001 := i.AddVendor("Nokia").AddPlatform("SRLinux").AddVersion("v0.0.1")
	NokiaSrlv001.SetFile(structs.Software, structs.NewFileEntry(structs.Filesystem, "Nokia/SRLinux/v0.0.1_version.txt"))
	NokiaSrlv001.SetFile(structs.Md5HashFile, structs.NewFileEntry(structs.Filesystem, "Nokia/SRLinux/v0.0.1_version_hash.txt"))

	NokiaSrlv101 := i.AddVendor("Nokia").AddPlatform("SRLinux").AddVersion("v1.0.1")
	NokiaSrlv101.SetFile(structs.Software, structs.NewFileEntry(structs.Filesystem, "Nokia/SRLinux/v1.0.1_version.txt"))
	NokiaSrlv101.SetFile(structs.Md5HashFile, structs.NewFileEntry(structs.Filesystem, "Nokia/SRLinux/v1.0.1_version_hash.txt"))

	dummydummyv001 := i.AddVendor("dummy").AddPlatform("dummy").AddVersion("v0.0.1")
	dummydummyv001.SetFile(structs.Software, structs.NewFileEntry(structs.Filesystem, "dummy/dummy/v0.0.1_version.txt"))
	dummydummyv001.SetFile(structs.Md5HashFile, structs.NewFileEntry(structs.Filesystem, "dummy/dummy/v0.0.1_version_hash.txt"))

	data, err := json.MarshalIndent(i.Vendors, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))
}

func TestLoadBackendReloadIndex(t *testing.T) {
	i := NewIndex()

	testdata := `
	{
		"Dummy":{
		   "Dummy":{
			  "v0.0.1":{
				 "files":{
					"software":{
					   "type":"filesystem",
					   "ref":"dummy/dummy/v0.0.1_version.txt"
					},
					"md5HashFile":{
					   "type":"filesystem",
					   "ref":"dummy/dummy/v0.0.1_version_hash.txt"
					}
				 }
			  }
		   }
		},
		"Nokia":{
		   "SRLinux":{
			  "v0.0.1":{
				 "files":{
					"software":{
					   "type":"filesystem",
					   "ref":"nokia/srlinux/v0.0.1_version.txt"
					},
					"md5HashFile":{
					   "type":"filesystem",
					   "ref":"nokia/srlinux/v0.0.1_version_hash.txt"
					}
				 }
			  },
			  "v1.0.1":{
				 "files":{
					"software":{
					   "type":"filesystem",
					   "ref":"nokia/srlinux/v1.0.1_version.txt"
					},
					"md5HashFile":{
					   "type":"filesystem",
					   "ref":"nokia/srlinux/v1.0.1_version_hash.txt"
					}
				 }
			  }
		   },
		   "SROS":{
			  "22.5.R1":{
				 "files":{
					"kernel.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/kernel.tim"
					},
					"both.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/both.tim"
					},
					"support.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/support.tim"
					},
					"iom.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/iom.tim"
					},
					"cpm.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/cpm.tim"
					},
					"boot.ldr":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/boot.ldr"
					},
					"hypervisor.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R1/hypervisor.tim"
					}
				 }
			  },
			  "22.5.R2":{
				 "files":{
					"kernel.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/kernel.tim"
					},
					"both.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/both.tim"
					},
					"support.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/support.tim"
					},
					"iom.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/iom.tim"
					},
					"cpm.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/cpm.tim"
					},
					"boot.ldr":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/boot.ldr"
					},
					"hypervisor.tim":{
					   "type":"filesystem",
					   "ref":"nokia/sros/22.5.R2/hypervisor.tim"
					}
				 }
			  }
		   }
		}
	  }
	`
	mfs := &fstest.MapFS{"index.json": {Data: []byte(testdata)}}

	err := i.LoadBackend(mfs)
	if err != nil {
		t.Error(err)
	}

	data, err := json.MarshalIndent(i.Vendors, "", "  ")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(data))

}
