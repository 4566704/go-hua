package ipquery

import (
	"errors"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"net"
	"regexp"
	"strings"
)

var cBuff []byte

func LoadIp(dbPath string) error {
	// 1、从 dbPath 加载整个 xdb 到内存
	var err error
	cBuff, err = xdb.LoadContentFromFile(dbPath)
	if err != nil {
		return err
	}
	return nil
}

type Data struct {
	Ip       string `json:"ip"`
	Country  string `json:"country"`
	Area     string `json:"area"`
	Province string `json:"province"`
	City     string `json:"city"`
	Isp      string `json:"isp"`
}

// RegexpDns 匹配域名
func RegexpDns(str string) error {
	reg, err := regexp.Compile(`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z]{0,62})\.?`)
	if err != nil {
		return err
	}
	s := reg.FindAllString(str, -1)
	if len(s) > 0 {
		return nil
	}
	return errors.New("不是域名")
}

func QueryString(ip string) string {
	data, err := Query(ip)
	if err != nil {
		return ""
	}
	if data.Country == "中国" {
		return fmt.Sprintf("%s%s %s", data.Province, data.City, data.Isp)
	}
	return fmt.Sprintf("%s%s %s", data.Country, data.Province, data.Isp)
}

func Query(ip string) (Data, error) {
	ipData := Data{}

	_, err := xdb.CheckIP(ip)
	if err != nil {
		err = RegexpDns(ip)
		if err != nil {
			return ipData, errors.New("不是IP或域名")
		}
		ipAddr, err := net.ResolveIPAddr("ip", ip)

		if err != nil {
			return ipData, errors.New("域名无法解析")
		}
		ipData.Ip = ipAddr.String()
	} else {
		ipData.Ip = ip
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	// 备注：并发使用，每个 goroutine 需要创建一个独立的 searcher 对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		return ipData, err
	}

	defer searcher.Close()

	// 进行搜索
	//var tStart = time.Now()
	region, err := searcher.SearchByStr(ipData.Ip)
	if err != nil {
		return ipData, err
	}

	//fmt.Printf("位置: %s, 耗时: %s\n", region, time.Since(tStart))
	if region == "" {
		return ipData, errors.New("搜索结果是空的")
	}

	arr := strings.Split(region, "|")
	if len(arr) < 5 {
		return ipData, errors.New("搜索结果字段不足")
	}

	ipData.Country = arr[0]
	ipData.Area = arr[1]
	ipData.Province = arr[2]
	ipData.City = arr[3]
	ipData.Isp = arr[4]

	if ipData.Country == "0" {
		ipData.Country = ""
	}
	if ipData.Area == "0" {
		ipData.Area = ""
	}
	if ipData.Province == "0" {
		ipData.Province = ""
	}
	if ipData.City == "0" {
		ipData.City = ""
	}
	if ipData.Isp == "0" {
		ipData.Isp = ""
	}
	return ipData, nil
}

func StringToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		//fmt.Println("Invalid IP address")
		return 0, errors.New("Invalid IP address")
	}

	ipBytes := ip.To4()
	if ipBytes == nil {
		//fmt.Println("Not an IPv4 address")
		return 0, errors.New("Not an IPv4 address")
	}

	ipNum := (uint32(ipBytes[0]) << 24) + (uint32(ipBytes[1]) << 16) + (uint32(ipBytes[2]) << 8) + uint32(ipBytes[3])
	//fmt.Println("IP address as number:", ipNum)
	return ipNum, nil
}
