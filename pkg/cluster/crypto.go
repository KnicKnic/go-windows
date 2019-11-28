package cluster

import (
	"syscall"
	"unsafe"

	"github.com/KnicKnic/go-failovercluster-api/pkg/errors"
	"github.com/KnicKnic/go-failovercluster-api/pkg/kernel32"
	"github.com/KnicKnic/go-failovercluster-api/pkg/ntdll"
	"golang.org/x/sys/windows"
)

var _ unsafe.Pointer

const (
	PROV_RSA_FULL      CryptographicServiceProviderType = 1
	PROV_RSA_SIG       CryptographicServiceProviderType = 2
	PROV_DSS           CryptographicServiceProviderType = 3
	PROV_FORTEZZA      CryptographicServiceProviderType = 4
	PROV_MS_EXCHANGE   CryptographicServiceProviderType = 5
	PROV_SSL           CryptographicServiceProviderType = 6
	PROV_RSA_SCHANNEL  CryptographicServiceProviderType = 12
	PROV_DSS_DH        CryptographicServiceProviderType = 13
	PROV_EC_ECDSA_SIG  CryptographicServiceProviderType = 14
	PROV_EC_ECNRA_SIG  CryptographicServiceProviderType = 15
	PROV_EC_ECDSA_FULL CryptographicServiceProviderType = 16
	PROV_EC_ECNRA_FULL CryptographicServiceProviderType = 17
	PROV_DH_SCHANNEL   CryptographicServiceProviderType = 18
	PROV_SPYRUS_LYNKS  CryptographicServiceProviderType = 20
	PROV_RNG           CryptographicServiceProviderType = 21
	PROV_INTEL_SEC     CryptographicServiceProviderType = 22
	PROV_REPLACE_OWF   CryptographicServiceProviderType = 23
	PROV_RSA_AES       CryptographicServiceProviderType = 24

	MS_ENH_RSA_AES_PROV string = "Microsoft Enhanced RSA and AES Cryptographic Provider"

	CLUS_CREATE_CRYPT_CONTAINER_NOT_FOUND OpenClusterCryptProviderFlags = 1
	CLUS_CREATE_CRYPT_NONE                OpenClusterCryptProviderFlags = 0
)

type (
	HCLUSCRYPTPROVIDER               uintptr
	CryptographicServiceProviderType uint32
	OpenClusterCryptProviderFlags    uint32
)

var (
	procOpenClusterCryptProvider   = resapi_dll.NewProc("OpenClusterCryptProvider")
	procCloseClusterCryptProvider  = resapi_dll.NewProc("CloseClusterCryptProvider")
	procClusterEncrypt             = resapi_dll.NewProc("ClusterEncrypt")
	procClusterDecrypt             = resapi_dll.NewProc("ClusterDecrypt")
	procFreeClusterCrypt           = resapi_dll.NewProc("FreeClusterCrypt")
	procOpenClusterCryptProviderEx = resapi_dll.NewProc("OpenClusterCryptProviderEx")
)

func openClusterCryptProvider(lpszResource *uint16, lpszProvider *uint16, dwType CryptographicServiceProviderType, dwFlags OpenClusterCryptProviderFlags) (HCLUSCRYPTPROVIDER, error) {
	r0, _, lastError := syscall.Syscall6(procOpenClusterCryptProvider.Addr(), 4, uintptr(unsafe.Pointer(lpszResource)), uintptr(unsafe.Pointer(lpszProvider)), uintptr(dwType), uintptr(dwFlags), 0, 0)
	handle := HCLUSCRYPTPROVIDER(r0)
	return handle, errors.NotNill(r0, lastError)
}

func openClusterCryptProviderEx(lpszResource *uint16, lpszKeyName *uint16, lpszProvider *uint16, dwType CryptographicServiceProviderType, dwFlags OpenClusterCryptProviderFlags) (HCLUSCRYPTPROVIDER, error) {
	r0, _, lastError := syscall.Syscall6(procOpenClusterCryptProviderEx.Addr(), 5, uintptr(unsafe.Pointer(lpszResource)), uintptr(unsafe.Pointer(lpszKeyName)), uintptr(unsafe.Pointer(lpszProvider)), uintptr(dwType), uintptr(dwFlags), 0)
	handle := HCLUSCRYPTPROVIDER(r0)
	return handle, errors.NotNill(r0, lastError)
}

func OpenClusterCryptProvider(Resource string, Provider string, dwType CryptographicServiceProviderType, dwFlags OpenClusterCryptProviderFlags) (handle HCLUSCRYPTPROVIDER, err error) {
	resource, err := windows.UTF16PtrFromString(Resource)
	if err != nil {
		return
	}
	provider, err := windows.UTF16PtrFromString(Provider)
	if err != nil {
		return
	}

	handle, err = openClusterCryptProvider(resource, provider, dwType, dwFlags)
	return
}

func (handle HCLUSCRYPTPROVIDER) CloseClusterCryptProvider() {
	syscall.Syscall(procCloseClusterCryptProvider.Addr(), 1, uintptr(handle), 0, 0)
	return
}

func encryptDecrypt(encryptDecryptFunc *windows.LazyProc, handle HCLUSCRYPTPROVIDER, data []byte) (encrypted []byte, err error) {

	// doing this as api is not clear if I need a valid pointer for 0 sized memory
	cData, err := ntdll.MemcpyLocalAlloc(data)
	if err != nil {
		return
	}
	defer kernel32.LocalFree(cData)

	dataSize := uint32(len(data))

	var cDest uintptr
	var destSize uint32

	r0, _, _ := syscall.Syscall6(encryptDecryptFunc.Addr(), 5, uintptr(handle), cData, uintptr(dataSize), uintptr(unsafe.Pointer(&cDest)), uintptr(unsafe.Pointer(&destSize)), 0)
	err = errors.NotZero(syscall.Errno(r0))
	if err != nil {
		return
	}
	defer freeClusterCrypt(cDest)

	encrypted = make([]byte, destSize)
	ntdll.MemcpySrcC(encrypted, cDest, uint64(destSize))

	return
}

func (handle HCLUSCRYPTPROVIDER) ClusterEncrypt(data []byte) ([]byte, error) {

	return encryptDecrypt(procClusterEncrypt, handle, data)
}
func (handle HCLUSCRYPTPROVIDER) ClusterDecrypt(data []byte) ([]byte, error) {

	return encryptDecrypt(procClusterDecrypt, handle, data)
}

func freeClusterCrypt(ptr uintptr) {
	_, _, _ = syscall.Syscall(procFreeClusterCrypt.Addr(), 1, uintptr(ptr), 0, 0)
}
