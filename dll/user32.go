//go:build (windows && amd64) || (windows && 386)

package dll

import (
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

var (
	user32Dll           *syscall.DLL
	procGetWindowTextW  *syscall.Proc
	procGetWindowClassW *syscall.Proc
	procFindWindowExW   *syscall.Proc
)

func init() {
	user32Dll = syscall.MustLoadDLL("user32.dll")
	procGetWindowTextW = user32Dll.MustFindProc("GetWindowTextW")
	procGetWindowClassW = user32Dll.MustFindProc("GetClassNameW")
	procFindWindowExW = user32Dll.MustFindProc("FindWindowExW")
}

func GetWindowText(hWnd win.HWND, str *uint16, maxCount int32) (len int32, err error) {
	//r0, _, e1 := procGetWindowTextW.Call(uintptr(hWnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hWnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetWindowTextString(hWnd win.HWND) (string, error) {
	buf := make([]uint16, 255)
	_, err := GetWindowText(hWnd, &buf[0], int32(len(buf)))
	if err != nil {
		return "", err
	}
	windowsTitle := syscall.UTF16ToString(buf)
	return windowsTitle, nil
}

func GetTopFolderPath() (string, error) {
	//lName := uint16(0)
	hWnd := win.FindWindow(syscall.StringToUTF16Ptr("CabinetWClass"), nil)
	if hWnd == 0 {
		return "", nil
	}
	hWnd = GetWindowFromClassName(hWnd, "Breadcrumb Parent")
	if hWnd == 0 {
		return "", nil
	}
	hWnd = FindWindowExW(hWnd, 0, syscall.StringToUTF16Ptr("ToolbarWindow32"), nil)
	if hWnd == 0 {
		return "", nil
	}
	buf := make([]uint16, 255)
	_, err := GetWindowText(hWnd, &buf[0], int32(len(buf)))
	if err != nil {
		return "", err
	}
	title := syscall.UTF16ToString(buf)
	return title, nil
}

func FindWindowExW(parent win.HWND, child win.HWND, lpClassName, lpWindowName *uint16) win.HWND {
	ret, _, _ := procFindWindowExW.Call(uintptr(parent), uintptr(child),
		uintptr(unsafe.Pointer(lpClassName)),
		uintptr(unsafe.Pointer(lpWindowName)),
		0)

	return win.HWND(ret)
}

// GetWindowFromClassName 使用类名获取窗口句柄 递归
func GetWindowFromClassName(parentWnd win.HWND, className string) win.HWND {
	// 寻找顶层窗口
	hWnd := FindWindowExW(parentWnd, 0, nil, nil)

	for hWnd != 0 {
		// 获取窗口类名
		name, _ := GetWindowClassString(hWnd)
		if name == className {
			return hWnd
		}
		// 递归查找
		phWnd := GetWindowFromClassName(hWnd, className)
		if phWnd != 0 {
			return phWnd
		}
		// 查找下一子窗口
		hWnd = FindWindowExW(parentWnd, hWnd, nil, nil)
	}
	return 0
}

// GetWindowClassString 取窗口类名 string
func GetWindowClassString(hWnd win.HWND) (string, error) {
	buf := make([]uint16, 255)
	_, err := GetWindowClass(hWnd, &buf[0], int32(len(buf)))
	if err != nil {
		return "", err
	}
	return syscall.UTF16ToString(buf), nil
}

// GetWindowClass 取窗口类名 API
func GetWindowClass(hWnd win.HWND, str *uint16, maxCount int32) (len int32, err error) {
	//r0, _, e1 := procGetWindowTextW.Call(uintptr(hWnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	r0, _, e1 := syscall.Syscall(procGetWindowClassW.Addr(), 3, uintptr(hWnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
