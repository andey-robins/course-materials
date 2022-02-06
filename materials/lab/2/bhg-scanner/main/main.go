package main

import (
	"flag"
	"log"

	"bhg-scanner/scanner"
)

func main(){
	address := flag.String("address", "", "Address is the url to scan")
	fname := flag.String("output", "", "The filename to store the output csv in")

	flag.Parse()

	if address == nil || *address == "" {
		log.Println("command line arg -address is required")
		return
	}

	scanner.PortScanner(*address, *fname)
}