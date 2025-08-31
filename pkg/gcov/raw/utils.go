package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"strings"
)

// ParseString 解析字符串，返回解析的字符串、占 data 的字节数、解析错误
//
// 参考：
//
// Strings are padded with 1 to 4 NUL bytes, to bring the length up to a multiple
// of 4. The number of 4 bytes is stored, followed by the padded string.
//
// string: int32:0 | int32:length char* char:0 padding
// padding: | char:0 | char:0 char:0 | char:0 char:0 char:0
func ParseString(data []byte) (string, int, error) {
	if len(data) < 4 {
		return "", 0, newDataTooShortError(len(data), 4, "length")
	}

	length := binary.LittleEndian.Uint32(data[:4])
	data = data[4:]

	if len(data) < int(length)*4 {
		return "", 0, newDataTooShortError(len(data), int(length)*4, "content")
	}

	return strings.TrimRight(string(data[:length*4]), "\x00"), int(length)*4 + 4, nil
}

// Bytes 原始字节
type Bytes []byte

var _ fmt.Stringer = Bytes(nil)
var _ encoding.TextMarshaler = Bytes(nil)

// String 返回字符串表示
func (bytes Bytes) String() string {
	ret := ""
	remaining := bytes
	for len(remaining) >= 4 {
		ret += fmt.Sprintf("0x%08x ", binary.LittleEndian.Uint32(remaining[:4]))
		remaining = remaining[4:]
	}

	return strings.TrimRight(ret, " ")
}

// MarshalText 序列化为文本
func (bytes Bytes) MarshalText() ([]byte, error) {
	return []byte(bytes.String()), nil
}

// HexUint32 以十六进制表示的 uint32
type HexUint32 uint32

var _ fmt.Stringer = HexUint32(0)
var _ encoding.TextMarshaler = HexUint32(0)

// String 返回字符串表示
func (v HexUint32) String() string {
	return fmt.Sprintf("0x%08x", uint32(v))
}

// MarshalText 序列化为文本
func (v HexUint32) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}
