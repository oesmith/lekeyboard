package lekeyboard

import (
	"flag"

	"github.com/paypal/gatt"
)

var (
	name = flag.String("name", "LeKeyboard", "Bluetooth device name")

	genericAccessID = gatt.UUID16(0x1800)

	appearanceID = gatt.UUID16(0x2a01)
	deviceNameID = gatt.UUID16(0x2a00)
	peripheralPrivacyID = gatt.UUID16(0x2a02)
	preferredParamsID = gatt.UUID16(0x2a04)
	reconnectionAddressID = gatt.UUID16(0x2a03)

	// HID keyboard (961).
	hidKeyboard = []byte{0x3, 0xc1}
	privacyDisabled = []byte{0}
	// TODO: work out what this means.
	preferredParamsValue = []byte{0x06, 0x00, 0x06, 0x00, 0x00, 0x00, 0xd0, 0x07}
	// TODO: work out what this means.
	reconnectionAddressValue = []byte{0, 0, 0, 0, 0, 0}
)

func NewAccessService() *gatt.Service {
	s := gatt.NewService(genericAccessID)
	s.AddCharacteristic(appearanceID).SetValue(hidKeyboard)
	s.AddCharacteristic(deviceNameID).SetValue([]byte(*name))
	s.AddCharacteristic(peripheralPrivacyID).SetValue(privacyDisabled)
	s.AddCharacteristic(preferredParamsID).SetValue(preferredParamsValue)
	s.AddCharacteristic(reconnectionAddressID).SetValue(reconnectionAddressValue)
	return s
}
