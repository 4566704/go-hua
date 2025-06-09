package common

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func CreateRandomStringHex(len int) string {
	var container string
	var str = "ABCDEF1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func CreateRandomStringInt(len int) string {
	var container string
	var str = "1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

// IsMatchDynamic 动态规划；时间复杂度O(n^2)，空间复杂度O(n^2)
func IsMatchDynamic(s string, p string) bool {
	n, m := len(s), len(p)
	dp := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]bool, m+1)
	}
	dp[0][0] = true
	for i := 1; i <= m; i++ {
		if p[i-1] == '*' {
			// 可以匹配任意字符串(包括空字符串)
			dp[0][i] = true
		} else {
			break
		}
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if p[j-1] == '*' {
				// dp[i][j-1]=>不使用这个*，dp[i-1][j]=>使用这个*
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else if p[j-1] == '?' || s[i-1] == p[j-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[n][m]
}

// IsMatchGreed 贪心；时间复杂度O(n^2)，空间复杂度O(1)
func IsMatchGreed(s string, p string) bool {
	i, j := 0, 0
	start, last := 0, 0
	for i = 0; i < len(s); {
		if j < len(p) && (s[i] == p[j] || p[j] == '?') {
			i++
			j++
		} else if j < len(p) && p[j] == '*' {
			last = i
			// 记录s下一个的位置
			start = j + 1
			// 记录*下一个的位置
			j++
		} else if start != 0 {
			last++
			i = last
			// 更新到s记录位置的下一个
			j = start
		} else {
			return false
		}
	}
	for ; j < len(p) && p[j] == '*'; j++ {

	}
	return j == len(p)
}

// GetBetweenStr 取中间字符串
func GetBetweenStr(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	if ending == "" {
		return str[s:]
	}
	e := strings.Index(str[s:], ending)

	if e < 0 {
		return str[s:]
	}
	return str[s : s+e]
}

// GetGUID 使用来自 crypto/rand 包的 rand.Read 函数来生成基本的 UUID。
func GetGUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// GetUUID 使用来自 crypto/rand 包的 rand.Read 函数来生成基本的 UUID。
func GetUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	uuid := fmt.Sprintf("%x%x%x%x%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

// CheckMobile 检验手机号
func CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)
}

// PhoneCheckRule 手机号验证规则
func PhoneCheckRule() string {
	return "^1[3|4|5|6|7|8|9][0-9]\\d{8}$"
}

// EmailCheckRule 邮箱验证规则
func EmailCheckRule() string {
	return "^[A-Z0-9._%+-]+@[A-Z0-9.-]+\\.[A-Z]{2,6}$"
}

// CheckIdCard 检验身份证
func CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// Buffer 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (b *Buffer) Append(i interface{}) *Buffer {
	switch val := i.(type) {
	case int:
		b.append(strconv.Itoa(val))
	case int64:
		b.append(strconv.FormatInt(val, 10))
	case uint:
		b.append(strconv.FormatUint(uint64(val), 10))
	case uint64:
		b.append(strconv.FormatUint(val, 10))
	case string:
		b.append(val)
	case []byte:
		b.Write(val)
	case rune:
		b.WriteRune(val)
	}
	return b
}

func (b *Buffer) append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			log.Println("*****内存不够了！******")
		}
	}()
	b.WriteString(s)
	return b
}

// Camel2Case 驼峰式写法转为下划线写法
func Camel2Case(name string) string {
	buffer := NewBuffer()
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.Append('_')
			}
			buffer.Append(unicode.ToLower(r))
		} else {
			buffer.Append(r)
		}
	}
	return buffer.String()
}
