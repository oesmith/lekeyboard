package lekeyboard

import (
	"log"

	"github.com/paypal/gatt"
)

var (
	gattID = gatt.UUID16(0x1801)
	serviceChangedID = gatt.UUID16(0x2a05)
)

func NewGattService() *gatt.Service {
	s := gatt.NewService(gattID)
	s.AddCharacteristic(serviceChangedID).HandleNotifyFunc(
		func(r gatt.Request, n gatt.Notifier) {
			log.Println("GATT service changed")
		})
	return s
}
