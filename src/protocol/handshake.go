package protocol

import (
	"io"
)

const (
	socks5Version = 0x05
	rsv           = 0x00
)

// MaxAddrLen is the maximum size of SOCKS address in bytes.
const MaxAddrLen = 1 + 1 + 255 + 2

type MethodNegotiationReq []byte

func ReadMethodNegotiationReq(r io.Reader) (MethodNegotiationReq, error) {
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	var (
		req MethodNegotiationReq
		buf = make([]byte, MaxAddrLen)
	)
	// read VER, NMETHODS
	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		return req, err
	}
	nmethods := buf[1]
	// read METHODS
	if _, err := io.ReadFull(r, buf[:nmethods]); err != nil {
		return req, err
	}
	return buf[:nmethods], nil
}

func WriteMethodNegotiationReply(method byte, w io.Writer) error {
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	_, err := w.Write([]byte{socks5Version, method})
	return err
}

var (
	NoAcceptableMethods = []byte{socks5Version, 0xff}
)

type CommandMessage struct {
	command  byte
	addrType byte
	addr     string
	port     int
}

// Socks Command Request
// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
// Socks Command Reply
// +----+-----+-------+------+----------+----------+
// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+
