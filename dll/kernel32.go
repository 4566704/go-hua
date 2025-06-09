package dll

import "syscall"

var (
	kernel32Dll      *syscall.DLL
	procCreateMutexW *syscall.Proc
	procOpenMutexW   *syscall.Proc
	procReleaseMutex *syscall.Proc
)

func init() {
	kernel32Dll = syscall.MustLoadDLL("kernel32.dll")
	procCreateMutexW = kernel32Dll.MustFindProc("CreateMutexW")
	procOpenMutexW = kernel32Dll.MustFindProc("OpenMutexW")
	procReleaseMutex = kernel32Dll.MustFindProc("ReleaseMutex")
}
