package raw

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestVersion_String 测试 Version.String 方法
func TestVersion_String(t *testing.T) {
	a := assert.New(t)

	a.Equal("B70e", Version(binary.LittleEndian.Uint32([]byte("e07B"))).String())
}
