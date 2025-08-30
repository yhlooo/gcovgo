package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// Version 版本
type Version uint32

var _ fmt.Stringer = Version(0)
var _ encoding.TextMarshaler = Version(0)

// String 返回字符串表示
func (v Version) String() string {
	return string(binary.BigEndian.AppendUint32(nil, uint32(v)))
}

// MarshalText 序列化为文本
func (v Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

const (
	Version8 Version = 'A'<<24 | '8'<<16 | '0'<<8 | '*'
)
