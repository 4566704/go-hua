package httpurl

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// 复用的 http.Transport 实例
var transport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

// 复用的 http.Client 实例
var client = &http.Client{
	Transport: transport,
	Timeout:   time.Second * 15,
}

func Post(url string, param []byte, header map[string]string) ([]byte, error) {
	payload := bytes.NewReader(param)
	// 创建新的 HTTP POST 请求
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	// 添加请求头
	for k, v := range header {
		req.Header.Add(k, v)
	}
	if header == nil || header["Content-Type"] == "" {
		req.Header.Add("Content-Type", "application/json")
	}
	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保响应体在函数结束时关闭
	defer res.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 检查响应状态码
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return body, nil
}

func Get(httpUrl string, header map[string]string) ([]byte, error) {
	// 创建新的 HTTP GET 请求
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return nil, err
	}
	// 添加请求头
	for k, v := range header {
		req.Header.Add(k, v)
	}
	req.Header.Add("Content-Type", "application/json")
	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保响应体在函数结束时关闭
	defer res.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 检查响应状态码
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return body, nil
}

func GetProxy(httpUrl string, header map[string]string, proxyAddr string) ([]byte, error) {

	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	// req.Header.Add("Content-Type", "application/json")
	/*req.Header.Add("Authorization", "BNbmgAAGI155F6MJ3N2Tk9ruL_6XQpx-uxkkg:tGCY3xCsgybHd5IjcDMi9yZXBvcy93aF9mbG93RGF0YVNvdXJjZTEiLCJleHBpcmVzIjoxNTM2NzU4NjQ3LCJjb250ZW5VudFR5cGUiOiIiLCJoZWFkZXJzIjoiIiwibWV0aG9kIjoiR0VUIn0=")*/

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if proxyAddr != "" {
		proxyUrl, err := url.Parse(proxyAddr)
		if err != nil {
			return nil, err
		}
		tr.Proxy = http.ProxyURL(proxyUrl)
	}

	client := &http.Client{Transport: tr}
	client.Timeout = time.Second * 15
	defer client.CloseIdleConnections()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//fmt.Println(res)
	//fmt.Println(string(body))
	if res.StatusCode != 200 {
		return body, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return body, nil

}

func Put(url string, param []byte, header map[string]string) ([]byte, error) {
	payload := bytes.NewReader(param)
	// 创建新的 HTTP GET 请求
	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, err
	}
	// 添加请求头
	for k, v := range header {
		req.Header.Add(k, v)
	}
	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保响应体在函数结束时关闭
	defer res.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 检查响应状态码
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return body, nil
}

func Delete(url string, param []byte, header map[string]string) ([]byte, error) {
	payload := bytes.NewReader(param)
	// 创建新的 HTTP DELETE 请求
	req, err := http.NewRequest("DELETE", url, payload)
	if err != nil {
		return nil, err
	}
	// 添加请求头
	for k, v := range header {
		req.Header.Add(k, v)
	}
	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// 确保响应体在函数结束时关闭
	defer res.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	// 检查响应状态码
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}
	return body, nil
}
