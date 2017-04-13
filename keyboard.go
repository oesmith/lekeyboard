package lekeyboard

import (
	"github.com/paypal/gatt"
)

var (
	hidID = gatt.UUID16(0x1812)

	protocolModeID = gatt.UUID16(0x2a4e)
)

struct KeyboardService {
	bootProtocolMode [1]byte

	service *gatt.Service
}

// TODO: Add a channel for sending reports?
func NewKeyboardService() *KeyboardService {
	ks := &KeyboardService{}
	// TODO: initialize defaults.
	return ks
}

func (ks *KeyboardService) GetService() *gatt.Service {
	if ks.service == nil {
		ks.service = gatt.NewService(hidID)
		// TODO: protocol mode characteristic.
		// TODO: report characteristics (input/output).
		// TODO: - client characteristic descriptors (input/output).
		// TODO: - report reference descriptors (input/output).
		// TODO: report map characteristic.
		// TODO: boot keyboard input report characteristic.
		// TODO: - client characteristic descriptor.
		// TODO: boot keyboard output report characteristic.
		// TODO: hid information characteristic.
		// TODO: hid control point characteristic.
	}
	return ks.service
}
