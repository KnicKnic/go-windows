package ntdll

import (
	"syscall"
	"unsafe"

	"github.com/KnicKnic/go-windows/pkg/kernel32"
	"golang.org/x/sys/windows"
)

// Do the interface allocations only once for common

var (
	modntdll   = windows.NewLazySystemDLL("ntdll.dll")
	procmemcpy = modntdll.NewProc("memcpy")
)

func Memcpy(dest uintptr, src uintptr, size uint64) (ptr uintptr) {
	r0, _, _ := syscall.Syscall(procmemcpy.Addr(), 3, uintptr(dest), uintptr(src), uintptr(size))
	ptr = uintptr(r0)
	return
}

func MemcpyDestC(dest uintptr, src []byte, size uint64) {
	if size != 0 {
		_, _, _ = syscall.Syscall(procmemcpy.Addr(), 3, uintptr(dest), uintptr(unsafe.Pointer(&src[0])), uintptr(size))
	}
}

func MemcpySrcC(dest []byte, src uintptr, size uint64) {
	if size != 0 {
		_, _, _ = syscall.Syscall(procmemcpy.Addr(), 3, uintptr(unsafe.Pointer(&dest[0])), uintptr(src), uintptr(size))
	}
}

func MemcpyLocalAlloc(data []byte) (ptr uintptr, err error) {

	size := uint32(len(data))
	ptr, err = kernel32.LocalAlloc(size)
	if err != nil {
		return
	}
	MemcpyDestC(ptr, data, uint64(size))
	return
}
