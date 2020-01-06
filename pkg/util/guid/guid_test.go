package guid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuid(t *testing.T) {

	guid, err := Generate()
	assert.Nil(t, err)
	fmt.Println(guid)

	bytes, err := guid.ToByte()
	assert.Nil(t, err)

	newGUID, err := FromBytes(bytes)
	assert.Nil(t, err)
	assert.Equal(t, guid, newGUID)
	fmt.Println(newGUID)

	str := guid.String()

	guidFromStr, err := FromString(str)
	assert.Nil(t, err)
	assert.Equal(t, guid, guidFromStr)
	fmt.Println(guidFromStr)

}
