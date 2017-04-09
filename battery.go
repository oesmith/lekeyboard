package lekeyboard

import (
	"github.com/paypal/gatt"
)

var (
	batteryID = gatt.UUID16(0x180f)
	batteryLevelID = gatt.UUID16(0x2a19)
	presentationFormatID = gatt.UUID16(0x2904)
	// Set a constant 100% battery level.
	fullBattery = []byte{100}
	// Describe the battery percentage format.
	// See https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.descriptor.gatt.characteristic_presentation_format.xml
	batteryPercentFormat = []byte{4, 1, 39, 173, 1, 0, 0}
)

func NewBatteryService() *gatt.Service {
	s := gatt.NewService(batteryID)
	c := s.AddCharacteristic(batteryLevelID)
	c.SetValue(fullBattery)
	c.AddDescriptor(presentationFormatID).SetValue(batteryPercentFormat)
	return s
}
