package deviceregistry

import (
	interf "github.com/yndd/ztp-webserver/pkg/deviceregistry/interfaces"
)

var deviceRegistry *DeviceRegistry

// DeviceRegistry is the instance that devices register to on instantiation.
type DeviceRegistry struct {
	devices []interf.RegistryDevice
}

// NewDeviceRegistry constructs a new DeviceRegistry
func newDeviceRegistry() *DeviceRegistry {
	return &DeviceRegistry{
		devices: []interf.RegistryDevice{},
	}
}

// RegisterDevice adds a new RegistryDevice to the registry
func (dr *DeviceRegistry) RegisterDevice(rd interf.RegistryDevice) {
	dr.devices = append(dr.devices, rd)
}

// GetRegistryDevices returns all the registered Devices
func (dr *DeviceRegistry) GetRegistryDevices() []interf.RegistryDevice {
	return dr.devices
}

// GetDeviceRegistry is the method used be the Devices to acquire a handle
// on the DeviceRegistry. It is implemented as a singleton and will therefore
//  always return the pointer to the same instance.
func GetDeviceRegistry() *DeviceRegistry {
	if deviceRegistry != nil {
		return deviceRegistry
	}
	deviceRegistry = newDeviceRegistry()
	return deviceRegistry
}
