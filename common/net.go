package common

import (
	"net"
)

const NetBufferSize = 0x1000

func GetLocalIp() ([]string, error) {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	ipList := make([]string, 0)
	for _, address := range addrList {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipList = append(ipList, ipNet.IP.String())
			}
		}
	}
	return ipList, err
}

// GetLocalMac 获取本机的MAC地址
func GetLocalMac() ([]string, error) {
	macList := make([]string, 0)
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, inter := range interfaces {
		macList = append(macList, inter.HardwareAddr.String())
	}
	return macList, nil
}
