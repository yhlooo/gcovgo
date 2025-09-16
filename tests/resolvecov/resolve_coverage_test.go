package resolvecov

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yhlooo/gcovgo/pkg/gcov"
	"github.com/yhlooo/gcovgo/samples/gcovdata"
)

// TestResolveCoverage 测试解算覆盖率
func TestResolveCoverage(t *testing.T) {
	for _, item := range gcovdata.Data {
		t.Run(item.NoteFile, testResolveSingleCoverage(item))
	}
}

// testResolveSingleCoverage 测试解算单个 note 和 data 文件覆盖率
func testResolveSingleCoverage(item gcovdata.Item) func(t *testing.T) {
	return func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		// 打开测试文件
		noteFile, err := gcovdata.FS.Open(item.NoteFile)
		r.NoError(err)
		defer func() { _ = noteFile }()
		var dataFile io.ReadCloser
		if item.DataFile != "" {
			dataFile, err = gcovdata.FS.Open(item.DataFile)
			r.NoError(err)
			defer func() { _ = dataFile.Close() }()
		}

		// 测试
		info, err := gcov.ResolveBinary(noteFile, dataFile)
		r.NoError(err)

		// 校验结果
		// TODO: 因暂未支持读取代码内容，暂不校验人类可读格式的输出
		//if item.HumanReadableOutputFile != "" {
		//	expected, err := gcovdata.FS.ReadFile(item.HumanReadableOutputFile)
		//	r.NoError(err)
		//	a.Equal(string(expected), info.HumanReadableText(t.Context()))
		//}
		if item.IntermediateOutputFile != "" {
			expected, err := gcovdata.FS.ReadFile(item.IntermediateOutputFile)
			r.NoError(err)

			if item.IntermediaJSON {
				infoJSON, err := json.Marshal(info)
				r.NoError(err)
				a.JSONEq(string(expected), string(infoJSON))
			} else {
				a.Equal(string(expected), info.IntermediateText(t.Context()))
			}
		}
	}
}
