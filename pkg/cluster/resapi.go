package cluster

import (
	"golang.org/x/sys/windows"
)

var (
	resapi_dll = windows.NewLazyDLL("resutils.dll")
)
