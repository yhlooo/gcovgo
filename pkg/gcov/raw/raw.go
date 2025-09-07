package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// Raw gcov 原始数据
type Raw struct {
	// 魔术码
	Magic Magic
	// 版本号
	Version Version
	// 时间戳
	Stamp uint32
	// 校验和
	Checksum HexUint32
	// 当前工作目录
	CurrenWorkingDirectory string `json:",omitempty"`

	SupportUnexecutedBlocks uint32 // TODO: 应该是 bool 类型？

	// 记录
	Records []Record `json:",omitempty"`
}

var _ encoding.BinaryUnmarshaler = (*Raw)(nil)

// UnmarshalBinary 从二进制反序列化
func (raw *Raw) UnmarshalBinary(data []byte) error {
	// magic version 和 stamp
	if len(data) < 12 {
		return newDataTooShortError(len(data), 12, "magic, version and stamp")
	}
	raw.Magic = Magic(binary.LittleEndian.Uint32(data[:4]))
	raw.Version = Version(binary.LittleEndian.Uint32(data[4:8]))
	raw.Stamp = binary.LittleEndian.Uint32(data[8:12])
	data = data[12:]

	if raw.Version >= Version12 {
		if len(data) < 4 {
			return newDataTooShortError(len(data), 4, "checksum")
		}
		raw.Checksum = HexUint32(binary.LittleEndian.Uint32(data[:4]))
		data = data[4:]
	}

	// support_unexecuted_blocks
	switch raw.Magic {
	case MagicNote:
		if raw.Version >= Version9 {
			cwd, n, err := ParseString(data, raw.Version)
			if err != nil {
				return fmt.Errorf("parse cwd error: %w", err)
			}
			raw.CurrenWorkingDirectory = cwd
			data = data[n:]
		}
		if raw.Version >= Version8 {
			if len(data) < 4 {
				return newDataTooShortError(len(data), 4, "support_unexecuted_blocks")
			}
			raw.SupportUnexecutedBlocks = binary.LittleEndian.Uint32(data[:4])
			data = data[4:]
		}
	case MagicData:
	default:
		return fmt.Errorf("unknown magic: %s", raw.Magic)
	}

	// records
	for len(data) >= 8 {
		record := Record{version: raw.Version}
		if err := record.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unmarshal record %d error: %w", len(raw.Records), err)
		}
		raw.Records = append(raw.Records, record)
		data = data[record.Size():]
	}

	return nil
}

// newDataTooShortError 创建数据太短错误
func newDataTooShortError(n, atLeast int, fieldName string) error {
	return fmt.Errorf("remaining data too short: %d, it should be at least %d bytes for %s", n, atLeast, fieldName)
}
