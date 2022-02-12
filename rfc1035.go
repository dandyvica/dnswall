// Main structure for DNS headers, flags, etc
package main

import (
	"encoding/binary"
	"io"
)

type DNSPacketHeader struct {
	Id uint16 // A 16 bit identifier assigned by the program that
	//   generates any kind of query.  This identifier is copied
	//   the corresponding reply and can be used by the requester
	//   to match up replies to outstanding queries.
	Flags    uint16
	Qd_count uint16 // an unsigned 16 bit integer specifying the number of
	//    entries in the question section.
	An_count uint16 // an unsigned 16 bit integer specifying the number of
	// resource records in the answer section.
	Ns_count uint16 // an unsigned 16 bit integer specifying the number of name
	// server resource records in the authority records section.
	Ar_count uint16 // an unsigned 16 bit integer specifying the number of
	// resource records in the additional records section.
}

// Convert a buffer to a DNSPacketHeader struct
func (header *DNSPacketHeader) FromNetworkBytes(rdr io.Reader) error {
	return binary.Read(rdr, binary.BigEndian, header)
}

// Question
type DNSQuestion struct {
	Domain string // a domain name represented as a sequence of labels, where
	//             each label consists of a length octet followed by that
	//             number of octets.  The domain name terminates with the
	//             zero length octet for the null label of the root.  Note
	//             that this field may be an odd number of octets; no
	//             padding is used.
	QType uint16 // a two octet code which specifies the type of the query.
	//             The values for this field include all codes valid for a
	//             TYPE field, together with some more general codes which
	//             can match more than one type of RR.
	QClass uint16 // a two octet code that specifies the class of the query.
	// For example, the QCLASS field is IN for the Internet.
}

// Read the question
func (question *DNSQuestion) FromNetworkBytes(rdr io.Reader) error {
	for {
		var buffer [1]byte

		_, err := rdr.Read(buffer[:])
		if err != nil {
			return err
		}

		// sentinel found
		if buffer[0] == 0 {
			break
		}

		// buffer[0] contains the length of label
		size := buffer[0]

		// so read 'size' bytes for the reader
		label := make([]byte, size)

		_, err = rdr.Read(label)
		if err != nil {
			return err
		}

		question.Domain += string(label) + "."

	}

	// now read QType and QClass: read buffer and convert to uint16
	var buffer [4]byte

	_, err := rdr.Read(buffer[:])
	if err != nil {
		return err
	}

	question.QType = binary.BigEndian.Uint16(buffer[0:2])
	question.QClass = binary.BigEndian.Uint16(buffer[2:4])

	return nil
}
