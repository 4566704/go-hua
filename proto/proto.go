package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const (
	PakHeadLen int32 = 8
)

// Encode 编码
func Encode(buf []byte) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	var id int32 = 5
	var length = int32(len(buf))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, id)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(buf))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// Decode 解码
func Decode(reader *bufio.Reader) ([]byte, error) {
	// 读取头部
	lengthByte, err := reader.Peek(int(PakHeadLen)) // 读取前4个字节的数据
	if err != nil {
		return nil, err
	}
	lengthBuff := bytes.NewBuffer(lengthByte)
	var id int32 = 0
	var length int32 = 0

	err = binary.Read(lengthBuff, binary.LittleEndian, &id)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	if length <= 0 {
		return nil, errors.New("decode length cannot be less than 0")
	}

	if length > 0x1000000 {
		return nil, errors.New("decode length cannot be greater than 0x1000000")
	}

	var data []byte

	// 剩余读长长度
	unread := int(PakHeadLen + length)
	readLen := 0
	for {
		if unread > 0x10000 {
			readLen = 0x10000
		} else {
			readLen = unread
		}

		_, err = reader.Peek(readLen) // 先窥探数据
		if err != nil {
			return nil, err
		}

		// Buffered返回缓冲中现有的可读取的字节数。
		if reader.Buffered() < readLen {
			return nil, err
		}

		// 读取
		buf := make([]byte, readLen)
		_, err = reader.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf...)

		unread -= readLen
		//fmt.Println("readLen", readLen)
		if unread <= 0 {
			break
		}

	}

	if len(data) < int(PakHeadLen) {
		return nil, nil
	}

	return data[PakHeadLen:], nil
}

// EncodeCrc32 编码 Crc32校验
func EncodeCrc32(buf []byte) ([]byte, error) {
	// 读取消息的长度，转换成int32类型（占4个字节）
	ieee := crc32.NewIEEE()
	ieee.Write(buf)
	var crc = ieee.Sum32()
	var length = int32(len(buf))
	var pkg = new(bytes.Buffer)
	// 写入消息头
	err := binary.Write(pkg, binary.LittleEndian, crc)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(buf))
	if err != nil {
		return nil, err
	}
	return pkg.Bytes(), nil
}

// DecodeCrc32 解码 Crc32校验
func DecodeCrc32(reader *bufio.Reader) ([]byte, error) {
	// 读取头部
	lengthByte, err := reader.Peek(int(PakHeadLen)) // 读取前4个字节的数据
	if err != nil {
		return nil, err
	}
	lengthBuff := bytes.NewBuffer(lengthByte)
	var crc uint32 = 0
	var length int32 = 0

	err = binary.Read(lengthBuff, binary.LittleEndian, &crc)
	if err != nil {
		return nil, err
	}
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	if length <= 0 {
		return nil, errors.New("decode length cannot be less than 0")
	}

	if length > 0x1000000 {
		return nil, errors.New("decode length cannot be greater than 0x1000000")
	}

	var data []byte

	// 剩余读长长度
	unread := int(PakHeadLen + length)
	readLen := 0
	for {
		if unread > 0x1000 {
			readLen = 0x1000
		} else {
			readLen = unread
		}

		_, err = reader.Peek(readLen) // 先窥探数据
		if err != nil {
			return nil, err
		}

		// Buffered返回缓冲中现有的可读取的字节数。
		if reader.Buffered() < readLen {
			return nil, err
		}

		// 读取
		buf := make([]byte, readLen)
		_, err = reader.Read(buf)
		if err != nil {
			return nil, err
		}
		data = append(data, buf...)

		unread -= readLen
		//fmt.Println("readLen", readLen)
		if unread <= 0 {
			break
		}

	}

	if len(data) < int(PakHeadLen) {
		return nil, nil
	}
	ieee := crc32.NewIEEE()
	ieee.Write(data[PakHeadLen:])
	if crc != ieee.Sum32() {
		return nil, errors.New("crc check error")
	}
	return data[PakHeadLen:], nil
}
