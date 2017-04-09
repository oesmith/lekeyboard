package lekeyboard

import (
	"github.com/paypal/gatt"
)

var (
	hidID = gatt.UUID16(0x1812)

	protocolModeID = gatt.UUID16(0x2a4e)

	bootProtocolMode = []byte{0}
)

// TODO: Add a channel for sending reports?
func NewKeyboardService() {
	s := gatt.NewService(hidID)
	s.AddCharacteristic(protocolModeID).setValue(bootProtocolMode)
	// TODO
	return s
}
