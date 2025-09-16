package gcov

import (
	"context"
	"encoding"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-logr/logr"
)

// CoverageInfo 覆盖情况信息
type CoverageInfo struct {
	// gcc 版本
	GCCVersion Version `json:"gcc_version"`
	// 格式版本
	FormatVersion string `json:"format_version,omitempty"`
	// 处理的数据文件名
	DataFile string `json:"data_file,omitempty"`
	// gcov note 文件名
	GcovNoteFile string `json:"-"`
	// gcov data 文件名
	GcovDataFile string `json:"-"`
	// 执行解析的工作目录
	CurrenWorkingDirectory string `json:"current_working_directory,omitempty"`
	// 文件覆盖情况
	Files []File `json:"files"`
}

// IntermediateText 输出中间文本形式
func (info *CoverageInfo) IntermediateText(ctx context.Context) string {
	ret := info.GCCVersion.IntermediateText(ctx)
	for _, file := range info.Files {
		ret += file.IntermediateText(ctx)
	}
	return ret
}

// HumanReadableText 输出人类可读的文本形式
func (info *CoverageInfo) HumanReadableText(ctx context.Context) string {
	logger := logr.FromContextOrDiscard(ctx)

	ret := ""
	for _, file := range info.Files {
		fileContent, err := os.ReadFile(file.Filename)
		if err != nil {
			logger.Info(fmt.Sprintf("WARN: read file %q error: %v", file.Filename, err))
		}
		ret += fmt.Sprintf(`        -:    0:Source:%s
        -:    0:Graph:%s
        -:    0:Data:%s
        -:    0:Runs:1
        -:    0:Programs:1
`, file.Filename, info.GcovNoteFile, info.GcovDataFile)
		ret += file.HumanReadableText(ctx, fileContent)
	}
	return ret
}

// Version 版本
type Version struct {
	// 主版本号
	Major int `json:"major"`
	// 子版本号
	Minor int `json:"minor"`
	// 发行状态
	// 比如 'e' experimental ， 'p' prerelease ， 'r' release
	Status string `json:"status,omitempty"`
}

var _ fmt.Stringer = (*Version)(nil)
var _ encoding.TextMarshaler = (*Version)(nil)

// IntermediateText 输出中间文本形式
func (v *Version) IntermediateText(_ context.Context) string {
	return fmt.Sprintf("version:%s\n", v.String())
}

// String 返回字符串表示
func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.0", v.Major, v.Minor)
}

// MarshalText 序列化为文本
func (v *Version) MarshalText() ([]byte, error) {
	return []byte(v.String()), nil
}

// File 文件覆盖情况信息
type File struct {
	// 文件名
	Filename string `json:"file"`
	// 文件中函数覆盖情况
	Functions []Function `json:"functions,omitempty"`
	// 文件中行覆盖情况
	Lines []Line `json:"lines,omitempty"`
}

// IntermediateText 输出中间文本形式
func (f *File) IntermediateText(ctx context.Context) string {
	ret := fmt.Sprintf("file:%s\n", f.Filename)
	for _, fn := range f.Functions {
		ret += fn.IntermediateText(ctx)
	}
	for _, ln := range f.Lines {
		ret += ln.IntermediateText(ctx)
	}
	return ret
}

// HumanReadableText 输出人类可读的文本形式
func (f *File) HumanReadableText(ctx context.Context, content []byte) string {
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	linesN := len(lines)
	if len(f.Lines) > 0 && int(f.Lines[len(f.Lines)-1].LineNumber) > linesN {
		linesN = int(f.Lines[len(f.Lines)-1].LineNumber)
	}

	fnI := 0
	lnI := 0

	ret := ""
	for i := 0; i < linesN; i++ {
		// 获取行内容
		lnContent := "/*EOF*/"
		if i < len(lines) {
			lnContent = lines[i]
		}

		// 获取行执行次数和分支信息
		count := "-"
		var brs []Branch
		var calls []Branch
		if lnI < len(f.Lines) {
			ln := f.Lines[lnI]
			if ln.LineNumber == uint32(i+1) {
				count = strconv.FormatUint(ln.Count, 10)
				brs = ln.Branches
				calls = ln.CallBranches
				lnI++
			}
		}

		// 获取分支信息
		if fnI < len(f.Functions) {
			fn := f.Functions[fnI]
			if fn.StartLine == uint32(i+1) {
				ret += fn.HumanReadableText(ctx)
				fnI++
			}
		}

		ret += fmt.Sprintf(" %8s: %4d:%s\n", count, i+1, lnContent)
		for j, br := range brs {
			ret += br.HumanReadableText(ctx, j, false)
		}
		for j, br := range calls {
			ret += br.HumanReadableText(ctx, j, true)
		}
	}

	return ret
}

