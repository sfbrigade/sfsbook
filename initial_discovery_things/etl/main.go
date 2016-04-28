package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.Println("foo")

	tr := &http.Transport{ }
	client := &http.Client{Transport: tr}

	// This is https://data.sfgov.org/Public-Safety/SFPD-Incidents-from-1-January-2003/tmnf-yvry
	// TODO(rjk): Extract a window.
	resp, err := client.Get("https://data.sfgov.org/resource/tmnf-yvry.json?$limit=100")
	if err != nil {
		log.Panic("Get had a sad! ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("ReadAll had a sad! ", err)
	}

	// TODO(rjk): Parse the JSON
	// TODO(rjk): Persist the data
	// TODO(rjk): Queries.

	log.Println("data: \n", string(body))
}
