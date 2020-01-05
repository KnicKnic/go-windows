package cluster

import (
	"golang.org/x/sys/windows"
)

var (
	clusapi_dll = windows.NewLazyDLL("clusapi.dll")
)
