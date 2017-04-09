package main

import (
	"flag"

	kb "github.com/oesmith/lekeyboard"
)

func main() {
	flag.Parse()
	kb.Run()
}
