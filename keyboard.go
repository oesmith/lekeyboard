package lekeyboard

import (
	"log"

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

	clientCharacteristicID = gatt.UUID16(0x2902)
	reportReferenceID = gatt.UUID16(0x2908)

	hidInfoID = gatt.UUID16(0x2a4a)
	reportMapID = gatt.UUID16(0x2a4b)
	hidControlPointID = gatt.UUID16(0x2a4c)
	reportID = gatt.UUID16(0x2a4d)
	protocolModeID = gatt.UUID16(0x2a4e)

	// HID info data.
	// * HID version: 1.1
	// * Country code: 0
	// * Remote wake: false
	// * Normally connectable: false
	hidInfo = []byte{1,1,0,0}

	// Report descriptor (default HID keyboard -- equivalent to boot protocol).
	hidDescriptor = []byte{
		0x05, 0x01, 0x09, 0x06, 0xa1, 0x01, 0x85, 0x01, 0x75, 0x01, 0x95, 0x08,
		0x05, 0x07, 0x19, 0xe0, 0x29, 0xe7, 0x15, 0x00, 0x25, 0x01, 0x81, 0x02,
		0x95, 0x01, 0x75, 0x08, 0x81, 0x03, 0x95, 0x05, 0x75, 0x01, 0x05, 0x08,
		0x19, 0x01, 0x29, 0x05, 0x91, 0x02, 0x95, 0x01, 0x75, 0x03, 0x91, 0x03,
		0x95, 0x06, 0x75, 0x08, 0x15, 0x00, 0x26, 0xff, 0x00, 0x05, 0x07, 0x19,
		0x00, 0x29, 0xff, 0x81, 0x00, 0xc0,
	}

	// Report references.
	inputReportRef = []byte{0x1, 0x1}
	outputReportRef = []byte{0x1, 0x2}
)

type KeyboardService struct {
	protocolMode byte
	inputClientConfig [2]byte
	outputClientConfig [2]byte
	inputReport [8]byte
	outputReport [1]byte

	service *gatt.Service
}

func makeReadFunc(data []byte, name string) func(gatt.ResponseWriter, *gatt.Request) {
	return func(w gatt.ResponseWriter, r *gatt.ReadRequest) {
		if r.Offset >= len(data) {
			log.Print("Invalid offset reading %s", name)
			w.SetStatus(7) // Invalid offset
			return
		}
		if _, err := w.Write(data[r.Offset:]); err != nil {
			log.Println("Error reading " + name, err)
			w.SetStatus(14) // "Unlikely"
			return
		}
		w.SetStatus(0)
	};
}

func makeWriteFunc(data []byte, name string) func(gatt.Request, []byte) byte {
	return func(r gatt.Request, data []byte) byte {
		if len(data) != len(data) {
			log.Printf("Bad write: %s\n", data)
			return 6; // Not supported.
		}
		ks.protocolMode = data[0];
		log.Printf("Protocol mode changed: %d\n", ks.protocolMode)
		return 0;
	};
}

func (ks *KeyboardService) Reset() {
	ks.protocolMode = 0
	for i := range(2) {
		ks.inputClientConfig[i] = 0
		ks.outputClientConfig[i] = 0
	}
	for i := range(8) {
		ks.inputReport[i] = 0
	}
	ks.outputReport[0] = 0
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
					return
				}
				if _, err := w.Write([]byte{ks.protocolMode}); err != nil {
					log.Println("Protocol mode read failed.", err)
					w.SetStatus(14) // "Unlikely"
					return
				}
				w.SetStatus(0)
			})
		c.HandleWriteFunc(
			func(r gatt.Request, data []byte) byte {
				if len(data) != 1 || data[0] > 1 {
					log.Printf("Bad protocol mode write: %s\n", data)
					return 6; // Not supported.
				}
				ks.protocolMode = data[0];
				log.Printf("Protocol mode changed: %d\n", ks.protocolMode)
				return 0;
			})

		c := ks.service.AddCharacteristic(reportID)
		d := c.AddDescriptor(reportReferenceID).SetValue(inputReportRef)
		d := c.AddDescriptor(clientCharacteristicID)
		d.HandleReadFunc(
			func(w gatt.ResponseWriter, r *gatt.ReadRequest) {
				if r.Offset > 1 {
					w.SetStatus(7) // Invalid offset
					return
				}
				if _, err := w.Write(ks.inputClientConfig[r.Offset:]); err != nil {
					log.Println("Input client config read failed.", err)
					w.SetStatus(14) // Unlikely
				}
				w.SetStatus(0)
			})
		c.HandleWriteFunc(
			func(r gatt.Request, data []byte) byte {
				if len(data) != 2 {
					log.Printf("Bad input client config write: %s\n", data)
					return 6; // Not supported.
				}
				copy(ks.inputClientConfig, data)
				log.Printf("Input client config changed: %d\n", ks.protocolMode)
				return 0;
			})
		c.HandleReadFunc(
			func(w gatt.ResponseWriter, r *gatt.ReadRequest) {
				if r.Offset > 7 {
					w.SetStatus(7)
					return
				}
				if _, err := w.Write(ks.inputReport[r.Offset:]); err != nil {
					log.Println("Input report read failed.", err)
					w.SetStatus(14)
					return
				}
				w.SetStatus(0)
			})

		c := ks.service.AddCharacteristic(reportID)
		d := c.AddDescriptor(reportReferenceID).SetValue(outputReportRef)
		d := c.AddDescriptor(clientCharacteristicID)
		d.HandleReadFunc(
			func(w gatt.ResponseWriter, r *gatt.ReadRequest) {
				if r.Offset > 1 {
					w.SetStatus(7) // Invalid offset
					return
				}
				if _, err := w.Write(ks.outputClientConfig[r.Offset:]); err != nil {
					log.Println("Output client config read failed.", err)
					w.SetStatus(14) // Unlikely
				}
				w.SetStatus(0)
			})
		c.HandleWriteFunc(
			func(r gatt.Request, data []byte) byte {
				if len(data) != 2 {
					log.Printf("Bad output client config write: %s\n", data)
					return 6; // Not supported.
				}
				copy(ks.outputClientConfig, data)
				log.Printf("Output client config changed: %d\n", ks.protocolMode)
				return 0;
			})

		// TODO: report characteristics (input/output).
		// TODO: - client characteristic descriptors (input/output).
		// TODO: - report reference descriptors (input/output).

		c := ks.service.AddCharacteristic(reportMapID).SetValue(hidDescriptor)

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
