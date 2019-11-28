package cluster

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

var (
	resapi_dll = windows.NewLazyDLL("resutils.dll")
)
