package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// RecordFunction 函数记录
type RecordFunction struct {
	// 标识
	Ident uint32

	LineNoChecksum HexUint32
	CfgChecksum    HexUint32

	// 函数名
	Name string `json:",omitempty"`
	// 函数所在文件名
	Source string `json:",omitempty"`
	// 函数起始行号
	StartLineNo uint32 `json:",omitempty"`
	// 函数起始列
	StartColumn uint32 `json:",omitempty"`
	// 函数结束行号
	EndLineNo uint32 `json:",omitempty"`
}

var _ encoding.BinaryUnmarshaler = (*RecordFunction)(nil)

// UnmarshalBinary 从二进制反序列化
//
// note:
//
//	announce_function: header int32:ident int32:lineno_checksum
//	    int32:cfg_checksum string:name string:source
//	    int32:start_lineno int32:start_column int32:end_lineno
//
// data:
//
//	announce_function: header int32:ident int32:lineno_checksum int32:cfg_checksum
func (r *RecordFunction) UnmarshalBinary(data []byte) error {
	if len(data) < 12 {
		return newDataTooShortError(len(data), 12, "ident, lineno_checksum and cfg_checksum")
	}
	r.Ident = binary.LittleEndian.Uint32(data[:4])
	r.LineNoChecksum = HexUint32(binary.LittleEndian.Uint32(data[4:8]))
	r.CfgChecksum = HexUint32(binary.LittleEndian.Uint32(data[8:12]))
	data = data[12:]

	if len(data) == 0 {
		// gcda 中没有下面其它字段
		return nil
	}

	var err error
	var n int
	r.Name, n, err = ParseString(data)
	if err != nil {
		return fmt.Errorf("parse name error: %w", err)
	}
	data = data[n:]

	// 文档中没有提及，但 name 和 source 之间可能有个 \x00000000
	if len(data) > 4 && string(data[:4]) == "\x00\x00\x00\x00" {
		data = data[4:]
	}

	r.Source, n, err = ParseString(data)
	if err != nil {
		return fmt.Errorf("parse source error: %w", err)
	}
	data = data[n:]

	if len(data) < 12 {
		return newDataTooShortError(len(data), 12, "start_lineno, start_column and end_lineno")
	}
	r.StartLineNo = binary.LittleEndian.Uint32(data[:4])
	r.StartColumn = binary.LittleEndian.Uint32(data[4:8])
	r.EndLineNo = binary.LittleEndian.Uint32(data[8:12])

	return nil
}
