package memory

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuid(t *testing.T) {

	guid, err := GenerateGuid()
	assert.Nil(t, err)
	fmt.Println(guid)

	bytes, err := guid.ToByte()
	assert.Nil(t, err)

	newGuid, err := GuidFromBytes(bytes)
	assert.Nil(t, err)
	assert.Equal(t, guid, newGuid)
	fmt.Println(newGuid)

	str := guid.String()

	guidFromStr, err := GuidFromString(str)
	assert.Nil(t, err)
	assert.Equal(t, guid, guidFromStr)
	fmt.Println(guidFromStr)

}
