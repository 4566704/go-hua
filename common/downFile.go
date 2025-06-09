package common

import (
	"compress/gzip"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadFile(url string, downPath string, fb func(length, downLen int64)) error {
	var (
		fileSize int64
		buf      = make([]byte, 32*1024)
		written  int64
	)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	//req.Header.Add("Authorization", "BNbmgAAGI155F6MJ3N2Tk9ruL_6XQpx-uxkkg:tGCY3xCsgybHd5IjcDMi9yZXBvcy93aF9mbG93RGF0YVNvdXJjZTEiLCJleHBpcmVzIjoxNTM2NzU4NjQ3LCJjb250ZW5VudFR5cGUiOiIiLCJoZWFkZXJzIjoiIiwibWV0aG9kIjoiR0VUIn0=")
	req.Header.Add("Accept-Encoding", "gzip")
	//req.Header.Set("Content-Encoding", "gzip")
	//req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1)")
	req.Header.Set("Accept", " */*")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	client.Timeout = time.Second * 300
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	//读取服务器返回的文件大小
	fileSize = resp.ContentLength
	//fmt.Println(resp.Header)
	if fileSize == -1 {
		fmt.Println("取文件长度失败", url)
	}

	folderPath := filepath.Dir(downPath)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return err
	}
	//创建文件
	file, err := os.Create(downPath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	body := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			fmt.Println("http resp unzip is failed,err: ", err)
		}
	}

	for {
		//读取bytes
		nr, er := body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		fb(fileSize, written)
	}
	return err
}
