// Main structures for DNS headers, flags, etc
// Coming from https://datatracker.ietf.org/doc/html/rfc1035
package main

import (
	//"bytes"
	"encoding/binary"
	"strings"

	//"fmt"
	"io"
)

// Utility function to convert a bool to an uint16: no standard conversion offered by Go
func bool2int16(b bool) uint16 {
	if b {
		return uint16(1)
	}
	return uint16(0)
}

// From RFC1035: https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
type DNSPacketHeader struct {
	Id uint16 // A 16 bit identifier assigned by the program that
	//   generates any kind of query.  This identifier is copied
	//   the corresponding reply and can be used by the requester
	//   to match up replies to outstanding queries.
	Flags    uint16 // see below
	Qd_count uint16 // an unsigned 16 bit integer specifying the number of
	//    entries in the question section.
	An_count uint16 // an unsigned 16 bit integer specifying the number of
	// resource records in the answer section.
	Ns_count uint16 // an unsigned 16 bit integer specifying the number of name
	// server resource records in the authority records section.
	Ar_count uint16 // an unsigned 16 bit integer specifying the number of
	// resource records in the additional records section.
}

// Convert a buffer to a DNSPacketHeader struct, from BigEndian
func (header *DNSPacketHeader) fromNetworkBytes(rdr io.Reader) error {
	return binary.Read(rdr, binary.BigEndian, header)
}

// From RFC1035: : https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
type DNSPacketFlags struct {
	QR     byte // A one bit field that specifies whether this message is a query (0), or a response (1).
	OpCode byte // A four bit field that specifies kind of query in this
	//  message.  This value is set by the originator of a query
	//  and copied into the response.  The values are:
	// 0               a standard query (QUERY)
	// 1               an inverse query (IQUERY)
	// 2               a server status request (STATUS)
	// 3-15            reserved for future use
	AA bool // Authoritative Answer - this bit is valid in responses,
	//and specifies that the responding name server is an
	//authority for the domain name in question section.
	//Note that the contents of the answer section may have
	//multiple owner names because of aliases.  The AA bit
	//corresponds to the name which matches the query name, or
	//the first owner name in the answer section.
	TC bool //    TrunCation - specifies that this message was truncated
	//    due to length greater than that permitted on the
	//    transmission channel.
	RD bool // Recursion Desired - this bit may be set in a query and
	// is copied into the response.  If RD is set, it directs
	// the name server to pursue the query recursively.
	// Recursive query support is optional.
	RA bool // Recursion Available - this be is set or cleared in a
	//  response, and denotes whether recursive query support is
	//  available in the name server.
	Z     bool // Reserved for future use.  Must be zero in all queries and responses.
	AD    bool
	CD    bool
	RCODE byte // Response code - this 4 bit field is set as part of
	//responses.  The values have the following
	//interpretation:
	//0               No error condition
	//1               Format error - The name server was
	//                unable to interpret the query.
	//2               Server failure - The name server was
	//                unable to process this query due to a
	//                problem with the name server.
	//3               Name Error - Meaningful only for
	//                responses from an authoritative name
	//                server, this code signifies that the
	//                domain name referenced in the query does
	//                not exist.
	//4               Not Implemented - The name server does
	//                not support the requested kind of query.
	//5               Refused - The name server refuses to
	//                perform the specified operation for
	//                policy reasons.  For example, a name
	//                server may not wish to provide the
	//                information to the particular requester,
	//                or a name server may not wish to perform
	//                a particular operation (e.g., zone
	//                transfer) for particular data.
	//6-15            Reserved for future use.
}

// Convert a buffer to a DNSPacketFlags struct, from BigEndian
func (flags *DNSPacketFlags) fromNetworkBytes(value uint16) {
	// decode all flags according to structure
	//                               1  1  1  1  1  1
	// 0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// |                      ID                       |
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	// |QR|   Opcode  |AA|TC|RD|RA|Z |AD|CD|   RCODE   |
	// +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
	flags.QR = byte(value >> 15)
	flags.OpCode = byte(value >> 11 & 0b1111)
	flags.AA = (value>>10)&1 == 1
	flags.TC = (value>>9)&1 == 1
	flags.RD = (value>>8)&1 == 1
	flags.RA = (value>>7)&1 == 1
	flags.Z = (value >> 6 & 1) == 1
	flags.AD = (value >> 5 & 1) == 1
	flags.CD = (value >> 4 & 1) == 1
	flags.RCODE = byte(value & 0b1111)
}

// Write a DNSPacketFlags struct to a buffer
// func (flags *DNSPacketFlags) ToNetworkBytes(buffer []byte) error {
// 	// define a new buf
// 	buf := new(bytes.Buffer)

// 	// build a uint16 integer from flags
// 	value := uint16(flags.QR) << 15;
// 	value |= uint16(flags.OpCode) << 11;
// 	value |= bool2int16(flags.AA) << 10
// 	value |= bool2int16(flags.TC) << 9
// 	value |= bool2int16(flags.RD) << 8
// 	value |= bool2int16(flags.RA) << 7
// 	value |= bool2int16(flags.Z) << 6
// 	value |= bool2int16(flags.AD) << 5
// 	value |= bool2int16(flags.CD) << 4
// 	value |= uint16(flags.RCODE)

