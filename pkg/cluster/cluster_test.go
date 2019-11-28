package cluster

import (
	"testing"

	"github.com/KnicKnic/go-failovercluster-api/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	zeroUintptr      uintptr = 0
	validClusterName string  = "TestCluster"
)

func TestOpenCluster(t *testing.T) {
	handle, err := OpenCluster()
	defer handle.Close()

	if err != nil {
		t.Log(err.Error())
	}

	assert.Zero(t, handle, "handle should be null")
	assert.Equal(t, errors.EPT_S_NOT_REGISTERED, err, "error should not be null")
}

func TestOpenRemoteCluster(t *testing.T) {
	handle, err := OpenRemoteCluster(".")
	defer handle.Close()

	if err != nil {
		t.Log(err.Error())
	}

	assert.Zero(t, handle, "handle should be null")
	assert.Equal(t, errors.RPC_S_SERVER_UNAVAILABLE, err, "error should not be null")
}
func TestOpenRemoteCluster_2(t *testing.T) {
	handle, err := OpenRemoteCluster(validClusterName)
	defer handle.Close()
	if err != nil {
		t.Log(err.Error())
	}

	assert.NotZero(t, handle, "handle should be null")
	assert.Equal(t, nil, err, "error should be null")
}
