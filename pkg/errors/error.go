package errors

import (
	"syscall"
)

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_INVALID_DATA       = 11
	errnoERROR_IO_PENDING         = 997
	errnoRPC_S_SERVER_UNAVAILABLE = 1722
	errnoEPT_S_NOT_REGISTERED     = 1753
)

var (
	ERROR_INVALID_DATA       error = syscall.Errno(errnoERROR_INVALID_DATA)
	ERROR_IO_PENDING         error = syscall.Errno(errnoERROR_IO_PENDING)
	RPC_S_SERVER_UNAVAILABLE error = syscall.Errno(errnoRPC_S_SERVER_UNAVAILABLE)
	EPT_S_NOT_REGISTERED     error = syscall.Errno(errnoEPT_S_NOT_REGISTERED)
)

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_INVALID_DATA:
		return ERROR_INVALID_DATA
	case errnoERROR_IO_PENDING:
		return ERROR_IO_PENDING
	case errnoRPC_S_SERVER_UNAVAILABLE:
		return RPC_S_SERVER_UNAVAILABLE
	case errnoEPT_S_NOT_REGISTERED:
		return EPT_S_NOT_REGISTERED
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

func Ensure(lastError syscall.Errno) error {
	if lastError != 0 {
		return lastError
	} else {
		return syscall.EINVAL
	}
}

func NotNill(ro uintptr, lastError syscall.Errno) error {
	if ro == 0 {
		return Ensure(lastError)
	}
	return nil

}

func NotZero(lastError syscall.Errno) error {
	if lastError != 0 {
		return Ensure(lastError)
	}
	return nil
}
