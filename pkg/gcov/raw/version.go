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

// Parse 解析版本号
//
// 版本号是一个 uint32 ，也可以表示为 4 个字节
//
// 对于 gcc6 及以下版本，第一个字母表示 major ，接着两个字母 minor （若 minor 小于 10 则包含一个前导 0 ），
// 最后一个字母表示发行状态（比如 'e' experimental ， 'p' prerelease ， 'r' release ）。
// 比如 304e 表示 3.4 experimental
//
// 对于 gcc7 及以上版本，前两个字母表示 major ，以 A0 开始（以与之前版本区分），接着一个字母 minor ，
// 最后一个字母表示发行状态（ e p r 等）。比如 B70e 表示 17.0 experimental
func (v Version) Parse() (major, minor int, status string) {
	char0 := byte((v >> 24) & 255)
	char1 := byte((v >> 16) & 255)
	char2 := byte((v >> 8) & 255)
	status = string(byte(v & 255))

	if char0 < 'A' {
		major = int(char0 - '0')
		minor = int(char1-'0')*10 + int(char2) - '0'
	} else {
		major = int(char0-'A')*10 + int(char1-'0')
		minor = int(char2 - '0')
	}

	return
}

const (
	Version8  Version = 'A'<<24 | '8'<<16 | '0'<<8 | '*'
	Version9  Version = 'A'<<24 | '9'<<16 | '0'<<8 | '*'
	Version12 Version = 'B'<<24 | '2'<<16 | '0'<<8 | '*'
)
