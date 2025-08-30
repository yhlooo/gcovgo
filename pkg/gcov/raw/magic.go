package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// Magic 文件魔术码，用于标识文件类型
type Magic uint32

var _ fmt.Stringer = Magic(0)
var _ encoding.TextMarshaler = Magic(0)

// String 返回字符串表示
func (magic Magic) String() string {
	switch magic {
	case MagicNote:
		return "gcno"
	case MagicData:
		return "gcda"
	}
	return fmt.Sprintf("%s(0x%08x)", string(binary.BigEndian.AppendUint32(nil, uint32(magic))), int32(magic))
}

// MarshalText 序列化为文本
func (magic Magic) MarshalText() ([]byte, error) {
	return []byte(magic.String()), nil
}

const (
	// MagicNote note 文件的 magic
	MagicNote Magic = 'g'<<24 | 'c'<<16 | 'n'<<8 | 'o'
	// MagicData data 文件的 magic
	MagicData Magic = 'g'<<24 | 'c'<<16 | 'd'<<8 | 'a'
)
