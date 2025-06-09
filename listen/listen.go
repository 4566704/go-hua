package listen

import (
	"net"
)

type Listener struct {
	listener net.Listener
}

func NewListen(network, address string) (*Listener, error) {
	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	l := new(Listener)
	l.listener = listener
	return l, nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	conn, err := l.listener.Accept()
	return conn, err
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	return l.listener.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}
