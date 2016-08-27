package main

import (
	"flag"
	"log"
	"os"

	"github.com/sfbrigade/sfsbook/dba"
	"github.com/sfbrigade/sfsbook/server"
	"github.com/sfbrigade/sfsbook/setup"
	"github.com/sfbrigade/sfsbook/dba/fieldmap"
)

func main() {
	flag.Parse()

	// TODO(rjk): make the logging configurable in a useful way.
	// TODO(rjk): make the log useful.
	log.Println("sfsbook starting")

	pth, err := os.Getwd()
	if err != nil {
		log.Fatalln("Wow! No CWD. Giving up.", err)
	}

	setup.ConstructNecessaryStartingState(pth)
	index, err := dba.OpenBleve(pth, fieldmap.RefGuide)
	if err != nil {
		log.Fatalln("No database! Giving up:", err)
	}

	srv := server.MakeServer(":10443", pth, index)
	server.Start(pth, srv)

}
