package protocol

import (
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	Socks5Version = 0x05
	Rsv           = 0x00

	AddrTypeIPv4       = 1
	AddrTypeDomainName = 3
	AddrTypeIPv6       = 4

	CommandNegoSucceed = 0x00

	portLen = 2
)

type NetworkError struct {
	err error
}

func (err NetworkError) Error() string {
	return err.err.Error()
}

// MaxAddrLen is the maximum size of SOCKS address in bytes.
const MaxAddrLen = 1 + 1 + 255 + 2

// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+

func ReadMethodNegotiationReq(conn net.Conn, buf []byte) ([]byte, error) {
	// read VER, NMETHODS
	if _, err := io.ReadFull(conn, buf[:2]); err != nil {
		return nil, NetworkError{err: err}
	}
	version := buf[0]
	if version != Socks5Version {
		return nil, unknownProtocol
	}

	nmethods := buf[1]
	// read METHODS
	if _, err := io.ReadFull(conn, buf[:nmethods]); err != nil {
		return nil, NetworkError{err: err}
	}
	return buf[:nmethods], nil
}

// +----+--------+
// |VER | METHOD |
// +----+--------+
// | 1  |   1    |
// +----+--------+

func WriteMethodNegotiationReply(method byte, conn net.Conn) error {
	_, err := conn.Write([]byte{Socks5Version, method})
	if err != nil {
		return NetworkError{err: err}
	}
	return nil
}

type Addr []byte

func (addr Addr) String() string {
	var host, port string

	switch addr[0] { // address type
	case AddrTypeDomainName:
		host = string(addr[2 : 2+int(addr[1])])
		port = strconv.Itoa((int(addr[2+int(addr[1])]) << 8) | int(addr[2+int(addr[1])+1]))
	case AddrTypeIPv4:
		host = net.IP(addr[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(addr[1+net.IPv4len]) << 8) | int(addr[1+net.IPv4len+1]))
	case AddrTypeIPv6:
		host = net.IP(addr[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(addr[1+net.IPv6len]) << 8) | int(addr[1+net.IPv6len+1]))
	}

	return net.JoinHostPort(host, port)
}

func ParseAddr(s string) Addr {
	var addr Addr
	host, port, err := net.SplitHostPort(s)
	if err != nil {
		return nil
	}
	if ip := net.ParseIP(host); ip != nil {
		if ip4 := ip.To4(); ip4 != nil {
			addr = make([]byte, 1+net.IPv4len+2)
			addr[0] = AddrTypeIPv4
			copy(addr[1:], ip4)
		} else {
			addr = make([]byte, 1+net.IPv6len+2)
			addr[0] = AddrTypeIPv6
			copy(addr[1:], ip)
		}
	} else {
		if len(host) > 255 {
			return nil
		}
		addr = make([]byte, 1+1+len(host)+2)
		addr[0] = AddrTypeDomainName
		addr[1] = byte(len(host))
		copy(addr[2:], host)
	}

	portnum, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return nil
	}

	addr[len(addr)-2], addr[len(addr)-1] = byte(portnum>>8), byte(portnum)

	return addr
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

func ReadCommandNegotiationReq(conn net.Conn, buf []byte) (byte, Addr, error) {
	// read VER CMD RSV
	if _, err := io.ReadFull(conn, buf[:3]); err != nil {
		return 0, nil, NetworkError{err: err}
	}
	command := buf[1]

	// read ATYP DST.ADDR DST.PORT
	addr, err := readAddr(conn, buf)
	if err != nil {
		return 0, nil, err
	}
	return command, addr, nil
}

func WriteCommandNegotiationReply(conn net.Conn, buf []byte, source string) error {
	buf[0] = Socks5Version
	buf[1] = CommandNegoSucceed
	buf[2] = Rsv

	addr := ParseAddr(source)
	if addr == nil {
		return GeneralSocksServerFailure
	}
	_, err := conn.Write(buf[:3])
	if err != nil {
		return NetworkError{err: err}
	}
	_, err = conn.Write(addr)
	if err != nil {
		return NetworkError{err: err}
	}
	return nil
}

func readAddr(r io.Reader, buf []byte) ([]byte, error) {
	// read ATYP
	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return nil, NetworkError{err: err}
	}

	switch buf[0] {
	case AddrTypeDomainName:
		// read one byte for domain length
		if _, err := io.ReadFull(r, buf[1:2]); err != nil {
			return nil, NetworkError{err: err}
		}
		domainLen := int(buf[1])
		addrBuf := buf[2 : 2+domainLen+portLen]
		if _, err := io.ReadFull(r, addrBuf); err != nil {
			return nil, NetworkError{err: err}
		}
		return buf[:1+1+domainLen+portLen], nil
	case AddrTypeIPv4:
		addrBuf := buf[1 : 1+net.IPv4len+portLen]
		if _, err := io.ReadFull(r, addrBuf); err != nil {
			return nil, NetworkError{err: err}
		}
		return buf[:1+net.IPv4len+portLen], nil
	case AddrTypeIPv6:
		addrBuf := buf[1 : 1+net.IPv6len+portLen]
		if _, err := io.ReadFull(r, addrBuf); err != nil {
			return nil, NetworkError{err: err}
		}
		return buf[:1+net.IPv6len+portLen], nil
	default:
		return nil, AddressTypeNotSupported
	}
}

type SocksError struct {
	msg   string
	cache []byte
}

func (err SocksError) Error() string {
	return fmt.Sprintf("Socks protocol err: %s", err.msg)
}

func (err SocksError) SendErrorReply(conn net.Conn) {
	_, _ = conn.Write(err.cache)
}

var (
	unknownProtocol  = fmt.Errorf("unknown protocol")
	NoAcceptedMethod = SocksError{
		msg:   "NO ACCEPTABLE METHODS",
		cache: []byte{Socks5Version, 0xff},
	}
	GeneralSocksServerFailure = commandNegotiationSocksError("general SOCKS server failure", 0x01)
	HostUnreachable           = commandNegotiationSocksError("Host unreachable", 0x04)
	CommandNotSupported       = commandNegotiationSocksError("Command not supported", 0x07)
	AddressTypeNotSupported   = commandNegotiationSocksError("Address type not supported", 0x08)
)

func commandNegotiationSocksError(msg string, code byte) SocksError {
	return SocksError{
		msg:   msg,
		cache: []byte{Socks5Version, code, Rsv, AddrTypeIPv4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	}
}
