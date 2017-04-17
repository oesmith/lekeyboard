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

	bootKeyboardInputReportID = gatt.UUID16(0x2a22)
	bootKeyboardOutputReportID = gatt.UUID16(0x2a32)
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
	protocolMode []byte
	inputClientConfig []byte
	bootInputClientConfig []byte
	inputReport []byte
	outputReport []byte

	service *gatt.Service
}

func makeReadFunc(data []byte, name string) func(w gatt.ResponseWriter, r *gatt.ReadRequest) {
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
	}
}

func makeWriteFunc(buffer []byte, name string) func(gatt.Request, []byte) byte {
	return func(r gatt.Request, data []byte) byte {
		if len(data) != len(buffer) {
			log.Printf("Bad write: %s\n", data)
			return 6; // Not supported.
		}
		copy(buffer, data)
		log.Printf("Write %s: %s\n", name, data)
		return 0;
	};
}

// TODO: Add a channel for sending reports?
func NewKeyboardService() *KeyboardService {
	ks := &KeyboardService{}
	ks.Reset()
	return ks
}

func (ks *KeyboardService) Reset() {
	ks.protocolMode = make([]byte, 1)

	ks.inputClientConfig = make([]byte, 2)
	ks.bootInputClientConfig = make([]byte, 2)
	ks.inputReport = make([]byte, 8)

	ks.outputReport = make([]byte, 1)
}

func (ks *KeyboardService) GetService() *gatt.Service {
	if ks.service == nil {
		ks.service = gatt.NewService(hidID)

		c := ks.service.AddCharacteristic(protocolModeID)
		c.HandleReadFunc(makeReadFunc(ks.protocolMode, "Protocol Mode"))
		c.HandleWriteFunc(makeWriteFunc(ks.protocolMode, "Protocol Mode"))
		// TODO: events on write.

		c = ks.service.AddCharacteristic(reportID)
		c.HandleReadFunc(makeReadFunc(ks.inputReport, "Input Report"))
		c.AddDescriptor(reportReferenceID).SetValue(inputReportRef)
		d := c.AddDescriptor(clientCharacteristicID)
		d.HandleReadFunc(makeReadFunc(ks.inputClientConfig, "Input Client Config"))
		d.HandleWriteFunc(makeWriteFunc(ks.inputClientConfig, "Input Client Config"))
		// TODO: events on write.

		c = ks.service.AddCharacteristic(reportID)
		c.HandleReadFunc(makeReadFunc(ks.outputReport, "Output Report"))
		c.HandleWriteFunc(makeWriteFunc(ks.outputReport, "Output Report"))
		// TODO: events on write.
		c.AddDescriptor(reportReferenceID).SetValue(outputReportRef)

		ks.service.AddCharacteristic(reportMapID).SetValue(hidDescriptor)

		c = ks.service.AddCharacteristic(bootKeyboardInputReportID)
		c.HandleReadFunc(makeReadFunc(ks.inputReport, "Boot Input Report"))
		d = c.AddDescriptor(clientCharacteristicID)
		d.HandleReadFunc(makeReadFunc(ks.bootInputClientConfig, "Boot Input Client Config"))
		d.HandleWriteFunc(makeWriteFunc(ks.bootInputClientConfig, "Boot Input Client Config"))
		// TODO: events on write.

		c = ks.service.AddCharacteristic(bootKeyboardOutputReportID)
		c.HandleReadFunc(makeReadFunc(ks.outputReport, "Boot Output Report"))
		c.HandleWriteFunc(makeWriteFunc(ks.outputReport, "Boot Output Report"))
		// TODO: events on write.
		
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
