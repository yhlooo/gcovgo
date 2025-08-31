package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// RecordArcs 弧记录
type RecordArcs struct {
	// 块编号
	BlockNo uint32
	// 块往外连接的弧
	Arcs []Arc
}

var _ encoding.BinaryUnmarshaler = (*RecordArcs)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	header int32:block_no arc*
func (r *RecordArcs) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return newDataTooShortError(len(data), 4, "block_no")
	}
	r.BlockNo = binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	r.Arcs = make([]Arc, len(data)/8)
	for i := range r.Arcs {
		if err := r.Arcs[i].UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unmarshal arc %d error: %w", i, err)
		}
		data = data[8:]
	}

	return nil
}

// Arc 弧
type Arc struct {
	// 目标块编号
	DestBlock uint32
	// 弧属性
	Flags uint32 // TODO: 应该是个 bitmap
}

var _ encoding.BinaryUnmarshaler = (*Arc)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	int32:dest_block int32:flags
func (arc *Arc) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return newDataTooShortError(len(data), 8, "dest_block and flags")
	}
	arc.DestBlock = binary.LittleEndian.Uint32(data[:4])
	arc.Flags = binary.LittleEndian.Uint32(data[4:8])

	return nil
}
