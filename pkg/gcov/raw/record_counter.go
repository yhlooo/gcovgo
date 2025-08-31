package raw

import (
	"encoding"
	"encoding/binary"
)

// RecordCounter 计数器记录
type RecordCounter struct {
	Counts []uint64
}

var _ encoding.BinaryUnmarshaler = (*RecordCounter)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	counts: header int64:count*
func (r *RecordCounter) UnmarshalBinary(data []byte) error {
	for len(data) >= 8 {
		r.Counts = append(r.Counts, binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
	}
	return nil
}
