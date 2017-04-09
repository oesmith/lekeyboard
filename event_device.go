package lekeyboard

import (
	"flag"
	"log"

	evdev "github.com/gvalkov/golang-evdev"
)

var devicePath = flag.String("device", "/dev/input/event0", "Input device")

// TODO: Should probably output events to a channel.

func startEventSource() error {
	dev, err := evdev.Open(*devicePath)
	if err != nil {
		return err
	}
	log.Printf("Opened device:\n%s", dev)
	go eventReadLoop(dev)
	return nil
}

func eventReadLoop(dev *evdev.InputDevice) {
	for {
		e, err := dev.ReadOne()
		if err != nil {
			log.Print(err)
			return
		}
		if e.Type != evdev.EV_KEY || e.Value > 1 {
			continue
		}
		log.Print(e)
	}
}
