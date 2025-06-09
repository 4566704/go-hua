package nettest

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"
)

type ICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

type Ping struct {
	addr    string
	port    int
	timeout int64
	icmp    ICMP
	size    int
}

func NewPing(addr string, port int, timeout int64) *Ping {
	p := new(Ping)
	p.addr = addr
	p.port = port
	p.timeout = timeout
	p.size = 32

	// icmp
	p.icmp.Type = 8
	p.icmp.Code = 0
	p.icmp.Checksum = 0
	p.icmp.Identifier = 1
	p.icmp.SequenceNum = 0
	return p
}

func (p *Ping) Set(addr string, port int, timeout int64) {
	p.addr = addr
	p.port = port
	p.timeout = timeout
}

func (p *Ping) Test() int {
	et := 0
	if p.port == 0 {
		et = p.icmpEcho()
	} else {
		et = p.portCheck()
	}
	return et
}

func (p *Ping) icmpEcho() int {

	conn, err := net.DialTimeout("ip:icmp", p.addr, time.Duration(p.timeout)*time.Millisecond)
	if err != nil {
		log.Println(err)
		return -1
	}

	defer conn.Close()

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, p.icmp) // 以大端模式写入
	data := make([]byte, p.size)                    //
	buffer.Write(data)
	data = buffer.Bytes()

	p.icmp.SequenceNum++
	// 检验和设为0
	data[2] = byte(0)
	data[3] = byte(0)

	data[6] = byte(p.icmp.SequenceNum >> 8)
	data[7] = byte(p.icmp.SequenceNum)
	p.icmp.Checksum = CheckSum(data)
	data[2] = byte(p.icmp.Checksum >> 8)
	data[3] = byte(p.icmp.Checksum)

	// 开始时间
	t1 := time.Now()
	err = conn.SetDeadline(t1.Add(time.Duration(time.Duration(p.timeout) * time.Millisecond)))
	if err != nil {
		//log.Println(err)
		return -1
	}
	_, err = conn.Write(data)
	if err != nil {
		//log.Println(err)
		return -1
	}
	buf := make([]byte, 65535)
	_, err = conn.Read(buf)
	if err != nil {
		//fmt.Println("请求超时。")
		return -1
	}
	et := int(time.Since(t1) / 1000000)
	return et
}

func CheckSum(data []byte) uint16 {
	var sum uint32
	var length = len(data)
	var index int
	for length > 1 { // 溢出部分直接去除
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length == 1 {
		sum += uint32(data[index])
	}
	// CheckSum的值是16位，计算是将高16位加低16位，得到的结果进行重复以该方式进行计算，直到高16位为0
	/*
	   sum的最大情况是：ffffffff
	   第一次高16位+低16位：ffff + ffff = 1fffe
	   第二次高16位+低16位：0001 + fffe = ffff
	   即推出一个结论，只要第一次高16位+低16位的结果，再进行之前的计算结果用到高16位+低16位，即可处理溢出情况
	*/
	sum = uint32(sum>>16) + uint32(sum)
	sum = uint32(sum>>16) + uint32(sum)
	return uint16(^sum)
}

func (p *Ping) portCheck() int {
	// 开始时间
	t1 := time.Now()
	conn, err := net.DialTimeout("tcp", p.addr+":"+strconv.Itoa(p.port), time.Duration(p.timeout)*time.Millisecond)
	et := int(time.Since(t1) / 1000000)
	if err != nil {
		//log.Println(err)
		//fmt.Println("请求超时。")
		return -1
	}
	defer conn.Close()
	return et
}
