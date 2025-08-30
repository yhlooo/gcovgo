package raw

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMagic_String 测试 Magic.String 方法
func TestMagic_String(t *testing.T) {
	a := assert.New(t)

	a.Equal("gcno", MagicNote.String())
	a.Equal("gcda", MagicData.String())
	a.Equal("abcd(0x61626364)", Magic(binary.LittleEndian.Uint32([]byte("dcba"))).String())
}
