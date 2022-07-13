package deviceregistry

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yndd/ztp-webserver/pkg/deviceregistry/interfaces"
	"github.com/yndd/ztp-webserver/pkg/mocks"
)

func TestNewDeviceRegistrySingleton(t *testing.T) {
	dr := GetDeviceRegistry()
	// make sure its a singleton
	if dr != GetDeviceRegistry() {
		t.Errorf("DeviceRegistry is meant to be a singleton. But received a different instance.")
	}
}

func TestRegisterDevice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dr := GetDeviceRegistry()
	devices := []interfaces.RegistryDevice{}

	// add device 1
	mockDevice1 := mocks.NewMockRegistryDevice(mockCtrl)
	devices = append(devices, mockDevice1)
	dr.RegisterDevice(mockDevice1)

	// add device 2
	mockDevice2 := mocks.NewMockRegistryDevice(mockCtrl)
	devices = append(devices, mockDevice2)
	dr.RegisterDevice(mockDevice2)

	// add device 3
	mockDevice3 := mocks.NewMockRegistryDevice(mockCtrl)
	devices = append(devices, mockDevice3)
	dr.RegisterDevice(mockDevice3)

	// check that we get all three back
	entryCount := len(dr.GetRegistryDevices())
	if entryCount < 3 {
		t.Errorf("Added three Devices to registry, got only %d back.", entryCount)
	}

	registryDevices := dr.GetRegistryDevices()
	for _, x := range devices {
		if !deviceInSlice(x, registryDevices) {
			t.Errorf("Missing Device in registry!")
		}
	}
}

func deviceInSlice(a interfaces.RegistryDevice, list []interfaces.RegistryDevice) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
