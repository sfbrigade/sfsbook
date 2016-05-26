package main

import (
	"log"
	"os"

	"github.com/sfbrigade/sfsbook/setup"
)

func main() {
	// TODO(rjk): make the logging configurable in a useful way.
	// TODO(rjk): make the log useful.
	log.Println("sfsbook starting")

	pth, err := os.Getwd()
	if err != nil {
		log.Fatalln("Wow! No CWD. Giving up.", err)
	}

	setup.ConstructNecessaryStartingState(pth)

}
