package dumpraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	gcovraw "github.com/yhlooo/gcovgo/pkg/gcov/raw"
	"github.com/yhlooo/gcovgo/samples/gcovdata"
)

// TestDumpRaw 测试导出原始数据
func TestDumpRaw(t *testing.T) {
	for _, item := range gcovdata.Data {
		t.Run(item.NoteFile, testDumpSingleRaw(item.NoteFile))
		if item.DataFile != "" {
			t.Run(item.DataFile, testDumpSingleRaw(item.DataFile))
		}
	}
}

// testDumpSingleRaw 测试导出单个原始数据的方法
func testDumpSingleRaw(path string) func(t *testing.T) {
	return func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		content, err := gcovdata.FS.ReadFile(path)
		r.NoError(err)

		raw := gcovraw.Raw{}
		r.NoError(raw.UnmarshalBinary(content))

		a.NotZero(raw.Version)
		a.NotZero(raw.Magic)
		a.NotZero(raw.Stamp)
		a.NotEmpty(raw.Records)

		// TODO: 需要有更具体的校验
	}
}
