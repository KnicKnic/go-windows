package cluster

import (
	"syscall"
	"unsafe"

	"github.com/KnicKnic/go-failovercluster-api/pkg/errors"
	"golang.org/x/sys/windows"
)

var (
	procnativeOpenCluster  = clusapi_dll.NewProc("OpenCluster")
	procnativeCloseCluster = clusapi_dll.NewProc("CloseCluster")
)

type (
	ClusterHandle uintptr
)

func openCluster() (ClusterHandle, error) {
	r0, _, lastError := syscall.Syscall(procnativeOpenCluster.Addr(), 1, uintptr(0), 0, 0)
	handle := ClusterHandle(r0)
	return handle, errors.NotNill(r0, lastError)
}
func openRemoteCluster(clusterName *uint16) (ClusterHandle, error) {
	r0, _, lastError := syscall.Syscall(procnativeOpenCluster.Addr(), 1, uintptr(unsafe.Pointer(clusterName)), 0, 0)
	handle := ClusterHandle(r0)
	return handle, errors.NotNill(r0, lastError)
}

func OpenCluster() (handle ClusterHandle, err error) {
	handle, err = openCluster()
	return
}
func OpenRemoteCluster(clusterName string) (handle ClusterHandle, err error) {
	cn, err := windows.UTF16PtrFromString(clusterName)
	if err != nil {
		return
	}
	handle, err = openRemoteCluster(cn)
	return
}

func closeCluster(handle ClusterHandle) error {
	_, _, lastError := syscall.Syscall(procnativeCloseCluster.Addr(), 1, uintptr(handle), 0, 0)
	return errors.NotZero(lastError)
}

func (handle ClusterHandle) Close() {
	_ = closeCluster(handle)
}
