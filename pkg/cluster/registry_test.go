package cluster

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/KnicKnic/go-failovercluster-api/pkg/memory"
	"github.com/stretchr/testify/assert"
)

func TestQueryInvalidValue(t *testing.T) {
	handle, err := OpenRemoteCluster(validClusterName)
	defer handle.Close()
	if err != nil {
		t.Log(err.Error())
	}

	assert.NotZero(t, handle, "handle should be null")
	assert.Equal(t, nil, err, "error should be null")

	res, err := handle.OpenResource(validResourceName)
	assert.Equal(t, nil, err, "error should be null")
	defer res.Close()

	key, err := res.GetKey(syscall.KEY_ALL_ACCESS)
	assert.Equal(t, nil, err, "error should be null")
	defer key.Close()

	_, err = key.QueryGuidValue("nonExistant")
	fmt.Println(err)
	assert.NotNil(t, err)
	assert.Equal(t, syscall.ERROR_FILE_NOT_FOUND, err)

	myGuidStr := "206994D6-C7B7-ABDB-D89E-AB9CBF3853C4"
	myGuid, err := memory.GuidFromString(myGuidStr)
	err = key.SetGuidValue("test_guid_value", myGuid)
	assert.Nil(t, err)

	bytes, err := myGuid.ToByte()
	assert.Nil(t, err)

	err = key.SetValue(myGuidStr, syscall.REG_BINARY, bytes)
	assert.Nil(t, err)
}
