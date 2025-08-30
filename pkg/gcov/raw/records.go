package raw

import (
	"encoding"
	"encoding/binary"
)

// Record 记录
type Record struct {
	// 记录类型标签
	Tag RecordTag
	// 记录数据长度（ uint32 的个数）
	Length uint32

	// 基于 Tag ，以下成员仅有一个值不为 nil

	// 原始数据，当 Tag 无法处理时有值
	Raw *RecordRaw `json:",omitempty"`
}

var _ encoding.BinaryUnmarshaler = (*Record)(nil)

// UnmarshalBinary 从二进制反序列化
func (r *Record) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return newDataTooShortError(len(data), 8, "tag and length")
	}

	r.Tag = RecordTag(binary.LittleEndian.Uint32(data[:4]))
	r.Length = binary.LittleEndian.Uint32(data[4:8])
	data = data[8:]

	if len(data) < int(r.Length)*4 {
		return newDataTooShortError(len(data), int(r.Length)*4, "items")
	}

	switch r.Tag {
	default:
		r.Raw = &RecordRaw{Data: make([]byte, 4*int(r.Length))}
		copy(r.Raw.Data, data)
	}

	return nil
}

// RecordTag 记录类型标签
type RecordTag uint32

// RecordRaw 记录原始数据
type RecordRaw struct {
	Data []byte
}
