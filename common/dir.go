package common

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// 获取工作目录
func GetWorkDir() string {
	str, _ := os.Getwd()
	return str
}

func trimmedPath(file string) string {
	idx := strings.LastIndexByte(file, '/')
	if idx == -1 {
		return file
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(file[:idx], '/')
	if idx == -1 {
		return file
	}
	return file[:idx+1]
}

// 最终方案-全兼容
func GetRunDir() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Contains(dir, getTmpDir()) {
		dir = GetWorkDir()
	}
	if dir[len(dir)-1:] != "/" || dir[len(dir)-1:] != "\\" {
		dir += string(filepath.Separator)
	}
	return dir
}

// 获取系统临时目录，兼容go run
func getTmpDir() string {
	dir := os.Getenv("TEMP")
	if dir == "" {
		dir = os.Getenv("TMP")
	}
	res, _ := filepath.EvalSymlinks(dir)
	return res
}

// 获取当前执行文件绝对路径 (go build)
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		return ""
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	abPath = trimmedPath(filename)
	return abPath
}
