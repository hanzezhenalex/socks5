package protocol

type SocksError struct {
}

func newSocksError() SocksError {
	return SocksError{}
}

func (err SocksError) Error() string {
	return ""
}
