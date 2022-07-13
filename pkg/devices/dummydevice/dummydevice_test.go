package dummydevice

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/yndd/ztp-webserver/pkg/mocks"
)

func TestSetWebserverSetupper(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	dd := GetDummyDevice()
	webserver := mocks.NewMockWebserverSetupper(mockCtrl)
	// the dummy device is expected to register 2 handlers
	webserver.EXPECT().AddHandler(gomock.Any(), gomock.Any()).MinTimes(2)
	dd.SetWebserverSetupper(webserver)
}
