package dll

import (
	"fmt"
	"github.com/lxn/win"
	"syscall"
)

func ShellExecute(dir string, file string, args string, hide bool) bool {
	lpOperation, _ := syscall.UTF16PtrFromString("open")
	lpFile, _ := syscall.UTF16PtrFromString(file)
	lpCwd, _ := syscall.UTF16PtrFromString(dir)
	lpArgs, _ := syscall.UTF16PtrFromString(args)
	showCmd := win.SW_SHOWNORMAL
	if hide {
		showCmd = win.SW_HIDE
	}
	return win.ShellExecute(0, lpOperation, lpFile, lpArgs, lpCwd, showCmd)
}

func test() {
	//ShellExecute("", "https://www.baidu.com", false)
}

func CreateProcess() {
	// 定义CreateProcess参数
	var si syscall.StartupInfo
	var pi syscall.ProcessInformation

	// 获取可执行文件路径
	exe := "C:\\Windows\\System32\\notepad.exe"

	// 转换为LPCTSTR类型
	exePtr, err := syscall.UTF16PtrFromString(exe)
	if err != nil {
		panic(err)
	}

	// 调用CreateProcess函数
	err = syscall.CreateProcess(
		exePtr, // 可执行文件路径
		nil,    // 命令行参数
		nil,    // 进程安全描述符
		nil,    // 线程安全描述符
		false,  // 是否继承句柄
		0,      // 创建标志
		nil,    // 环境变量
		nil,    // 当前目录
		&si,    // StartupInfo结构体指针
		&pi,    // ProcessInformation结构体指针
	)
	if err != nil {
		panic(err)
	}

	// 关闭句柄
	_ = syscall.CloseHandle(pi.Process)
	_ = syscall.CloseHandle(pi.Thread)

	fmt.Println("Process ID:", pi.ProcessId)
}
