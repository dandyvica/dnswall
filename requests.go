package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
)

const (
	DEFAULT_BUFFER_SIZE = 2048
)

// This functions is call by the UDP server to server requests
func handleDNSRequest(packet net.PacketConn, requesterAddress net.Addr, buffer []byte, options *CliOptions) {
	// get DNS question from initial request
	question, err := GetDomainQuestion(buffer, options)
	if err != nil {
		return
	}
	log.Printf("received request <%s> for domain: <%s> for requester: <%v>", GetQType(question.QType), question.Domain, requesterAddress)

	// send question to resolver and wait for its answer
	answerBuffer, nbReadBytes, err := GetAnswerFromResolver(buffer, options, requesterAddress)
	if err != nil {
		return
	}

	// send back answer coming from resolver to requester
	nbWrittenBytes, err := packet.WriteTo(answerBuffer[:nbReadBytes], requesterAddress)
	if err != nil {
		log.Printf("error: <%v> when writing back to DNS requester", err)
		return
	}
	if options.debug {
		log.Printf("%v bytes written back to requester on behalf of <%s>", nbWrittenBytes, requesterAddress)
	}
}

// Get domain name from the request coming from the client
func GetDomainQuestion(buffer []byte, options *CliOptions) (*DNSQuestion, error) {
	// define a new reader
	rdr := bytes.NewReader(buffer)

	// read DNS header
	header := new(DNSPacketHeader)
	err := header.FromNetworkBytes(rdr)
	if err != nil {
		log.Printf("error: <%v> when converting buffer to DNS packet header", err)
		return nil, err
	}

	if options.debug {
		log.Printf("header=%+v", header)
		flags := new(DNSPacketFlags)
		flags.FromNetworkBytes(header.Flags)
		log.Printf("flags=%+v", flags)
	}

	// now read questions, but normally there's only one question
	if header.Qd_count != 1 {
		log.Printf("several questions (%d) in a single query, we don't expect this!", header.Qd_count)
	}

	// retrieve question
	question := new(DNSQuestion)
	err = question.FromNetworkBytes(rdr)
	if err != nil {
		log.Printf("error: <%v> when converting buffer to DNS question", err)
		return nil, err
	}
	if options.debug {
		log.Printf("question=%+v", question)
	}

	return question, nil
}

// Send request to resolver and wait for request
func GetAnswerFromResolver(buffer []byte, options *CliOptions, requesterAddress net.Addr) ([]byte, int, error) {
	// open connection to DNS resolver
	conn, err := net.Dial("udp", options.resolverAddress)
	if err != nil {
		log.Printf("error: <%v> when connecting to DNS resolver: <%s>", options.resolverAddress, err)
		return nil, 0, err
	}

	// forward DNS request coming from client to the resolver
	nbWrittenBytes, err := conn.Write(buffer)
	if err != nil {
		log.Printf("error: <%v> when writing to DNS resolver", err)
		return nil, 0, err
	}
	if options.debug {
		log.Printf("%v bytes sent to resolver on behalf of <%s>", nbWrittenBytes, requesterAddress)
	}

	// wait for answer from resolver
	answerBuffer := make([]byte, DEFAULT_BUFFER_SIZE)
	nbReadBytes, err := bufio.NewReader(conn).Read(answerBuffer)
	if err != nil {
		log.Printf("error: <%v> when reading from DNS resolver", err)
		return nil, 0, err
	}
	if options.debug {
		log.Printf("%v bytes read from resolver on behalf of <%s>", nbReadBytes, requesterAddress)
	}

	return answerBuffer, nbReadBytes, nil
}
