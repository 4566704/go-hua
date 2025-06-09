package bytepacket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Packet struct {
	buffer *bytes.Buffer
}

// NewPacket 新建封包
func NewPacket(b []byte) *Packet {
	return &Packet{buffer: bytes.NewBuffer(b)}
}

// SetData 设置数据
func (p *Packet) SetData(b []byte) {
	p.buffer = bytes.NewBuffer(b)
}

// GetData 获取数据
func (p *Packet) GetData() []byte {
	return p.buffer.Bytes()
}

// ReadInt8 读取int8
func (p *Packet) ReadInt8() int8 {
	var val int8 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteInt8 写入int8
func (p *Packet) WriteInt8(val int8) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadUInt8 读取uint8
func (p *Packet) ReadUInt8() uint8 {
	var val uint8 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteUInt8 写入uint8
func (p *Packet) WriteUInt8(val uint8) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadInt16 读取int16
func (p *Packet) ReadInt16() int16 {
	var val int16 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteInt16 写入int16
func (p *Packet) WriteInt16(val int16) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadUInt16 读取uint16
func (p *Packet) ReadUInt16() uint16 {
	var val uint16 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteUInt16 写入uint16
func (p *Packet) WriteUInt16(val uint16) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadInt32 读取int32
func (p *Packet) ReadInt32() int32 {
	var val int32 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteUInt32 写入uint32
func (p *Packet) WriteUInt32(val uint32) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadUInt32 读取uint32
func (p *Packet) ReadUInt32() uint32 {
	var val uint32 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteInt32 写入int32
func (p *Packet) WriteInt32(val int32) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadInt64 读取int64
func (p *Packet) ReadInt64() int64 {
	var val int64 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteInt64 写入int64
func (p *Packet) WriteInt64(val int64) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadUInt64 读取uint64
func (p *Packet) ReadUInt64() uint64 {
	var val uint64 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteUInt64 写入uint64
func (p *Packet) WriteUInt64(val uint64) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadFloat32 读取float32
func (p *Packet) ReadFloat32() float32 {
	var val float32 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteFloat32 写入float32
func (p *Packet) WriteFloat32(val float32) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadFloat64 读取float64
func (p *Packet) ReadFloat64() float64 {
	var val float64 = 0
	binary.Read(p.buffer, binary.LittleEndian, &val)
	return val
}

// WriteFloat64 写入float64
func (p *Packet) WriteFloat64(val float64) {
	binary.Write(p.buffer, binary.LittleEndian, &val)
}

// ReadStringGbk 读取GBK字符串
func (p *Packet) ReadStringGbk() (string, error) {
	var len int16 = 0
	binary.Read(p.buffer, binary.LittleEndian, &len)
	buf := p.buffer.Next(int(len))
	if len == 0 {
		return "", nil
	}
	if len < 0 {
		return "", fmt.Errorf("invalid string length: %d", len)
	}
	bytes, err := GbkToUtf8(buf)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// WriteStringGbk 写入GBK字符串
func (p *Packet) WriteStringGbk(str string) error {
	bytes, err := Utf8ToGbk([]byte(str))
	if err != nil {
		return err
	}
	var len int16 = int16(len(bytes))
	err = binary.Write(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return err
	}
	_, err = p.buffer.Write(bytes)
	if err != nil {
		return err
	}
	return err
}

// ReadString 读取字符串
func (p *Packet) ReadString() (string, error) {
	var len int16 = 0
	err := binary.Read(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return "", err
	}
	if len == 0 {
		return "", nil
	}
	if len < 0 {
		return "", fmt.Errorf("invalid string length: %d", len)
	}
	buf := p.buffer.Next(int(len))
	return string(buf), nil
}

// WriteString 写入字符串
func (p *Packet) WriteString(str string) error {
	bytes := []byte(str)
	var len int16 = int16(len(bytes))
	err := binary.Write(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return err
	}
	_, err = p.buffer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// ReadBytes 读取字节集
func (p *Packet) ReadBytes() []byte {
	var len int32 = 0
	binary.Read(p.buffer, binary.LittleEndian, &len)
	if len == 0 {
		return nil
	}
	if len < 0 {
		return nil
	}
	buf := p.buffer.Next(int(len))
	return buf
}

// WriteBytes 写入字节集
func (p *Packet) WriteBytes(data []byte) error {
	var len int32 = int32(len(data))
	err := binary.Write(p.buffer, binary.LittleEndian, &len)
	if err != nil {
		return err
	}
	_, err = p.buffer.Write(data)
	if err != nil {
		return err
	}
	return err
}

// ReadElangDateTime 读取易语言日期时间
func (p *Packet) ReadElangDateTime() time.Time {
	var f float64
	binary.Read(p.buffer, binary.LittleEndian, &f)
	t := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Local)
	msStr := fmt.Sprint("+", uint64(math.Floor(f*86400*1000+0.5)), "ms")
	ms, _ := time.ParseDuration(msStr)
	return t.Add(ms)
}

// GbkToUtf8 GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk UTF-8 转 GBK
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
