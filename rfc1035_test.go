package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNSPacketHeader(t *testing.T) {
	assert := assert.New(t)

	buffer := []byte{0x30, 0x5c, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}
	header := new(DNSPacketHeader)
	err := header.FromNetworkBytes(bytes.NewReader(buffer))

	assert.NotNil(t, err)

	assert.Equal(header.Id, uint16(0x305c))
	assert.Equal(header.Flags, uint16(0x0100))
	assert.Equal(header.Qd_count, uint16(1))
	assert.Equal(header.An_count, uint16(0))
	assert.Equal(header.Ns_count, uint16(0))
	assert.Equal(header.Ar_count, uint16(1))
}

func TestDNSFlags(t *testing.T) {
	assert := assert.New(t)

	flags := new(DNSPacketFlags)
	flags.FromNetworkBytes(uint16(0b1000_1111_1111_0000))

	assert.Equal(flags.QR, uint8(1))
	assert.Equal(flags.OpCode, uint8(1))
	assert.True(flags.AA)
	assert.True(flags.TC)
	assert.True(flags.RD)
	assert.True(flags.RA)
	assert.True(flags.Z)
	assert.True(flags.AD)
	assert.True(flags.CD)
	assert.Equal(flags.RCODE, uint8(0))

}

func TestDNSQuestion(t *testing.T) {
	assert := assert.New(t)

	buffer := []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x000, 0x01}
	question := new(DNSQuestion)
	err := question.FromNetworkBytes(bytes.NewReader(buffer))

	assert.NotNil(t, err)

	assert.Equal(question.Domain, "www.google.com.")
	assert.Equal(question.QType, uint16(1))
	assert.Equal(question.QClass, uint16(1))
}
