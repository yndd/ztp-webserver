package interfaces

import "github.com/yndd/ztp-dhcp/pkg/backend"

type WebserverOperations interface {
	Run(port int, storageFolder string)
	SetBackend(backend.ZtpBackend)
}
