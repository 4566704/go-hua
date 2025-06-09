package spmux

import (
	"bufio"
	"net"
	"time"
)

type PortConn struct {
	conn   net.Conn
	reader *bufio.Reader
}

func newPortConn(conn net.Conn) *PortConn {
	c := &PortConn{
		conn: conn,
	}
	c.reader = bufio.NewReaderSize(conn, 0x1000)
	return c
}

// Peek 窥探数据
func (p *PortConn) Peek(n int) ([]byte, error) {
	return p.reader.Peek(n)
}

// Buffered 从当前缓冲区读取的字节数
func (p *PortConn) Buffered() int {
	return p.reader.Buffered()
}

func (p *PortConn) Read(b []byte) (n int, err error) {
	return p.reader.Read(b)
}

func (p *PortConn) Write(b []byte) (n int, err error) {
	return p.conn.Write(b)
}

func (p *PortConn) Close() error {
	return p.conn.Close()
}

func (p *PortConn) LocalAddr() net.Addr {
	return p.conn.LocalAddr()
}

func (p *PortConn) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *PortConn) SetDeadline(t time.Time) error {
	return p.conn.SetDeadline(t)
}

func (p *PortConn) SetReadDeadline(t time.Time) error {
	return p.conn.SetReadDeadline(t)
}

func (p *PortConn) SetWriteDeadline(t time.Time) error {
	return p.conn.SetWriteDeadline(t)
}
