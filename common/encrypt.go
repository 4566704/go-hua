package common

import (
	"crypto/md5"
	"encoding/hex"
	"os"
)

func GetMd5ByStr(str string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(str))
	if err != nil {
		return "", err
	}
	s := hex.EncodeToString(h.Sum(nil))
	return s, nil
}

func GetMd5ByBytes(data []byte) (string, error) {
	h := md5.New()
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	s := hex.EncodeToString(h.Sum(nil))
	return s, nil
}

func GetFileMd5(savePath string) (string, error) {
	buf, err := os.ReadFile(savePath)
	if err != nil {
		return "", err
	}
	str, err := GetMd5ByBytes(buf)
	if err != nil {
		return "", err
	}
	return str, nil
}
