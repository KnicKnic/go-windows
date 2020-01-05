package kernel32

import (
	"syscall"

	"github.com/KnicKnic/go-failovercluster-api/pkg/errors"
	"golang.org/x/sys/windows"
)

// Do the interface allocations only once for common

var (
	modkernel32    = windows.NewLazySystemDLL("kernel32.dll")
	procLocalAlloc = modkernel32.NewProc("LocalAlloc")
	procLocalFree  = modkernel32.NewProc("LocalFree")
)

const (
	LocalAlloc_LPTR uint32 = 0x40
)

func LocalAlloc(length uint32) (ptr uintptr, err error) {
	ptr, _, lastError := syscall.Syscall(procLocalAlloc.Addr(), 2, uintptr(LocalAlloc_LPTR), uintptr(length), 0)
	err = errors.NotNill(ptr, lastError)
	return
}

func LocalFree(mem uintptr) {
	syscall.Syscall(procLocalFree.Addr(), 1, uintptr(mem), 0, 0)
	return
}
