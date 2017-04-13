package lekeyboard

import (
	"github.com/paypal/gatt"
)

const (
	// Protocol modes.
	reportMode = byte(0)
	bootMode = byte(1)
)

var (
	// Assigned IDs
	hidID = gatt.UUID16(0x1812)
	hidInfoID = gatt.UUID16(0x2a4a)
	hidControlPointID = gatt.UUID16(0x2a4c)
	protocolModeID = gatt.UUID16(0x2a4e)

	// HID info data.
	// * HID version: 1.1
	// * Country code: 0
	// * Remote wake: false
	// * Normally connectable: false
	hidInfo = []byte{1,1,0,0}
)

struct KeyboardService {
	protocolMode byte
	service *gatt.Service
}

// TODO: Add a channel for sending reports?
func NewKeyboardService() *KeyboardService {
	ks := &KeyboardService{}
	ks.protocolMode = reportMode
	// TODO: initialize defaults.
	return ks
}

func (ks *KeyboardService) GetService() *gatt.Service {
	if ks.service == nil {
		ks.service = gatt.NewService(hidID)

		c := ks.service.AddCharacteristic(protocolModeID)
		c.HandleReadFunc(
			func(w gatt.ResponseWriter, r *gatt.ReadRequest) {
				if r.Offset > 0 {
					w.SetStatus(7) // Invalid offset
				if _, err := w.Write([]byte{ks.protocolMode}); err != nil {
					log.Println("Protocol mode read failed.", err)
					w.SetStatus(14) // "Unlikely"
					return
				}
				w.SetStatus(0)
			})
		c.HandleWriteFunc(
			func(r gatt.Request, data []byte) {
				if len(data) != 1 || data[0] > 1 {
					log.Printf("Bad protocol mode write: %s\n", data)
					return 6; // Not supported.
				}
				ks.protocolMode = data[0];
				log.Printf("Protocol mode changed: %d\n", ks.protocolMode)
				return 0;
			})

		// TODO: report characteristics (input/output).
		// TODO: - client characteristic descriptors (input/output).
		// TODO: - report reference descriptors (input/output).

		// TODO: report map characteristic.

		// TODO: boot keyboard input report characteristic.
		// TODO: - client characteristic descriptor.

		// TODO: boot keyboard output report characteristic.
		
		ks.service.AddCharacteristic(hidInfoID).SetValue(hidInfo)

		c = ks.service.AddCharacteristic(hidControlPointID)
		c.HandleWriteFunc(
			func(r gatt.Request, data []byte) byte {
				if len(data) != 1 || data[0] > 1 {
					log.Printf("Bad HID control point request: %s\n", data)
					return 6; // Not supported
				}
				switch data[0] {
				case 0:
					ks.HandleSuspend()
				case 1:
					ks.HandleExitSuspend()
				}
				return 0;
			})
	}
	return ks.service
}

func (ks *KeyboardService) HandleSuspend() {
	log.Println("Suspend")
}

func (ks *KeyboardService) HandleExitSuspend() {
	log.Println("Exit suspend")
}