package conn

import (
	"hua/flow"
	"hua/rate"
	"net"
	"sync"
	"time"
)

// 主要就是为了统计 带宽，流量

type Conn struct {
	Conn net.Conn   // 多路复用流
	Rate *rate.Rate //带宽
	Flow *flow.Flow //流量
	wg   *sync.WaitGroup
	once sync.Once
}

func NewConn(conn net.Conn, rate *rate.Rate, flow *flow.Flow, wg *sync.WaitGroup) *Conn {
	c := new(Conn)
	c.Conn = conn
	c.Rate = rate
	c.Flow = flow
	c.wg = wg
	return c
}

func (c *Conn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if c.Rate != nil {
		c.Rate.Add(n)
	}
	if c.Flow != nil {
		c.Flow.Add(int64(n), 0)
	}
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if c.Rate != nil {
		c.Rate.Add(n)
	}
	if c.Flow != nil {
		c.Flow.Add(int64(n), 0)
	}
	return
}

func (c *Conn) Close() (err error) {
	// 只执行一次
	c.once.Do(func() {
		if c.wg != nil {
			c.wg.Done()
		}
		err = c.Conn.Close()
	})
	return
}

func (c *Conn) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}
