package memory

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

type GUID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

func (guid GUID) String() string {
	return fmt.Sprintf("%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		guid.Data1, guid.Data2, guid.Data3,
		guid.Data4[0], guid.Data4[1], guid.Data4[2], guid.Data4[3],
		guid.Data4[4], guid.Data4[5], guid.Data4[6], guid.Data4[7])
}
func GuidFromString(g string) (guid GUID, err error) {
	_, err = fmt.Sscanf(g, "%08X-%04X-%04X-%02X%02X-%02X%02X%02X%02X%02X%02X",
		&guid.Data1, &guid.Data2, &guid.Data3,
		&guid.Data4[0], &guid.Data4[1], &guid.Data4[2], &guid.Data4[3],
		&guid.Data4[4], &guid.Data4[5], &guid.Data4[6], &guid.Data4[7])
	return
}

func GenerateGuidByte() (data []byte, err error) {

	data = make([]byte, 16)
	_, err = rand.Read(data)
	return
}

func GenerateGuid() (guid GUID, err error) {
	data, err := GenerateGuidByte()
	if err != nil {
		return
	}
	guid, err = GuidFromBytes(data)
	return
}

func (data GUID) ToByte() (output []byte, err error) {

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		return
	}
	output = buf.Bytes()
	return
}
func GuidFromBytes(data []byte) (guid GUID, err error) {

	buf := bytes.NewReader(data)
	err = binary.Read(buf, binary.LittleEndian, &guid)
	return
}
