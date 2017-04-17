package lekeyboard

import (
	"log"

	"github.com/paypal/gatt"
)

var (
  deviceInfoID = gatt.UUID16(0x180a)
  manufacturerNameID = gatt.UUID16(0x2a29)
  pnpID = gatt.UUID16(0x2a50)
  
  manufacturerName = []byte("Olly")
  pnpData = []byte{2, 0, 0, 0, 0, 0, 0} // TODO
)

func NewDeviceInfoService() *gatt.Service {
	s := gatt.NewService(deviceInfoID)
	s.AddCharacteristic(manufacturerNameID).SetValue(manufacturerName)
	s.AddCharacteristic(pnpID).SetValue(pnpData)
	return s
}