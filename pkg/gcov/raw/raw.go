package raw

import (
	"encoding/binary"
	"fmt"
)

// Raw gcov 原始数据
type Raw struct {
	// 魔术码
	Magic Magic
	// 版本号
	Version Version
}

// UnmarshalBinary 从二进制反序列化
func (raw *Raw) UnmarshalBinary(data []byte) error {
	// magic 和 version
	if len(data) < 8 {
		return fmt.Errorf("data too short: %d, it should be at least 8 bytes for magic and version", len(data))
	}
	raw.Magic = Magic(binary.LittleEndian.Uint32(data[:4]))
	raw.Version = Version(binary.LittleEndian.Uint32(data[4:8]))

	switch raw.Magic {
	case MagicNote, MagicData:
	default:
		return fmt.Errorf("unknown magic: %s", raw.Magic)
	}

	// TODO: ...

	return nil
}
