package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// Record 记录
type Record struct {
	// 记录类型标签
	Tag RecordTag
	// 记录数据长度（ uint32 的个数）
	Length uint32

	// 基于 Tag ，以下成员仅有一个值不为 nil

	// 函数，当 Tag 为 TagFunction 时有值
	Function *RecordFunction `json:",omitempty"`
	// 块，当 Tag 为 TagBlocks 时有值
	Blocks *RecordBlocks `json:",omitempty"`
	// 弧，当 Tag 为 TagArcs 时有值
	Arcs *RecordArcs `json:",omitempty"`
	// 行，当 Tag 为 TagLines 时有值
	Lines *RecordLines `json:",omitempty"`
	// 程序摘要，当 Tag 为 TagProgramSummary 时有值
	ProgramSummary *RecordProgramSummary `json:",omitempty"`
	// 计数器，当 Tag 为 TagCounter 时有值
	Counter *RecordCounter `json:",omitempty"`
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

	if len(data) < int(r.Length)*4 {
		return newDataTooShortError(len(data), int(r.Length)*4, "items")
	}
	data = data[8 : 8+int(r.Length)*4]

	var recordData encoding.BinaryUnmarshaler
	switch r.Tag {
	case TagFunction:
		r.Function = &RecordFunction{}
		recordData = r.Function
	case TagBlocks:
		r.Blocks = &RecordBlocks{}
		recordData = r.Blocks
	case TagArcs:
		r.Arcs = &RecordArcs{}
		recordData = r.Arcs
	case TagLines:
		r.Lines = &RecordLines{}
		recordData = r.Lines
	case TagProgramSummary:
		r.ProgramSummary = &RecordProgramSummary{}
		recordData = r.ProgramSummary
	case TagCounter:
		r.Counter = &RecordCounter{}
		recordData = r.Counter
	default:
		r.Raw = &RecordRaw{}
		recordData = r.Raw
	}

	if err := recordData.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("unmarshal %s record error: %w", r.Tag, err)
	}

	return nil
}

// RecordTag 记录类型标签
type RecordTag uint32

var _ fmt.Stringer = RecordTag(0)
var _ encoding.TextMarshaler = RecordTag(0)

const (
	// 通用记录类型
	// 以[01..3f] 开头

	TagFunction RecordTag = 0x01000000
	TagBlocks   RecordTag = 0x01410000
	TagArcs     RecordTag = 0x01430000
	TagLines    RecordTag = 0x01450000
	TagCounter  RecordTag = 0x01a10000

	// Note 的记录类型
	// 以 [41..9f] 开头

	// Data 的记录类型
	// 以 [a1..ff] 开头

	TagObjectSummary  RecordTag = 0xa1000000
	TagProgramSummary RecordTag = 0xa3000000
	TagAfdoFileNames  RecordTag = 0xaa000000
	TagAfdoFunction   RecordTag = 0xac000000
	TagAfdoWorkingSet RecordTag = 0xaf000000
)

// String 返回字符串表示
func (tag RecordTag) String() string {
	switch tag {
	case TagFunction:
		return "Function"
	case TagBlocks:
		return "Blocks"
	case TagArcs:
		return "Arcs"
	case TagLines:
		return "Lines"
	case TagCounter:
		return "Counter"
	case TagObjectSummary:
		return "ObjectSummary"
	case TagProgramSummary:
		return "ProgramSummary"
	case TagAfdoFileNames:
		return "AfdoFileNames"
	case TagAfdoFunction:
		return "AfdoFunction"
	case TagAfdoWorkingSet:
		return "AfdoWorkingSet"
	}

	return fmt.Sprintf("0x%08x", uint32(tag))
}

// MarshalText 序列化为文本
func (tag RecordTag) MarshalText() ([]byte, error) {
	return []byte(tag.String()), nil
}

// RecordRaw 记录原始数据
type RecordRaw struct {
	Data Bytes
}

var _ encoding.BinaryUnmarshaler = (*RecordRaw)(nil)

// UnmarshalBinary 从二进制反序列化
func (r *RecordRaw) UnmarshalBinary(data []byte) error {
	r.Data = make([]byte, len(data))
	copy(r.Data, data)
	return nil
}
