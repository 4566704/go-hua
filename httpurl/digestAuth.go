package httpurl

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type RequestArg struct {
	Host     string //http://localhost
	Uri      string
	Method   string
	Header   map[string]string //请求头
	PostBody []byte
	Username string //用户名
	Password string //密码
}

// DigestAuthRequest digestAuth 请求
func DigestAuthRequest(arg *RequestArg) ([]byte, error) {
	url := arg.Host + arg.Uri
	req, err := http.NewRequest(arg.Method, url, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		resStr, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return resStr, nil
	}
	parts := digestParts(resp)
	parts["uri"] = arg.Uri
	parts["method"] = arg.Method
	parts["username"] = arg.Username
	parts["password"] = arg.Password
	parts["algorithm"] = "MD5"
	req, err = http.NewRequest(arg.Method, url, bytes.NewBuffer(arg.PostBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", getDigestAuthorization(parts))
	for k, v := range arg.Header {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, err
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return all, nil
}

// 第一次请求的响应header 获取第二次请求所需要的信息 Www-Authenticate
func digestParts(resp *http.Response) map[string]string {
	result := map[string]string{}
	if len(resp.Header["Www-Authenticate"]) > 0 {
		wantedHeaders := []string{"nonce", "realm", "qop", "opaque"}
		responseHeaders := strings.Split(resp.Header["Www-Authenticate"][0], ",")
		for _, r := range responseHeaders {
			for _, w := range wantedHeaders {
				if strings.Contains(r, w) {
					result[w] = strings.Split(r, `"`)[1]
				}
			}
		}
	}
	return result
}

// 第二次请求的请求header Authorization 值
func getDigestAuthorization(digestParts map[string]string) string {
	d := digestParts
	//ha1=md5(username:realm:password)
	ha1 := getMD5(d["username"] + ":" + d["realm"] + ":" + d["password"])
	//ha2=md5(method:uri)
	ha2 := getMD5(d["method"] + ":" + d["uri"])
	nonceCount := "00000001"
	cnonce := getCnonce()

	//response=md5(ha1:nonce:nc:cnonce:qop:ha2)
	response := getMD5(fmt.Sprintf("%s:%s:%s:%s:%s:%s", ha1, d["nonce"], nonceCount, cnonce, d["qop"], ha2))

	//qop和nc的值不能加引号
	authorization := fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", response="%s", opaque="%s", algorithm=MD5", qop=%s, nc=%s, cnonce="%s"`,
		d["username"], d["realm"], d["nonce"], d["uri"], response, d["opaque"], d["qop"], nonceCount, cnonce)

	return authorization
}

// 字符串MD5加密
func getMD5(text string) string {
	hash := md5.New()
	hash.Write([]byte(text))
	return hex.EncodeToString(hash.Sum(nil))
}

// 获取 cnonce
// 客户端提供的不透明带引号的字符串值，客户端和服务器都使用它来避免选定的明文攻击、提供相互身份验证以及提供一些消息完整性保护
func getCnonce() string {
	b := make([]byte, 8)
	io.ReadFull(rand.Reader, b)
	return fmt.Sprintf("%x", b)[:16]
}
