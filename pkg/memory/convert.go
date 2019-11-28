package memory

import (
	"encoding/binary"
	"unsafe"
)

func Uint32ToByte(data uint32) (output []byte) {
	output = make([]byte, int(unsafe.Sizeof(data)))
	binary.LittleEndian.PutUint32(output, data)
	return
}
func Uint64ToByte(data uint64) (output []byte) {
	output = make([]byte, int(unsafe.Sizeof(data)))
	binary.LittleEndian.PutUint64(output, data)
	return
}
