package main

import (
	"bytes"
	//"fmt"
	"log"
	"net"
	"sync"
	"bufio"
)

func main() {
	// this will sync our go routines
	var wg sync.WaitGroup

	// start udp server on local address
	server, err := net.ListenPacket("udp4", "127.0.0.1:53")
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		// read data from client
		buf := make([]byte, 1024)
		n, client_addr, err := server.ReadFrom(buf)
		log.Printf("%d bytes received from address: %v\n", n, client_addr)
		if err != nil {
			continue
		}
		go serve(server, client_addr, buf[:n])
	}

	wg.Wait()
}

func serve(pc net.PacketConn, addr net.Addr, buffer []byte) {
	// define a new reader
	rdr := bytes.NewReader(buffer)

	// read DNS header
	header := new(DNSPacketHeader)
	_ = header.FromNetworkBytes(rdr)
	log.Printf("header=%+v", header)

	// now read questions
	for i := uint16(0); i < header.Qd_count; i++ {
		question := new(DNSQuestion)
		question.FromNetworkBytes(rdr)
		log.Printf("question=%+v", question)
	}

	// send back question to DNS resolver
    conn, err := net.Dial("udp", "1.1.1.1:53")
    if err != nil {
        log.Printf("error in connection to DNS resolver: %v", err)
        return
    }	
	conn.Write(buffer)

	// wait answer from resolver
	p :=  make([]byte, 2048)
	_, err = bufio.NewReader(conn).Read(p)

	// send back to client
	pc.WriteTo(p, addr)
}
