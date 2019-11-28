package cluster

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	validResourceName string = "r1"
)

func TestCrypto(t *testing.T) {
	var data = []byte{1, 2, 3, 4, 5}
	fmt.Println("Unencrypted: ", data)
	handle, err := OpenClusterCryptProvider(validResourceName, MS_ENH_RSA_AES_PROV, PROV_RSA_AES, CLUS_CREATE_CRYPT_CONTAINER_NOT_FOUND)
	defer handle.CloseClusterCryptProvider()

	if err != nil {
		t.Log(err.Error())
	}

	assert.NotZero(t, handle, "handle should not be null")

	encrypted, err := handle.ClusterEncrypt(data)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Encrypted: ", encrypted)

	decrypted, err := handle.ClusterDecrypt(encrypted)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("Decrypted: ", decrypted)

	clus, err := OpenCluster()
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("clus:", clus)
	defer clus.Close()

	res, err := clus.OpenResource(validResourceName)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)
	fmt.Println("res:", res)
	defer res.Close()

	key, err := res.GetKey(syscall.KEY_ALL_ACCESS)
	if err != nil {
		t.Log(err.Error())
	}
	fmt.Println("key:", key)
	assert.Nil(t, err)
	defer key.Close()

	err = key.SetValue("BeforeEncrypt", syscall.REG_BINARY, data)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)

	err = key.SetValue("AfterEncrypt", syscall.REG_BINARY, encrypted)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)

	err = key.SetValue("Decrypted", syscall.REG_BINARY, decrypted)
	if err != nil {
		t.Log(err.Error())
	}
	assert.Nil(t, err)

}
