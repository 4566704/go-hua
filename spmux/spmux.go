package spmux

import (
	"fmt"
	"hua/http"
	"hua/log"
	"net"
	"sync"
	"time"
)

const AcceptTimeout = 100
const ChannelBufferSize = 1
const FirstPeekSize = 1

type PortMux struct {
	listener  net.Listener
	logger    *log.Logger
	once      sync.Once
	port      int
	tcpConn   chan *PortConn
	httpConn  chan *PortConn
	httpsConn chan *PortConn
}

func NewPortMux(port int, logger *log.Logger) *PortMux {
	p := &PortMux{
		port:      port,
		logger:    logger,
		tcpConn:   make(chan *PortConn),
		httpConn:  make(chan *PortConn),
		httpsConn: make(chan *PortConn),
	}
	return p
}

func (p *PortMux) Start() error {
	addr := fmt.Sprintf("0.0.0.0:%d", p.port)
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	p.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	go p.listen()
	return nil
}

func (p *PortMux) listen() {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			p.logger.Errorf("接受连接失败 端口:%d 错误:%s", p.port, err.Error())
			p.Close()
			break
		}
		p.logger.Debugf("新连接 地址:%s", conn.RemoteAddr().String())
		go p.process(conn)
	}
}

func (p *PortMux) process(conn net.Conn) {
	newConn := newPortConn(conn)
	buf, err := newConn.Peek(FirstPeekSize)
	if err != nil {
		conn.Close()
		return
	}

	n := newConn.Buffered()
	if n >= 8 {
		buf, err = newConn.Peek(n)
		if err != nil {
			conn.Close()
			return
		}
	}

	if n >= 8 && http.IsHttp(buf) {
		timer := time.NewTimer(AcceptTimeout)
		select {
		case <-timer.C:
			conn.Close()
			p.logger.Errorf("接受超时 端口:%d 协议:%s", p.port, "http")
		case p.httpConn <- newConn:
		}
	} else if n >= 8 && http.IsHttps(buf) {
		timer := time.NewTimer(AcceptTimeout)
		select {
		case <-timer.C:
			conn.Close()
			p.logger.Errorf("接受超时 端口:%d 协议:%s", p.port, "https")
		case p.httpsConn <- newConn:
		}
	} else {
		timer := time.NewTimer(AcceptTimeout)
		select {
		case <-timer.C:
			conn.Close()
			p.logger.Errorf("接受超时 端口:%d 协议:%s", p.port, "tcp")
		case p.tcpConn <- newConn:
		}
	}
}

func (p *PortMux) Close() (err error) {
	p.once.Do(func() {
		close(p.tcpConn)
		close(p.httpsConn)
		close(p.httpConn)
		err = p.listener.Close()
	})
	return
}

// TcpListener TCP监听
// Deprecated: 一些网络通讯(RFB协议)是由服务器先发送数据的，会导致无法识别成TCP
func (p *PortMux) TcpListener() net.Listener {
	return NewPortListener(p.tcpConn, p.listener.Addr())
}

func (p *PortMux) HttpListener() net.Listener {
	return NewPortListener(p.httpConn, p.listener.Addr())
}

func (p *PortMux) HttpsListener() net.Listener {
	return NewPortListener(p.httpsConn, p.listener.Addr())
}
