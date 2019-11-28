package cluster

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

var (
	clusapi_dll = windows.NewLazyDLL("clusapi.dll")
)
