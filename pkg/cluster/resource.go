package cluster

import (
	"syscall"
	"unsafe"

	"github.com/KnicKnic/go-failovercluster-api/pkg/errors"
	"golang.org/x/sys/windows"
)

var (
	procnativeOpenClusterResource   = clusapi_dll.NewProc("OpenClusterResource")
	procnativeCloseClusterResource  = clusapi_dll.NewProc("CloseClusterResource")
	procnativeGetClusterResourceKey = clusapi_dll.NewProc("GetClusterResourceKey")
)

type (
	ResourceHandle uintptr
)

func openClusterResource(cluster ClusterHandle, resourceName *uint16) (ResourceHandle, error) {
	r0, _, lastError := syscall.Syscall(procnativeOpenClusterResource.Addr(), 2, uintptr(cluster), uintptr(unsafe.Pointer(resourceName)), 0)
	handle := ResourceHandle(r0)
	return handle, errors.NotNill(r0, lastError)
}

func (cluster ClusterHandle) OpenResource(resourceName string) (handle ResourceHandle, err error) {
	rn, err := windows.UTF16PtrFromString(resourceName)
	if err != nil {
		return
	}
	handle, err = openClusterResource(cluster, rn)
	return
}

func closeClusterResource(handle ResourceHandle) error {
	_, _, lastError := syscall.Syscall(procnativeCloseClusterResource.Addr(), 1, uintptr(handle), 0, 0)
	return errors.NotZero(lastError)
}

func (handle ResourceHandle) Close() {
	_ = closeClusterResource(handle)
}

// GetKey gets a cluster registry key for the resource
// for samDesired use syscall.KEY_ALL_ACCESS KEY_READ KEY_WRITE KEY_SET_VALUE
func (handle ResourceHandle) GetKey(samDesired int) (KeyHandle, error) {
	r0, _, lastError := syscall.Syscall(procnativeGetClusterResourceKey.Addr(), 2, uintptr(handle), uintptr(samDesired), 0)
	key := KeyHandle(r0)
	return key, errors.NotNill(r0, lastError)
}
