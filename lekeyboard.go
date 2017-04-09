package lekeyboard

import (
	"log"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/linux/cmd"
)

func Run() {
	startEventSource()

	d, err := gatt.NewDevice(
		gatt.LnxMaxConnections(1),
		gatt.LnxDeviceID(-1, true),
		gatt.LnxSetAdvertisingParameters(
			&cmd.LESetAdvertisingParameters{
				AdvertisingIntervalMin: 0x00f4,
				AdvertisingIntervalMax: 0x00f4,
				AdvertisingChannelMap:  0x07,
			}))
	if err != nil {
		log.Println("Failed to open device", err)
		return
	}
	d.Handle(
		gatt.CentralConnected(func(c gatt.Central) { log.Println("Connect:", c.ID()) }),
		gatt.CentralDisconnected(func(c gatt.Central) { log.Println("Disconnect:", c.ID()) }),
	)
	d.Init(onStateChanged)
	select {}
}

func onStateChanged(d gatt.Device, s gatt.State) {
	log.Println("State:", s)
	if s != gatt.StatePoweredOn {
		return
	}

	d.AddService(NewGattService())
	d.AddService(NewAccessService())

	b := NewBatteryService()
	d.AddService(b)

	k := NewKeyboardService()
	d.AddService(k)

	d.AdvertiseNameAndServices("PiZero", []gatt.UUID{b.UUID()})
}

