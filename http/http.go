package http

import (
	"bytes"
	"fmt"
	"strings"
)

func IsHttps(data []byte) bool {
	if len(data) == 0 {
		return false
	} else if data[0] != 0x16 {
		return false
	} else if data[1] != 0x3 {
		return false
	} else if data[3] != 0x1 && data[3] != 0x2 && data[3] != 0x3 && data[3] != 0x4 {
		return false
	} else if !bytes.ContainsAny(data, "http") {
		return false
	}
	return true
}

func IsHttp(data []byte) bool {
	fmt.Print(string(data))
	if len(data) < 8 {
		fmt.Println("<8")
		return false
	} else if strings.EqualFold(string(data[:3]), "GET") {
		fmt.Println("GET")
		return true
	} else if strings.EqualFold(string(data[:4]), "HEAD") {
		fmt.Println("HEAD")
		return true
	} else if strings.EqualFold(string(data[:4]), "POST") {
		fmt.Println("POST")
		return true
	} else if strings.EqualFold(string(data[:3]), "PUT") {
		fmt.Println("PUT")
		return true
	} else if strings.EqualFold(string(data[:6]), "DELETE") {
		fmt.Println("DELETE")
		return true
	} else if strings.EqualFold(string(data[:7]), "CONNECT") {
		fmt.Println("CONNECT")
		return true
	} else if strings.EqualFold(string(data[:7]), "OPTIONS") {
		fmt.Println("OPTIONS")
		return true
	} else if strings.EqualFold(string(data[:6]), "TRACE") {
		fmt.Println("TRACE", string(data[:6]))
		return true
	}
	return false
}
