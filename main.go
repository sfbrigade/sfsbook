package main

import (
	"flag"
	"log"
	"os"

	"github.com/sfbrigade/sfsbook/server"
	"github.com/sfbrigade/sfsbook/setup"
)

var port = flag.String("port", ":10443", "Set a port to listen on")

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

	srv := server.MakeServer(*port, handlerfactory, certfactory)

	// TODO(sa): enable this once we've a cert
	//if err := srv.ListenAndServeTLS(
	//  certfactory.GetCertfFileName(),
	//  certfactory.GetKeyFileName()); err != nil {
	//  log.Fatal("serving went wrong", err)
	//}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("serving went wrong", err)
	}
}
