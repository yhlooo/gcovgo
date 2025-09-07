package gcov

import (
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

// IntermediateText 输出中间文本形式
func (v *Version) IntermediateText() string {
	return fmt.Sprintf("version:%d.%d.0\n", v.Major, v.Minor)
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
	Name        string `json:"name"`
	StartLine   uint32 `json:"start_line"`
	StartColumn uint32 `json:"start_column,omitempty"`
	EndLine     uint32 `json:"end_line"`
	EndColumn   uint32 `json:"end_column,omitempty"`

	ExecutionCount uint64 `json:"execution_count"`
	Blocks         uint32 `json:"blocks"`
	BlocksExecuted uint32 `json:"blocks_executed"`

	DemangledName string `json:"demangled_name,omitempty"`
}

// IntermediateText 输出中间文本形式
func (fn *Function) IntermediateText() string {
	return fmt.Sprintf("function:%d,%d,%d,%s\n", fn.StartLine, fn.EndLine, fn.ExecutionCount, fn.Name)
}

// Line 覆盖情况信息
type Line struct {
	LineNumber      uint32   `json:"line_number"`
	Count           uint64   `json:"count"`
	Branches        []Branch `json:"branches,omitempty"`
	UnexecutedBlock bool     `json:"unexecuted,omitempty"`
	FunctionName    string   `json:"function_name,omitempty"`
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
	Count       uint64 `json:"count"`
	Fallthrough bool   `json:"fallthrough,omitempty"`
	Throw       bool   `json:"throw,omitempty"`
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
