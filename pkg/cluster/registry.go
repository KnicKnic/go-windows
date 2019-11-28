package cluster

import (
	"syscall"
	"unsafe"

	"github.com/KnicKnic/go-failovercluster-api/pkg/errors"
	"github.com/KnicKnic/go-failovercluster-api/pkg/memory"
	"golang.org/x/sys/windows"
)

type (
	KeyHandle uintptr
)

const (
	REG_CREATED_NEW_KEY uint32 = 0
)

var (
	procnativeClusterRegCreateKey = clusapi_dll.NewProc("ClusterRegCreateKey")
	// procnativeClusterOpenCreateKey  = clusapi_dll.NewProc("ClusterOpenCreateKey")
	procnativeClusterRegCloseKey   = clusapi_dll.NewProc("ClusterRegCloseKey")
	procnativeClusterRegSetValue   = clusapi_dll.NewProc("ClusterRegSetValue")
	procnativeClusterRegEnumValue  = clusapi_dll.NewProc("ClusterRegEnumValue")
	procnativeClusterRegQueryValue = clusapi_dll.NewProc("ClusterRegQueryValue")
)

func closeClusterKey(handle KeyHandle) error {
	_, _, lastError := syscall.Syscall(procnativeClusterRegCloseKey.Addr(), 1, uintptr(handle), 0, 0)
	return errors.NotZero(lastError)
}

func (handle KeyHandle) Close() {
	_ = closeClusterKey(handle)
}

func clusterRegSetValue(handle KeyHandle, lpszValueName *uint16, dwType uint32, data []byte) error {

	var r0 uintptr
	var dataSize uint32 = uint32(len(data))

	if dataSize == 0 {
		// use dataSize pointer as address for data because why not it won't be looked at
		r0, _, _ = syscall.Syscall6(procnativeClusterRegSetValue.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(lpszValueName)), uintptr(dwType), uintptr(unsafe.Pointer(&dataSize)), uintptr(dataSize), 0)
	} else {
		r0, _, _ = syscall.Syscall6(procnativeClusterRegSetValue.Addr(), 5, uintptr(handle), uintptr(unsafe.Pointer(lpszValueName)), uintptr(dwType), uintptr(unsafe.Pointer(&data[0])), uintptr(dataSize), 0)
	}
	lastError := syscall.Errno(r0)
	return errors.NotZero(lastError)
}

// SetValue sets a value on a key
// for dwType either see "golang.org/x/sys/windows/registry".BINARY (and other values)
// or use syscall.REG_BINARY & other values
func (handle KeyHandle) SetValue(value string, dwType uint32, data []byte) error {
	vn, err := windows.UTF16PtrFromString(value)
	if err != nil {
		return err
	}
	return clusterRegSetValue(handle, vn, dwType, data)
}

// SetByteValue sets a value on a key
func (handle KeyHandle) SetByteValue(value string, data []byte) error {
	return handle.SetValue(value, syscall.REG_BINARY, data)
}

// SetGuidValue sets a value on a key
func (handle KeyHandle) SetGuidValue(value string, guid memory.GUID) error {
	data, err := guid.ToByte()
	if err != nil {
		return err
	}
	return handle.SetByteValue(value, data)
}

func clusterRegCreateKey(handle KeyHandle, lpszKeyName *uint16, samDesired int) (KeyHandle, bool, error) {

	var r0 uintptr
	var disposition uint32

	var keyHandle uintptr

	r0, _, _ = syscall.Syscall9(procnativeClusterRegCreateKey.Addr(),
		7,
		uintptr(handle),
		uintptr(unsafe.Pointer(lpszKeyName)),
		uintptr(0), /*REG_OPTION_NON_VOLATILE*/
		uintptr(samDesired),
		uintptr(0),
		uintptr(unsafe.Pointer(&keyHandle)),
		uintptr(unsafe.Pointer(&disposition)),
		0,
		0)

	lastError := syscall.Errno(r0)
	created := disposition == REG_CREATED_NEW_KEY
	return KeyHandle(keyHandle), created, errors.NotZero(lastError)
}

