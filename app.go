// A straightforward but efficient way to block (or authorize) DNS domains.
// TODO: add ip filtering
// TODO: add auto-reload of blocklists

package main

import (
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	// this will sync our go routines
	var wg sync.WaitGroup

	// get command line arguments
	conf := readCliArgs()
	defer conf.logFileHAndle.Close()

	if conf.debug {
		log.Printf("%v", conf)
	}
	log.Printf("using resolver: %s", conf.resolver)

	// listen to this local address
	serverAddress := net.UDPAddr{
		Port: 53,
		IP:   net.ParseIP("127.0.0.1"),
	}

	// start udp UDPServer on previously defined address
	UDPServer, err := net.ListenUDP("udp", &serverAddress)
	if err != nil {
		log.Fatalf("error: <%v> when creating udp server on 127.0.0.1:53", err)
	}
	defer UDPServer.Close()
	log.Printf("listening to DNS requests")

	// launch goroutine to regularly update the blocklists
	go updateBlockLists(&conf)

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
		go handleDNSRequest(UDPServer, clientAddr, buf[:nbBytes], &conf)
	}

	wg.Wait()
}

// Update the blocklist regularly
func updateBlockLists(conf *Config) {
	for {
		// sleep before reading
		time.Sleep(time.Duration(conf.timeout) * time.Second)

		log.Printf("updating blocklists\n")
		conf.readBlocklists()
	}

}
