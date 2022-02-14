package main

import (
	"log"
	"net"
	"sync"
)

func main() {
	// this will sync our go routines
	var wg sync.WaitGroup

	// get command line arguments
	options := CliArgs()
	defer options.logFileHAndle.Close()

	if options.debug {
		log.Printf("%v", options)
	}

	// load yaml file
	var config Config
	config.getConfig(options.configFile)

	log.Printf("using resolver: %s", options.resolver)

	// start udp UDPServer on local address
	UDPServer, err := net.ListenPacket("udp4", "127.0.0.1:53")
	if err != nil {
		log.Fatalf("error: <%v> when creating udp server on 127.0.0.1:53", err)
	}
	defer UDPServer.Close()
	log.Printf("listening to DNS requests")

	// handle DNS requests from clients
	for {
		// read data from client
		buf := make([]byte, 1024)
		nbBytes, clientAddr, err := UDPServer.ReadFrom(buf)
		if err != nil {
			log.Printf("error: <%v> reading bytes from address: <%v>\n", nbBytes, clientAddr)
			continue
		}

		// serve request
		log.Printf("%d bytes received from address: %v\n", nbBytes, clientAddr)
		go handleDNSRequest(UDPServer, clientAddr, buf[:nbBytes], &options)
	}

	wg.Wait()
}
