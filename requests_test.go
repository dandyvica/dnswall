package main

import (
	//"fmt"
	"bytes"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDomainQuestion(t *testing.T) {
	assert := assert.New(t)

	// a query taken from Wireshark
	buffer := []byte{0x30, 0x5c, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x000, 0x01}
	options := new(CliOptions)
	domain, _ := GetDomainQuestion(buffer, options)

	assert.Equal(domain.Domain, "www.google.com")
	assert.Equal(domain.QType, uint16(1))
	assert.Equal(domain.QClass, uint16(1))
}

func TestQueryResolver(t *testing.T) {
	assert := assert.New(t)

	// a query taken from Wireshark using dig: $> dig @8.8.8.8 A www.google.com
	query := []byte{0xbd, 0x73, 0x01, 0x20, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}

	options := new(CliOptions)
	options.resolverAddress = "8.8.8.8:53"
	addr := net.UDPAddr{
		IP: net.ParseIP("0.0.0.0"),
	}

	buffer, _, err := QueryResolver(query, options, &addr)
	assert.Nil(err)

	// define a new reader
	rdr := bytes.NewReader(buffer)

	// read DNS header
	header := new(DNSPacketHeader)
	err = header.FromNetworkBytes(rdr)
	assert.Nil(err)
	flags := new(DNSPacketFlags)
	flags.FromNetworkBytes(header.Flags)
	fmt.Printf("%+v\n", flags)
	assert.Equal(flags.QR, byte(1))
	assert.Equal(flags.RCODE, byte(0))

}
