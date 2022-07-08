package deviceregistry

import webserverIf "github.com/yndd/ztp-webserver/pkg/webserver/interfaces"

var deviceRegistry *DeviceRegistry

// DeviceRegistry is the instance that devices register to on instantiation.
type DeviceRegistry struct {
	devices []RegistryDevice
}

// NewDeviceRegistry constructs a new DeviceRegistry
func NewDeviceRegistry() *DeviceRegistry {
	return &DeviceRegistry{
		devices: []RegistryDevice{},
	}
}

// RegisterDevice adds a new RegistryDevice to the registry
func (dr *DeviceRegistry) RegisterDevice(rd RegistryDevice) {
	dr.devices = append(dr.devices, rd)
}

// GetRegistryDevices returns all the registered Devices
func (dr *DeviceRegistry) GetRegistryDevices() []RegistryDevice {
	return dr.devices
}

// RegistryDevice is the Interface of instances held in the DeviceRegistry
type RegistryDevice interface {
	SetWebserverSetupper(webserverIf.WebserverSetupper)
}

// GetDeviceRegistry is the method used be the Devices to acquire a handle
// on the DeviceRegistry. It is implemented as a singleton and will therefore
//  always return the pointer to the same instance.
func GetDeviceRegistry() *DeviceRegistry {
	if deviceRegistry != nil {
		return deviceRegistry
	}
	deviceRegistry = NewDeviceRegistry()
	return deviceRegistry
}
