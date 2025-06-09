package spmux

import (
	"errors"
	"net"
	"sync/atomic"
)

type PortListener struct {
	ch      chan *PortConn
	addr    net.Addr
	isClose int32
}

func NewPortListener(connCh chan *PortConn, addr net.Addr) *PortListener {
	return &PortListener{
		ch:      connCh,
		addr:    addr,
		isClose: 0,
	}
}

func (pListener *PortListener) Accept() (net.Conn, error) {
	if atomic.LoadInt32(&pListener.isClose) == 1 {
		return nil, errors.New("监听已经关闭")
	}
	if pListener.ch == nil {
		return nil, errors.New("监听已经关闭")
	}
	conn := <-pListener.ch
	if conn != nil {
		return conn, nil
	}
	return nil, errors.New("the listener has closed")
}

func (pListener *PortListener) Close() error {
	//close
	if atomic.LoadInt32(&pListener.isClose) == 1 {
		return errors.New("监听已经关闭")
	}
	atomic.StoreInt32(&pListener.isClose, 1)
	return nil
}

func (pListener *PortListener) Addr() net.Addr {
	return pListener.addr
}
