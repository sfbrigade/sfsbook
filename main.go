package main

import (
	"flag"
	"log"
	"os"

	"github.com/sfbrigade/sfsbook/server"
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

	handlerfactory, err := server.MakeHandlerFactory(pth)
	if err != nil {
		log.Fatalln("Can't make HandlerFactory:", err)
	}

	// I don't think that I actually use the certificates properly.
	srv := server.MakeServer(":10443", handlerfactory)
	server.Start(pth, srv)

}
