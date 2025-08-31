package raw

import "encoding/binary"

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

	return string(data[:length*4]), int(length)*4 + 4, nil
}