// 	fmt.Printf("value=%d\n", value)
// 	err := binary.Write(buf, binary.BigEndian, value)
// 	buf.Write(buffer)
// 	fmt.Printf("buffer in=%+v\n", buf.Bytes())
// 	return err
// }

// From RFC1035: https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.2
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

// Read the question: due to the way labels are defined, need to read the buffer and build the domain
// as a collection of labels
func (question *DNSQuestion) fromNetworkBytes(rdr io.Reader) error {
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

	// as "." was added (ex: www.google.com.) which normally is the way domains are expected, this is not convenient
	// for blacklisting. So delete last char
	question.Domain = strings.TrimSuffix(question.Domain, ".")

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

// Get the QType string from its numeric value
// RR type codes: https://www.iana.org/assignments/dns-parameters/dns-parameters.xhtml#dns-parameters-4
func qType(value uint16) string {
	switch value {
	case 1:
		return "A" // a host address	[RFC1035]
	case 2:
		return "NS" // an authoritative name server	[RFC1035]
	case 3:
		return "MD" // a mail destination (OBSOLETE - use MX)	[RFC1035]
	case 4:
		return "MF" // a mail forwarder (OBSOLETE - use MX)	[RFC1035]
	case 5:
		return "CNAME" // the canonical name for an alias	[RFC1035]
	case 6:
		return "SOA" // marks the start of a zone of authority	[RFC1035]
	case 7:
		return "MB" // a mailbox domain name (EXPERIMENTAL)	[RFC1035]
	case 8:
		return "MG" // a mail group member (EXPERIMENTAL)	[RFC1035]
	case 9:
		return "MR" // a mail rename domain name (EXPERIMENTAL)	[RFC1035]
	case 10:
		return "NULL" // a null RR (EXPERIMENTAL)	[RFC1035]
	case 11:
		return "WKS" // a well known service description	[RFC1035]
	case 12:
		return "PTR" // a domain name pointer	[RFC1035]
	case 13:
		return "HINFO" // host information	[RFC1035]
	case 14:
		return "MINFO" // mailbox or mail list information	[RFC1035]
	case 15:
		return "MX" // mail exchange	[RFC1035]
	case 16:
		return "TXT" // text strings	[RFC1035]
	case 17:
		return "RP" // for Responsible Person	[RFC1183]
	case 18:
		return "AFSDB" // for AFS Data Base location	[RFC1183][RFC5864]
	case 19:
		return "X25" // for X.25 PSDN address	[RFC1183]
	case 20:
		return "ISDN" // for ISDN address	[RFC1183]
	case 21:
		return "RT" // for Route Through	[RFC1183]
	case 22:
		return "NSAP" // for NSAP address, NSAP style A record	[RFC1706]
	case 23:
		return "NSAPPTR" // for domain name pointer, NSAP style	[RFC1706]
	case 24:
		return "SIG" // for security signature	[RFC2536][RFC2931][RFC3110][RFC4034]
	case 25:
		return "KEY" // for security key	[RFC2536][RFC2539][RFC3110][RFC4034]
	case 26:
		return "PX" // X.400 mail mapping information	[RFC2163]
	case 27:
		return "GPOS" // Geographical Position	[RFC1712]
	case 28:
		return "AAAA" // IP6 Address	[RFC3596]
	case 29:
		return "LOC" // Location Information	[RFC1876]
	case 30:
		return "NXT" // Next Domain (OBSOLETE)	[RFC2535][RFC3755]
	case 31:
		return "EID" // Endpoint Identifier	[Michael_Patton][http://ana-3.lcs.mit.edu/~jnc/nimrod/dns.txt]		1995-06
	case 32:
		return "NIMLOC" // Nimrod Locator	[1][Michael_Patton][http://ana-3.lcs.mit.edu/~jnc/nimrod/dns.txt]		1995-06
	case 33:
		return "SRV" // Server Selection	[1][RFC2782]
	case 34:
		return "ATMA" // ATM Address	[ ATM Forum Technical Committee, "ATM Name System, V2.0", Doc ID: AF-DANS-0152.000, July 2000. Available from and held in escrow by IANA.]
	case 35:
		return "NAPTR" // Naming Authority Pointer	[RFC3403]
	case 36:
		return "KX" // Key Exchanger	[RFC2230]
	case 37:
		return "CERT" // CERT	[RFC4398]
	case 38:
		return "A6" // A6 (OBSOLETE - use AAAA)	[RFC2874][RFC3226][RFC6563]
	case 39:
		return "DNAME" // DNAME	[RFC6672]
	case 40:
		return "SINK" // SINK	[Donald_E_Eastlake][draft-eastlake-kitchen-sink]		1997-11
	case 41:
		return "OPT" // OPT	[RFC3225][RFC6891]
	case 42:
		return "APL" // APL	[RFC3123]
	case 43:
		return "DS" // Delegation Signer	[RFC4034]
	case 44:
		return "SSHFP" // SSH Key Fingerprint	[RFC4255]
	case 45:
		return "IPSECKEY" // IPSECKEY	[RFC4025]
	case 46:
		return "RRSIG" // RRSIG	[RFC4034]
	case 47:
		return "NSEC" // NSEC	[RFC4034][RFC9077]
	case 48:
		return "DNSKEY" // DNSKEY	[RFC4034]
	case 49:
		return "DHCID" // DHCID	[RFC4701]
	case 50:
		return "NSEC3" // NSEC3	[RFC5155][RFC9077]
	case 51:
		return "NSEC3PARAM" // NSEC3PARAM	[RFC5155]
	case 52:
		return "TLSA" // TLSA	[RFC6698]
	case 53:
		return "SMIMEA" // S/MIME cert association	[RFC8162]	SMIMEA/smimea-completed-template	2015-12-01
	case 54:
		return "Unassigned" //
	case 55:
		return "HIP" // Host Identity Protocol	[RFC8005]
	case 56:
		return "NINFO" // NINFO	[Jim_Reid]	NINFO/ninfo-completed-template	2008-01-21
	case 57:
		return "RKEY" // RKEY	[Jim_Reid]	RKEY/rkey-completed-template	2008-01-21
	case 58:
		return "TALINK" // Trust Anchor LINK	[Wouter_Wijngaards]	TALINK/talink-completed-template	2010-02-17
	case 59:
		return "CDS" // Child DS	[RFC7344]	CDS/cds-completed-template	2011-06-06
	case 60:
		return "CDNSKEY" // DNSKEY(s) the Child wants reflected in DS	[RFC7344]		2014-06-16
	case 61:
		return "OPENPGPKEY" // OpenPGP Key	[RFC7929]	OPENPGPKEY/openpgpkey-completed-template	2014-08-12
	case 62:
		return "CSYNC" // Child-To-Parent Synchronization	[RFC7477]		2015-01-27
	case 63:
		return "ZONEMD" // Message Digest Over Zone Data	[RFC8976]	ZONEMD/zonemd-completed-template	2018-12-12
	case 64:
		return "SVCB" // Service Binding	[draft-ietf-dnsop-svcb-https-00]	SVCB/svcb-completed-template	2020-06-30
	case 65:
		return "HTTPS" // HTTPS Binding	[draft-ietf-dnsop-svcb-https-00]	HTTPS/https-completed-template	2020-06-30
	// Unassigned	66-98
	case 99:
		return "SPF" // [RFC7208]
	case 100:
		return "UINFO" // [IANA-Reserved]
	case 101:
		return "UID" // [IANA-Reserved]
	case 102:
		return "GID" // [IANA-Reserved]
	case 103:
		return "UNSPEC" // [IANA-Reserved]
	case 104:
		return "NID" // [RFC6742]	ILNP/nid-completed-template
	case 105:
		return "L32" // [RFC6742]	ILNP/l32-completed-template
	case 106:
		return "L64" // [RFC6742]	ILNP/l64-completed-template
	case 107:
		return "LP" // [RFC6742]	ILNP/lp-completed-template
	case 108:
		return "EUI48" // an EUI-48 address	[RFC7043]	EUI48/eui48-completed-template	2013-03-27
	case 109:
		return "EUI64" // an EUI-64 address	[RFC7043]	EUI64/eui64-completed-template	2013-03-27
	// Unassigned	110-248
	case 249:
		return "TKEY" // Transaction Key	[RFC2930]
	case 250:
		return "TSIG" // Transaction Signature	[RFC8945]
	case 251:
		return "IXFR" // incremental transfer	[RFC1995]
	case 252:
		return "AXFR" // transfer of an entire zone	[RFC1035][RFC5936]
	case 253:
		return "MAILB" // mailbox-related RRs (MB, MG or MR)	[RFC1035]
	case 254:
		return "MAILA" // mail agent RRs (OBSOLETE - see MX)	[RFC1035]
	case 255:
		return "ANY" // A request for some or all records the server has available	[RFC1035][RFC6895][RFC8482]
	case 256:
		return "URI" // URI	[RFC7553]	URI/uri-completed-template	2011-02-22
	case 257:
		return "CAA" // Certification Authority Restriction	[RFC8659]	CAA/caa-completed-template	2011-04-07
	case 258:
		return "AVC" // Application Visibility and Control	[Wolfgang_Riedel]	AVC/avc-completed-template	2016-02-26
	case 259:
		return "DOA" // Digital Object Architecture	[draft-durand-doa-over-dns]	DOA/doa-completed-template	2017-08-30
	case 260:
		return "AMTRELAY" // Automatic Multicast Tunneling Relay	[RFC8777]	AMTRELAY/amtrelay-completed-template	2019-02-06
	// Unassigned	261-32767
	case 32768:
		return "TA" // DNSSEC Trust Authorities	[Sam_Weiler][http://cameo.library.cmu.edu/][ Deploying DNSSEC Without a Signed Root. Technical Report 1999-19, Information Networking Institute, Carnegie Mellon University, April 2004.]		2005-12-13
	case 32769:
		return "DLV" // DNSSEC Lookaside Validation (OBSOLETE)	[RFC8749][RFC4431]
	}
	return ""
}
