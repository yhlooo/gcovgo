package gcov

import (
	"encoding"
	"fmt"
)

// CoverageInfo 覆盖情况信息
type CoverageInfo struct {
	// gcc 版本
	GCCVersion Version `json:"gcc_version"`
	// 格式版本
	FormatVersion string `json:"format_version,omitempty"`
	// GCDA 文件名
	DataFile string `json:"data_file,omitempty"`
	// 执行解析的工作目录
	CurrenWorkingDirectory string `json:"current_working_directory,omitempty"`
	// 文件覆盖情况
	Files []File `json:"files"`
}

// IntermediateText 输出中间文本形式
func (info *CoverageInfo) IntermediateText() string {
	ret := info.GCCVersion.IntermediateText()
	for _, file := range info.Files {
		ret += file.IntermediateText()
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
func (v *Version) IntermediateText() string {
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
func (f *File) IntermediateText() string {
	ret := fmt.Sprintf("file:%s\n", f.Filename)
	for _, fn := range f.Functions {
		ret += fn.IntermediateText()
	}
	for _, ln := range f.Lines {
		ret += ln.IntermediateText()
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
}

// IntermediateText 输出中间文本形式
func (fn *Function) IntermediateText() string {
	return fmt.Sprintf("function:%d,%d,%d,%s\n", fn.StartLine, fn.EndLine, fn.ExecutionCount, fn.Name)
}

// Line 覆盖情况信息
type Line struct {
	// 行号
	LineNumber uint32 `json:"line_number"`
	// 执行次数
	Count uint64 `json:"count"`
	// 分支
	Branches []Branch `json:"branches"`
	// 该行是否包含未执行的块
	UnexecutedBlock bool `json:"unexecuted_block"`
	// 函数名
	FunctionName string `json:"function_name"`
}

// IntermediateText 输出中间文本形式
func (ln *Line) IntermediateText() string {
	unexecutedBlock := "0"
	if ln.UnexecutedBlock {
		unexecutedBlock = "1"
	}
	ret := fmt.Sprintf("lcount:%d,%d,%s\n", ln.LineNumber, ln.Count, unexecutedBlock)
	for _, br := range ln.Branches {
		ret += br.IntermediateText(ln.LineNumber)
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
func (br *Branch) IntermediateText(lineNo uint32) string {
	coverageType := "nottaken"
	if br.Count > 0 {
		coverageType = "taken"
	}
	// TODO: notexec
	return fmt.Sprintf("branch:%d,%s\n", lineNo, coverageType)
}