// Function 函数覆盖情况信息
type Function struct {
	// 函数名
	Name string `json:"name"`
	// 去混淆的函数名
	DemangledName string `json:"demangled_name,omitempty"`

	// 起始行号
	StartLine uint32 `json:"start_line"`
	// 起始列号
	StartColumn uint32 `json:"start_column,omitempty"`
	// 结束行号
	EndLine uint32 `json:"end_line"`
	// 结束列号
	EndColumn uint32 `json:"end_column,omitempty"`

	// 函数中基本块数目
	Blocks uint32 `json:"blocks"`
	// 函数中有执行的基本块数目
	BlocksExecuted uint32 `json:"blocks_executed"`
	// 函数执行次数
	ExecutionCount uint64 `json:"execution_count"`
	// 函数返回次数
	ReturnCount uint64 `json:"-"`
}

// IntermediateText 输出中间文本形式
func (fn *Function) IntermediateText(_ context.Context) string {
	return fmt.Sprintf("function:%d,%d,%d,%s\n", fn.StartLine, fn.EndLine, fn.ExecutionCount, fn.Name)
}

// HumanReadableText 输出人类可读的文本形式
func (fn *Function) HumanReadableText(_ context.Context) string {
	returned := uint64(0)
	if fn.ExecutionCount != 0 {
		returned = fn.ReturnCount * 100 / fn.ExecutionCount
	}
	executed := uint32(0)
	if fn.Blocks != 0 {
		executed = fn.BlocksExecuted * 100 / fn.Blocks
	}
	return fmt.Sprintf(
		"function %s called %d returned %d%% blocks executed %d%%\n",
		fn.Name, fn.ExecutionCount, returned, executed,
	)
}

// Line 覆盖情况信息
type Line struct {
	// 行号
	LineNumber uint32 `json:"line_number"`
	// 执行次数
	Count uint64 `json:"count"`
	// 分支
	Branches []Branch `json:"branches"`
	// 调用其它函数分支
	CallBranches []Branch `json:"-"`
	// 该行是否包含未执行的块
	UnexecutedBlock bool `json:"unexecuted_block"`
	// 函数名
	FunctionName string `json:"function_name"`
}

// IntermediateText 输出中间文本形式
func (ln *Line) IntermediateText(ctx context.Context) string {
	unexecutedBlock := "0"
	if ln.UnexecutedBlock {
		unexecutedBlock = "1"
	}
	ret := fmt.Sprintf("lcount:%d,%d,%s\n", ln.LineNumber, ln.Count, unexecutedBlock)
	for _, br := range ln.Branches {
		ret += br.IntermediateText(ctx, ln.LineNumber)
	}
	return ret
}

// Branch 分支覆盖情况信息
type Branch struct {
	// 执行次数
	Count uint64 `json:"count"`
	// 是否直落分支
	Fallthrough bool `json:"fallthrough"`
	// 是否异常分支
	Throw bool `json:"throw"`

	blockNo uint32
}

// IntermediateText 输出中间文本形式
func (br *Branch) IntermediateText(_ context.Context, lineNo uint32) string {
	coverageType := "nottaken"
	if br.Count > 0 {
		coverageType = "taken"
	}
	// TODO: notexec
	return fmt.Sprintf("branch:%d,%s\n", lineNo, coverageType)
}

// HumanReadableText 输出人类可读的文本形式
func (br *Branch) HumanReadableText(_ context.Context, i int, call bool) string {
	if call {
		return fmt.Sprintf("call %4d returned %d\n", i, br.Count)
	}
	suffix := ""
	if br.Fallthrough {
		suffix = " (fallthrough)"
	}
	return fmt.Sprintf("branch %2d taken %d%s\n", i, br.Count, suffix)
}
