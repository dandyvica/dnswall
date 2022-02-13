package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
)

func handleDNSRequest(packet net.PacketConn, requesterAddress net.Addr, buffer []byte, options *CliOptions) {
	// define a new reader
	rdr := bytes.NewReader(buffer)

	// read DNS header
	header := new(DNSPacketHeader)
	err := header.FromNetworkBytes(rdr)
	if err != nil {
		log.Printf("error: <%v> when converting buffer to DNS packet header", err)
		return
	}

	if options.debug {
		log.Printf("header=%+v", header)
		flags := new(DNSPacketFlags)
		flags.FromNetworkBytes(header.Flags)
		log.Printf("flags=%+v", flags)
	}

	// now read questions
	for i := uint16(0); i < header.Qd_count; i++ {
		question := new(DNSQuestion)
		err = question.FromNetworkBytes(rdr)
		if err != nil {
			log.Printf("error: <%v> when converting buffer to DNS question", err)
			return
		}
		if options.debug {
			log.Printf("question=%+v", question)
		}
		log.Printf("received request <%s> for domain: <%s> for requester: <%v>", GetQType(question.QType), question.Domain, requesterAddress)
	}

	// open connection to DNS resolver
	conn, err := net.Dial("udp", options.resolverAddress)
	if err != nil {
		log.Printf("error: <%v> when connecting to DNS resolver: <%s>", options.resolverAddress, err)
		return
	}

	// forward DNS request coming from client to the resolver
	nbWrittenBytes, err := conn.Write(buffer)
	if err != nil {
		log.Printf("error: <%v> when writing to DNS resolver", err)
		return
	}
	if options.debug {
		log.Printf("%v bytes sent to resolver on behalf of <%s>", nbWrittenBytes, requesterAddress)
	}

	// wait for answer from resolver
	answerBuffer := make([]byte, 2048)
	nbReadBytes, err := bufio.NewReader(conn).Read(answerBuffer)
	if err != nil {
		log.Printf("error: <%v> when reading from DNS resolver", err)
		return
	}
	if options.debug {
		log.Printf("%v bytes read from resolver on behalf of <%s>", nbReadBytes, requesterAddress)
	}

	// send back answer coming from resolver to requester
	nbWrittenBytes, err = packet.WriteTo(answerBuffer[:nbReadBytes], requesterAddress)
	if err != nil {
		log.Printf("error: <%v> when writing back to DNS requester", err)
		return
	}
	if options.debug {
		log.Printf("%v bytes written back to requester on behalf of <%s>", nbWrittenBytes, requesterAddress)
	}
}
