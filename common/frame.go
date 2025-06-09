package common

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const FrameVersion byte = 1

// 帧头部
const (
	sizeOfVer    = 1
	sizeOfCmd    = 1
	sizeOfLength = 2
	sizeOfSid    = 4
	headerSize   = sizeOfVer + sizeOfCmd + sizeOfSid + sizeOfLength
)

const (
	frameBufSize = 0x10000
)

type Frame struct {
	ver  byte
	cmd  byte
	sid  uint32
	data []byte
}

func NewFrame(version byte, cmd byte, sid uint32) Frame {
	return Frame{ver: version, cmd: cmd, sid: sid}
}

func ReadFrame(conn net.Conn) (Frame, error) {
	f := NewFrame(0, 0, 0)
	hdr := RawHeader{}
	// 读取头
	_, err := conn.Read(hdr[:])
	if err != nil {
		return f, err
	}

	f.ver = hdr.Version()
	f.cmd = hdr.Cmd()
	f.sid = hdr.StreamID()
	// 读取数据
	reader := bufio.NewReaderSize(conn, frameBufSize)
	// 窥探数据长度
	_, err = reader.Peek(int(hdr.Length()))
	if err != nil {
		return f, err
	}
	// 返回可以从当前缓冲区读取的字节数。
	if reader.Buffered() < int(hdr.Length()) {
		return f, err
	}

	data := make([]byte, int(hdr.Length()))
	// 读取数据
	_, err = reader.Read(data)
	if err != nil {
		return f, err
	}
	f.data = data
	return f, err
}

func WriteFrame(conn net.Conn, frame Frame) error {
	hdr := RawHeader{}
	hdr[0] = frame.ver
	hdr[1] = frame.cmd
	binary.LittleEndian.PutUint16(hdr[2:], uint16(len(frame.data)))
	binary.LittleEndian.PutUint32(hdr[4:], uint32(frame.sid))

	buf := hdr.Bytes()
	buf = append(buf, frame.data...)
	_, err := conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (f *Frame) Version() byte {
	return f.ver
}

func (f *Frame) Cmd() byte {
	return f.cmd
}

func (f *Frame) Length() uint16 {
	return uint16(len(f.data))
}

func (f *Frame) Data() []byte {
	return f.data
}

func (f *Frame) StreamID() uint32 {
	return f.sid
}

func (f *Frame) SetData(data []byte) error {
	if len(data) > frameBufSize-1 {
		fmt.Errorf("超出最大长度:%d", frameBufSize-1)
	}
	f.data = data
	return nil
}

func (f *Frame) Marshal(v any) error {
	buf, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return f.SetData(buf)
}

func (f *Frame) Unmarshal(v any) error {
	if f.data == nil || len(f.data) == 0 {
		return errors.New("没有数据")
	}
	return json.Unmarshal(f.data, v)
}

func (f *Frame) String() string {
	return fmt.Sprintf("Version:%d Cmd:%d StreamID:%d Length:%d",
		f.Version(), f.Cmd(), f.StreamID(), f.Length())
}

type RawHeader [headerSize]byte

func (h RawHeader) Version() byte {
	return h[0]
}

func (h RawHeader) Cmd() byte {
	return h[1]
}

func (h RawHeader) Length() uint16 {
	return binary.LittleEndian.Uint16(h[2:])
}

func (h RawHeader) StreamID() uint32 {
	return binary.LittleEndian.Uint32(h[4:])
}

func (h RawHeader) Bytes() []byte {
	buf := make([]byte, headerSize)
	for i, _ := range buf {
		buf[i] = h[i]
	}
	return buf
}
