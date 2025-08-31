package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// RecordLines 行记录
type RecordLines struct {
	// 块编号
	BlockNo uint32
	// 块中的行
	Lines []FileOrLine
}

var _ encoding.BinaryUnmarshaler = (*RecordLines)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	lines: header int32:block_no line* int32:0 string:NULL
func (r *RecordLines) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return newDataTooShortError(len(data), 4, "block_no")
	}
	r.BlockNo = binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	for {
		fl := FileOrLine{}
		if err := fl.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unmarshal line %d error: %w", len(r.Lines), err)
		}
		if fl.LineNo != 0 {
			data = data[4:]
		} else if fl.Filename != "" {
			fileNameLen := binary.LittleEndian.Uint32(data[4:8])
			data = data[8+int(fileNameLen)*4:]
		} else {
			// 解析完了
			return nil
		}
		r.Lines = append(r.Lines, fl)
	}
}

// FileOrLine 行
type FileOrLine struct {
	// 行号
	LineNo uint32 `json:",omitempty"`
	// 行所属文件
	Filename string `json:",omitempty"`
}

// UnmarshalBinary 从二进制反序列化
//
//	int32:line_no | int32:0 string:filename
func (fl *FileOrLine) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return newDataTooShortError(len(data), 4, "line_no")
	}
	fl.LineNo = binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	if fl.LineNo != 0 {
		// 行信息
		return nil
	}

	// 文件名信息
	var err error
	fl.Filename, _, err = ParseString(data)
	if err != nil {
		return fmt.Errorf("parse filename error: %w", err)
	}

	return nil
}
