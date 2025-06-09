package dll

import (
	"fmt"
	"github.com/lxn/win"
	"syscall"
	"unsafe"
)

func CreateMutex(bInitialOwner bool, mutexName string) (syscall.Handle, error) {

	handle, _, err := syscall.SyscallN(procCreateMutexW.Addr(), 0, uintptr(win.BoolToBOOL(bInitialOwner)), uintptr(unsafe.Pointer(&mutexName)))
	if err != 0 {
		println(err.Error())

		lastError := win.GetLastError()
		fmt.Println("lastError", lastError)
		if lastError == 183 {
			fmt.Println("Mutex 已经存在")
		}

		return 0, err
	}
	//defer syscall.CloseHandle(syscall.Handle(m))
	if syscall.Handle(handle) != 0 {
		lastError := win.GetLastError()
		if lastError == 183 {
			//fmt.Println("Mutex 已经存在")
		}

	}
	// do something with the mutex
	return syscall.Handle(handle), nil
}

const MUTANT_ALL_ACCESS = 2031617

func OpenMutex(bInheritHandle bool, mutexName string) (syscall.Handle, error) {
	desiredAccess := MUTANT_ALL_ACCESS
	handle, _, err := syscall.SyscallN(procOpenMutexW.Addr(), uintptr(desiredAccess), uintptr(win.BoolToBOOL(bInheritHandle)), uintptr(unsafe.Pointer(&mutexName)))
	if err != 0 {
		return 0, err
	}
	defer syscall.CloseHandle(syscall.Handle(handle))
	if syscall.Handle(handle) == 0 {
		e2 := win.GetLastError()
		if e2 == 183 {
			fmt.Println("打开失败")
		}
		fmt.Println(e2)

	} // do something with the mutex
	fmt.Println("Mutex open", syscall.Handle(handle))
	return syscall.Handle(handle), nil
}

func ReleaseMutex(handle syscall.Handle) bool {
	ret, _, e1 := syscall.SyscallN(procReleaseMutex.Addr(), uintptr(unsafe.Pointer(&handle)))
	if e1 != 0 {
		println(e1.Error())
		return false
	}
	return ret != 0
}
