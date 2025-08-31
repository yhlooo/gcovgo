package raw

import (
	"encoding"
	"encoding/binary"
)

// RecordBlocks 块记录
type RecordBlocks struct {
	Flags []uint32
}

var _ encoding.BinaryUnmarshaler = (*RecordBlocks)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	basic_block: header int32:flags*
func (r *RecordBlocks) UnmarshalBinary(data []byte) error {
	for len(data) >= 4 {
		r.Flags = append(r.Flags, binary.LittleEndian.Uint32(data[:4]))
		data = data[4:]
	}
	return nil
}
