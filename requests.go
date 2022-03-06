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
func handleDNSRequest(conn *net.UDPConn, requesterAddress net.Addr, buffer []byte, conf *Config) {
	//defer conf.mu.Unlock()

	// get DNS question from initial request
	question, err := getDomainQuestion(buffer, conf)
	if err != nil {
		return
	}
	log.Printf("received request <%s> for domain: <%s> for requester: <%v>", qType(question.QType), question.Domain, requesterAddress)

	// if domain name is in the whitelist => pass
	// if not, if in blacklist => reject
	// otherwise => pass
	conf.mu.Lock()
	if conf.filters.isFiltered(question.Domain) && !conf.dontFilter {
		err = rejectDomain(conn, buffer, requesterAddress)
		if err != nil {
			return
		}
		log.Printf("domain <%s> is blacklisted", question.Domain)
		return
	}
	conf.mu.Unlock()

	// send question to resolver and wait for its answer
	answerBuffer, nbReadBytes, err := queryResolver(buffer, conf, requesterAddress)
	if err != nil {
		return
	}

	// send back answer coming from resolver to requester
	nbWrittenBytes, err := conn.WriteTo(answerBuffer[:nbReadBytes], requesterAddress)
	if err != nil {
		log.Printf("error: <%v> when writing back to DNS requester", err)
		return
	}
	if conf.debug {
		log.Printf("%v bytes written back to requester on behalf of <%s>", nbWrittenBytes, requesterAddress)
	}
}

// Get domain name from the request coming from the client
func getDomainQuestion(buffer []byte, conf *Config) (*DNSQuestion, error) {
	// define a new reader
	rdr := bytes.NewReader(buffer)

	// read DNS header
	header := new(DNSPacketHeader)
	err := header.fromNetworkBytes(rdr)
	if err != nil {
		log.Printf("error: <%v> when converting buffer to DNS packet header", err)
		return nil, err
	}

	if conf.debug {
		log.Printf("header=%+v", header)
		flags := new(DNSPacketFlags)
		flags.fromNetworkBytes(header.Flags)
		log.Printf("flags=%+v", flags)
	}

	// now read questions, but normally there's only one question
	if header.Qd_count != 1 {
		log.Printf("several questions (%d) in a single query, we don't expect this!", header.Qd_count)
	}

	// retrieve question
	question := new(DNSQuestion)
	err = question.fromNetworkBytes(rdr)
	if err != nil {
		log.Printf("error: <%v> when converting buffer to DNS question", err)
		return nil, err
	}
	if conf.debug {
		log.Printf("question=%+v", question)
	}

	return question, nil
}

// Send request to resolver and wait for request
func queryResolver(buffer []byte, conf *Config, requesterAddress net.Addr) ([]byte, int, error) {
	// open connection to DNS resolver
	conn, err := net.Dial("udp", conf.resolverAddress)
	if err != nil {
		log.Printf("error: <%v> when connecting to DNS resolver: <%s>", conf.resolverAddress, err)
		return nil, 0, err
	}

	// forward DNS request coming from client to the resolver
	nbWrittenBytes, err := conn.Write(buffer)
	if err != nil {
		log.Printf("error: <%v> when writing to DNS resolver", err)
		return nil, 0, err
	}
	if conf.debug {
		log.Printf("%v bytes sent to resolver on behalf of <%s>", nbWrittenBytes, requesterAddress)
	}

	// wait for answer from resolver
	answerBuffer := make([]byte, DEFAULT_BUFFER_SIZE)
	nbReadBytes, err := bufio.NewReader(conn).Read(answerBuffer)
	if err != nil {
		log.Printf("error: <%v> when reading from DNS resolver", err)
		return nil, 0, err
	}
	if conf.debug {
		log.Printf("%v bytes read from resolver on behalf of <%s>", nbReadBytes, requesterAddress)
	}

	return answerBuffer, nbReadBytes, nil
}

// Respond with a NXDOMAIN to the requester to mean domain is not existing
func rejectDomain(conn *net.UDPConn, buffer []byte, requesterAddress net.Addr) error {
	// set flags: QR =1 (it's a response) and RCODE to 3 = NXDOMAIN
	// flags is at index 2 and 3 in the buffer
	buffer[2] |= 0b1000_0000 // set QR = Response = 1
	buffer[3] |= 0b1000_0011 // set RCODE = NXDOMAIN

	// send back answer coming from resolver to requester
	_, err := conn.WriteTo(buffer, requesterAddress)
	if err != nil {
		log.Printf("error: <%v> when writing NXDOMAIN DNS requester", err)
		return err
	}
	return nil
}
