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
	log.Println("sfsbook starting")

	pth, err := os.Getwd()
	if err != nil {
		log.Fatalln("Wow! No CWD. Giving up.", err)
	}

	certfactory, err := setup.MakeCertFactory(pth)
	if err != nil {
		log.Fatalln("Can't make CertFactory:", err)
	}

	handlerfactory, err := server.MakeHandlerFactory(pth)
	if err != nil {
		log.Fatalln("Can't make HandlerFactory:", err)
	}

	// I don't think that I actually use the certificates properly.
	srv := server.MakeServer(":10443", handlerfactory)

	if err := srv.ListenAndServeTLS(
		certfactory.GetCertfFileName(),
		certfactory.GetKeyFileName()); err != nil {
		log.Fatal("serving went wrong", err)
	}

}
