package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	zeroUintptr      uintptr = 0
	validClusterName string  = "localhost"
)

func TestOpenCluster(t *testing.T) {
	handle, err := OpenCluster()
	defer handle.Close()

	if err != nil {
		t.Log(err.Error())
	}

	assert.NotZero(t, handle, "handle should not be zero")
	assert.Nil(t, err, "error should be null")
}

func TestOpenRemoteCluster(t *testing.T) {
	handle, err := OpenRemoteCluster(validClusterName)
	defer handle.Close()

	if err != nil {
		t.Log(err.Error())
	}

	assert.NotZero(t, handle, "handle should not be zero")
	assert.Nil(t, err, "error should be null")
}
func TestOpenRemoteCluster_2(t *testing.T) {
	handle, err := OpenRemoteCluster(validClusterName)
	defer handle.Close()
	if err != nil {
		t.Log(err.Error())
	}

	assert.NotZero(t, handle, "handle should not be zero")
	assert.Nil(t, err, "error should be null")
}
