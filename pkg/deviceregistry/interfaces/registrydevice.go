package interfaces

import webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"

// RegistryDevice is the Interface of instances held in the DeviceRegistry
type RegistryDevice interface {
	SetWebserverSetupper(webserverIf.WebserverSetupper)
}
