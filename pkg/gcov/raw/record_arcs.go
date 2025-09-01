package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"strings"
)

// RecordArcs 边记录
type RecordArcs struct {
	// 块编号
	BlockNo uint32
	// 块往外连接的边
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

// Arc 边
type Arc struct {
	// 目标块编号
	DestBlock uint32
	// 边属性
	Flags ArcFlag
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
	arc.Flags = ArcFlag(binary.LittleEndian.Uint32(data[4:8]))

	return nil
}

// ArcFlag 边属性
type ArcFlag uint32

var _ fmt.Stringer = ArcFlag(0)

const (
	ArcFlagOnTree      ArcFlag = 1
	ArcFlagFake        ArcFlag = 1 << 1
	ArcFlagFallthrough ArcFlag = 1 << 2
)

// OnTree 是否包含 OnTree flag
func (flag ArcFlag) OnTree() bool {
	return flag&ArcFlagOnTree != 0
}

// Fake 是否包含 Fake flag
func (flag ArcFlag) Fake() bool {
	return flag&ArcFlagFake != 0
}

// Fallthrough 是否包含 Fallthrough flag
func (flag ArcFlag) Fallthrough() bool {
	return flag&ArcFlagFallthrough != 0
}

// String 返回字符串表示
func (flag ArcFlag) String() string {
	ret := ""
	remaining := flag
	if remaining.OnTree() {
		ret += "OnTree|"
		remaining &= ^ArcFlagOnTree
	}
	if remaining.Fake() {
		ret += "Fake|"
		remaining &= ^ArcFlagFake
	}
	if remaining.Fallthrough() {
		ret += "Fallthrough|"
		remaining &= ^ArcFlagFallthrough
	}
	if remaining != 0 {
		ret += fmt.Sprintf("%032b", remaining)
	}
	return strings.TrimRight(ret, "|")
}
