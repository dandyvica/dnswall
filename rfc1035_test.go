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
	err := header.fromNetworkBytes(bytes.NewReader(buffer))

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

	// FromNetworkBytes
	flags := new(DNSPacketFlags)
	flags.fromNetworkBytes(uint16(0b1000_1111_1111_0000))

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

	// ToNetworkBytes
	// flags.QR = 1
	// flags.OpCode = 1
	// flags.AA = true
	// flags.TC = true
	// flags.RD = true
	// flags.RA = true
	// flags.Z = true
	// flags.AD = true
	// flags.CD = true
	// flags.RCODE = 3
	// buffer := make([]byte, 2)
	// err := flags.ToNetworkBytes(buffer)
	// assert.Nil(err)
	// fmt.Printf("buffer=%+v\n", buffer)
	// assert.Equal(buffer, []byte{0b1111_1111, 0b1111_0011})

}

func TestDNSQuestion(t *testing.T) {
	assert := assert.New(t)

	buffer := []byte{0x03, 0x77, 0x77, 0x77, 0x06, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01}
	question := new(DNSQuestion)
	err := question.fromNetworkBytes(bytes.NewReader(buffer))

	assert.NotNil(t, err)

	assert.Equal(question.Domain, "www.google.com")
	assert.Equal(question.QType, uint16(1))
	assert.Equal(question.QClass, uint16(1))
}
