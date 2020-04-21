package cluster

import (
	"crypto/rand"
	"fmt"
	"reflect"
	"syscall"
	"testing"

	"github.com/KnicKnic/go-windows/pkg/util/guid"
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
	myGuid, err := guid.FromString(myGuidStr)
	err = key.SetGuidValue("test_guid_value", myGuid)
	assert.Nil(t, err)

	bytes, err := myGuid.ToByte()
	assert.Nil(t, err)

	err = key.SetValue(myGuidStr, syscall.REG_BINARY, bytes)
	assert.Nil(t, err)
}

func TestLoadValues(t *testing.T) {
	clusterHandle, err := OpenCluster()
	assert.Nil(t, err)
	defer clusterHandle.Close()

	// open resource
	resourceHandle, err := clusterHandle.OpenResource(validResourceName)
	assert.Nil(t, err)
	defer resourceHandle.Close()

	// open the root key of the registry
	rootKeyHandle, err := resourceHandle.GetKey(syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer rootKeyHandle.Close()

	// Create a test key
	key, _, err := rootKeyHandle.CreateKey("test", syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer key.Close()

	// set a random value
	value := make([]byte, 8)
	data := make([]byte, 8)
	rand.Read(value)
	rand.Read(data)
	valueString := fmt.Sprintf("%x", value)

	key.SetByteValue(valueString, data)

	// Call Load Values
	values, err := key.LoadValues()
	assert.Nil(t, err)

	// check that it is returned
	if retValue, ok := values[valueString]; !ok {
		t.Errorf("Expected value is not present")
	} else {
		if !reflect.DeepEqual(data, retValue.Data) {
			t.Errorf("Loaded value do not match input")
		}
	}

	// Delete the value
	err = key.DeleteValue(valueString)
	assert.Nil(t, err)
}

func TestDeleteValues(t *testing.T) {
	clusterHandle, err := OpenCluster()
	assert.Nil(t, err)
	defer clusterHandle.Close()

	// open resource
	resourceHandle, err := clusterHandle.OpenResource(validResourceName)
	assert.Nil(t, err)
	defer resourceHandle.Close()

	// open the root key of the registry
	rootKeyHandle, err := resourceHandle.GetKey(syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer rootKeyHandle.Close()

	// Create a test key
	key, _, err := rootKeyHandle.CreateKey("test", syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer key.Close()

	// set a random value
	value := make([]byte, 8)
	data := make([]byte, 8)
	rand.Read(value)
	rand.Read(data)
	valueString := fmt.Sprintf("%x", value)

	key.SetByteValue(valueString, data)

	// Call Load Values
	values, err := key.LoadValues()
	assert.Nil(t, err)

	// check that it is returned
	if retValue, ok := values[valueString]; !ok {
		t.Errorf("Expected value is not present")
	} else {
		if !reflect.DeepEqual(data, retValue.Data) {
			t.Errorf("Loaded value do not match input")
		}
	}

	// Delete the value
	err = key.DeleteValue(valueString)
	assert.Nil(t, err)

	// check that it is not returned
	values, err = key.LoadValues()
	assert.Nil(t, err)

	if _, ok := values[valueString]; ok {
		t.Errorf("Deleted value is present")
	}
}

func TestCreateBatch(t *testing.T) {
	clusterHandle, err := OpenCluster()
	assert.Nil(t, err)
	defer clusterHandle.Close()

	// open resource
	resourceHandle, err := clusterHandle.OpenResource(validResourceName)
	assert.Nil(t, err)
	defer resourceHandle.Close()

	// open the root key of the registry
	rootKeyHandle, err := resourceHandle.GetKey(syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer rootKeyHandle.Close()

	// Create a test key
	key, _, err := rootKeyHandle.CreateKey("test", syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer key.Close()

	// Create a test key
	subkey, _, err := key.CreateKey("test-subkey", syscall.KEY_ALL_ACCESS)
	assert.Nil(t, err)
	defer key.Close()

	t.Run("ConditionIsEqual", func(t *testing.T) {
		BatchTestConditionIsEqual(t, key, subkey)
	})

	t.Run("ConditionNotExists", func(t *testing.T) {
		BatchTestConditionNotExists(t, key)
	})

	t.Run("CloseWithoutExecuting", func(t *testing.T) {
		BatchTestCloseWithoutExecuting(t, key)
	})

	// Delete the values created in the test
	err = key.DeleteValue("guid")
	assert.Nil(t, err)
}

func BatchTestConditionIsEqual(t *testing.T, key KeyHandle, subkey KeyHandle) {
	data := make([]byte, 8)
	rand.Read(data)
	err := key.SetByteValue("guid", data)
	assert.Nil(t, err)

	testID := make([]byte, 8)
	testData := make([]byte, 8)
	rand.Read(testID)
	rand.Read(testData)
	testIDString := fmt.Sprintf("%x", testID)

	batchHandle, err := key.CreateBatch()
	assert.Nil(t, err)

	err = batchHandle.BatchAddCommand(CLUSREG_CONDITION_IS_EQUAL, "guid", syscall.REG_BINARY, data)
	assert.Nil(t, err)
	err = batchHandle.BatchAddCommand(CLUSREG_CREATE_KEY, "test-subkey", 0, nil)
	assert.Nil(t, err)
	batchHandle.BatchAddCommand(CLUSREG_SET_VALUE, testIDString, syscall.REG_BINARY, testData)
	assert.Nil(t, err)
	err, _ = batchHandle.CloseBatch(true)
	assert.Nil(t, err)

	// Call Load Values
	values, err := subkey.LoadValues()
	assert.Nil(t, err)

	// check that it is present and equal to the set value
	if retValue, ok := values[testIDString]; !ok {
		t.Errorf("Expected value not present")
	} else {
		if !reflect.DeepEqual(testData, retValue.Data) {
			t.Errorf("Loaded value does not match input")
		}
	}
	err = subkey.DeleteValue(testIDString)
	assert.Nil(t, err)
}

func BatchTestConditionNotExists(t *testing.T, key KeyHandle) {
	// Call Load Values
	values, err := key.LoadValues()
	assert.Nil(t, err)
	if _, ok := values["guid"]; ok {
		err = key.DeleteValue("guid")
		assert.Nil(t, err)
	}

	testData := make([]byte, 8)
	rand.Read(testData)

	batchHandle, err := key.CreateBatch()
	assert.Nil(t, err)
	err = batchHandle.BatchAddCommand(CLUSREG_CONDITION_NOT_EXISTS, "guid", 0, nil)
	assert.Nil(t, err)
	err = batchHandle.BatchAddCommand(CLUSREG_SET_VALUE, "guid", syscall.REG_BINARY, testData)
	assert.Nil(t, err)
	err, _ = batchHandle.CloseBatch(true)
	assert.Nil(t, err)

	// Call Load Values
	values, err = key.LoadValues()
	assert.Nil(t, err)

	// check that it is returned
	if retValue, ok := values["guid"]; !ok {
		t.Errorf("Expected value not present")
	} else {
		if !reflect.DeepEqual(testData, retValue.Data) {
			t.Errorf("Loaded value do not match input")
		}
	}
}

func BatchTestCloseWithoutExecuting(t *testing.T, key KeyHandle) {
	data := make([]byte, 8)
	rand.Read(data)
	err := key.SetByteValue("guid", data)
	assert.Nil(t, err)

	testData := make([]byte, 8)
	rand.Read(testData)

	batchHandle, err := key.CreateBatch()
	assert.Nil(t, err)
	err = batchHandle.BatchAddCommand(CLUSREG_SET_VALUE, "guid", syscall.REG_BINARY, testData)
	assert.Nil(t, err)
	err, _ = batchHandle.CloseBatch(false)
	assert.Nil(t, err)

	// Call Load Values
	values, err := key.LoadValues()
	assert.Nil(t, err)

	// check that it is returned
	if retValue, ok := values["guid"]; !ok {
		t.Errorf("Expected value not present")
	} else {
		if !reflect.DeepEqual(data, retValue.Data) {
			t.Errorf("Loaded value do not match input")
		}
	}
}
