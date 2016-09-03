package main

import (
	"flag"
	"log"
	"os"

	"github.com/sfbrigade/sfsbook/server"
	"github.com/sfbrigade/sfsbook/setup"
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

	global, err := setup.MakeGlobalState(pth)
	if err != nil {
		log.Fatalln("Can't make global state:", err)
	}

	// I don't think that I actually use the keys
	srv := server.MakeServer(":10443", global)
	server.Start(pth, srv)

}