// CreateKey creates a subkey
// for samDesired use syscall.KEY_ALL_ACCESS KEY_READ KEY_WRITE KEY_SET_VALUE
func (handle KeyHandle) CreateKey(keyName string, samDesired int) (key KeyHandle, created bool, err error) {
	kn, err := windows.UTF16PtrFromString(keyName)
	if err != nil {
		return
	}
	key, created, err = clusterRegCreateKey(handle, kn, samDesired)
	return
}

// clusterRegEnumValue
func clusterRegEnumValue(handle KeyHandle, index uint32) (keyName string, dwType uint32, data []byte, err error) {

	nameCCh := uint32(50)
	dataCB := uint32(248)

	lastError := uintptr(syscall.ERROR_MORE_DATA)
	var keyNameArr []uint16

	for lastError == uintptr(syscall.ERROR_MORE_DATA) {
		// increase values to ensure not zero & space for extra nulls
		nameCCh += 2
		dataCB += 8
		data = make([]byte, dataCB)
		keyNameArr = make([]uint16, nameCCh)
		lastError, _, _ = syscall.Syscall9(procnativeClusterRegEnumValue.Addr(),
			7,
			uintptr(handle),
			uintptr(index),
			uintptr(unsafe.Pointer(&keyNameArr[0])),
			uintptr(unsafe.Pointer(&nameCCh)),
			uintptr(unsafe.Pointer(&dwType)),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(unsafe.Pointer(&dataCB)),
			0,
			0)
	}

	err = errors.NotZero(syscall.Errno(lastError))
	if err != nil {
		return
	}
	// resize arrays to appropriate return sizes
	data = append([]byte(nil), data[:dataCB]...)
	// add 1 for null
	keyNameArr = append([]uint16(nil), keyNameArr[:nameCCh+1]...)

	keyName = windows.UTF16ToString(keyNameArr)

	return
}

// also need to add batches
// Test if batch returns error when violate a condition

// need query value

func clusterRegQueryValue(handle KeyHandle, value *uint16) (dwType uint32, data []byte, err error) {

	dataCB := uint32(248)

	lastError := uintptr(syscall.ERROR_MORE_DATA)

	for lastError == uintptr(syscall.ERROR_MORE_DATA) {
		// increase values to ensure not zero & space for extra nulls
		dataCB += 8
		data = make([]byte, dataCB)
		lastError, _, _ = syscall.Syscall6(procnativeClusterRegQueryValue.Addr(),
			5,
			uintptr(handle),
			uintptr(unsafe.Pointer(value)),
			uintptr(unsafe.Pointer(&dwType)),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(unsafe.Pointer(&dataCB)),
			0)
	}

	err = errors.NotZero(syscall.Errno(lastError))
	if err != nil {
		return
	}
	// resize arrays to appropriate return sizes
	data = append([]byte(nil), data[:dataCB]...)

	return
}

// QueryValue returns syscall.ERROR_FILE_NOT_FOUND if value does not exist
// for dwType either see "golang.org/x/sys/windows/registry".BINARY (and other values)
// or use syscall.REG_BINARY & other values
func (handle KeyHandle) QueryValue(valueName string) (dwType uint32, data []byte, err error) {
	vn, err := windows.UTF16PtrFromString(valueName)
	if err != nil {
		return
	}

	dwType, data, err = clusterRegQueryValue(handle, vn)
	return
}

// QueryByteValue returns syscall.ERROR_FILE_NOT_FOUND if value does not exist
func (handle KeyHandle) QueryByteValue(valueName string) (data []byte, err error) {
	dwType, data, err := handle.QueryValue(valueName)
	if err != nil {
		return
	}
	if dwType != syscall.REG_BINARY {
		err = errors.ERROR_INVALID_DATA
		return
	}
	return
}

// QueryGuidValue returns syscall.ERROR_FILE_NOT_FOUND if value does not exist
func (handle KeyHandle) QueryGuidValue(valueName string) (data memory.GUID, err error) {
	dataBuf, err := handle.QueryByteValue(valueName)
	if err != nil {
		return
	}
	if len(dataBuf) != int(unsafe.Sizeof(data)) {
		err = errors.ERROR_INVALID_DATA
		return
	}

	data, err = memory.GuidFromBytes(dataBuf)
	return
}
