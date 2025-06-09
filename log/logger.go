package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// 前景 背景 颜色
// ---------------------------------------
// 30 40 黑色
// 31 41 红色
// 32 42 绿色
// 33 43 黄色
// 34 44 蓝色
// 35 45 紫红色
// 36 46 青蓝色
// 37 47 白色
//
// 代码 意义
// -------------------------
// 0 终端默认设置
// 1 高亮显示
// 4 使用下划线
// 5 闪烁
// 7 反白显示
// 8 不可见

const (
	Black   int = iota + 30 //黑色
	Red                     //红色
	Green                   //绿色
	Yellow                  //黄色
	Blue                    //蓝色
	Magenta                 //紫红色
	Cyan                    //青蓝色
	White                   //白色
)

func init() {
	logColour["panic"] = logColourType{5, Red + 10, White}
	logColour["fatal"] = logColourType{5, Red + 10, White}
	logColour["error"] = logColourType{5, Red + 10, White}
	logColour["warning"] = logColourType{1, Yellow + White, White}
	logColour["info"] = logColourType{1, Blue + 10, White}
	logColour["debug"] = logColourType{0, Magenta + 10, White}
	logColour["trace"] = logColourType{0, Magenta + 10, White}
}

type logColourType struct {
	d int
	b int
	f int
}

var logColour = make(map[string]logColourType)

type Logger struct {
	mu         sync.Mutex
	out        io.Writer
	isDiscard  int32
	rootDir    string
	filePrefix string
	enableSave bool
	logFiles   map[string]*os.File
	logDir     string
}

func New(out io.Writer) *Logger {
	l := &Logger{
		out:      out,
		logFiles: make(map[string]*os.File),
	}
	if out == io.Discard {
		l.isDiscard = 1
	}
	return l
}

func (l *Logger) SetRootDir(dir string) {
	l.rootDir = dir

}

func (l *Logger) EnableSave(logDir string, filePrefix string) error {
	l.filePrefix = filePrefix

	err := os.MkdirAll(logDir, 0666)
	if err != nil {
		return err
	}
	// 获取当前日期
	now := time.Now()
	// 格式化日期为 YYYYMMDD 的形式
	dateStr := now.Format("20060102")
	logFile, ok := l.logFiles[dateStr]
	if !ok || logFile == nil {
		logFileName := filepath.Join(logDir, l.filePrefix+dateStr+".log")
		// 打开日志文件，如果文件不存在则创建
		logFileHandle, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		l.logFiles[dateStr] = logFileHandle
		l.Debugf("日志文件:%s", logFileName)
		l.enableSave = true
		l.logDir = logDir
	}

	return nil
}

func (l *Logger) Debugf(format string, a ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output("debug", fmt.Sprintf(format, a...))
}

func (l *Logger) Infof(format string, a ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output("info", fmt.Sprintf(format, a...))
}

func (l *Logger) Warnf(format string, a ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output("warning", fmt.Sprintf(format, a...))
}

func (l *Logger) Errorf(format string, a ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	l.output("error", fmt.Sprintf(format, a...))
}

func (l *Logger) output(level string, s string) error {

	_, file, line, _ := runtime.Caller(2)
	if len(l.rootDir) > 0 {
		file = strings.Replace(file, l.rootDir, "", 1)
	} else {
		file = filepath.Base(file)
	}

	add := fmt.Sprintf("%s:%d ", file, line)
	// 获取当前日期
	now := time.Now()
	timestamp := now.Local().Format("2006-01-02 15:04:05")
	levelStr := fmt.Sprintf("\x1b[%d;%d;%dm%-8s\x1b[0m", logColour[level].d, logColour[level].b, logColour[level].f, strings.ToUpper(level))
	msg := fmt.Sprintf("%s  %s  %s\n", timestamp, levelStr, add+s)
	logText := fmt.Sprintf("%s  %s  %s\n", timestamp, strings.ToUpper(level), add+s)
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.out.Write([]byte(msg))
	// 是否记录日志到文件
	if l.enableSave {
		// 格式化日期为 YYYYMMDD 的形式
		dateStr := now.Format("20060102")
		logFile, ok := l.logFiles[dateStr]
		// 文件是否已经在存，是否已经打开
		if ok && logFile != nil {
			// 已经存在 直接写入
			logFile.WriteString(logText)
		} else {
			// 文件不存在，创建新文件
			logFileName := filepath.Join(l.logDir, l.filePrefix+dateStr+".log")
			// 打开日志文件，如果文件不存在则创建
			logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			l.logFiles[dateStr] = logFile
			logFile.WriteString(logText)
		}
	}
	return err
}
