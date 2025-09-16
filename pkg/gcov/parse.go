package gcov

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yhlooo/gcovgo/pkg/gcov/cfg"
	"github.com/yhlooo/gcovgo/pkg/gcov/raw"
)

// ResolveBinaryFile 解析 gcov 二进制文件
func ResolveBinaryFile(noteFileName, dataFileName string) (*CoverageInfo, error) {
	noteFile, err := os.Open(noteFileName)
	if err != nil {
		return nil, fmt.Errorf("open note file %q error: %w", noteFileName, err)
	}
	defer func() { _ = noteFile.Close() }()

	var dataReader io.Reader
	if dataFileName != "" {
		dataFile, err := os.Open(dataFileName)
		if err != nil {
			return nil, fmt.Errorf("open data file %q error: %w", dataFileName, err)
		}
		defer func() { _ = dataFile.Close() }()
		dataReader = dataFile
	}

	return ResolveBinary(noteFile, dataReader)
}

// ResolveBinary 解析 gcov 二进制
func ResolveBinary(note, data io.Reader) (*CoverageInfo, error) {
	// 反序列化 note
	noteContent, err := io.ReadAll(note)
	if err != nil {
		return nil, fmt.Errorf("read note error: %w", err)
	}
	noteObj := &raw.Raw{}
	if err := noteObj.UnmarshalBinary(noteContent); err != nil {
		return nil, fmt.Errorf("unmarshal note error: %w", err)
	}
	if !noteObj.IsNote() {
		return nil, fmt.Errorf("not a valid note magic: %q", noteObj.Magic.String())
	}

	// 反序列化 data ，获取计数器
	var counters map[uint32][]uint64
	if data != nil {
		dataContent, err := io.ReadAll(data)
		if err != nil {
			return nil, fmt.Errorf("read data error: %w", err)
		}
		dataObj := &raw.Raw{}
		if err := dataObj.UnmarshalBinary(dataContent); err != nil {
			return nil, fmt.Errorf("unmarshal data error: %w", err)
		}
		if !dataObj.IsData() {
			return nil, fmt.Errorf("not a valid data magic: %q", dataObj.Magic.String())
		}
		counters = dataObj.FunctionCounters()
	}

	major, minor, status := noteObj.Version.Parse()
	ret := &CoverageInfo{
		GCCVersion: Version{
			Major:  major,
			Minor:  minor,
			Status: status,
		},
		FormatVersion:          "1",
		CurrenWorkingDirectory: noteObj.CurrenWorkingDirectory,
	}
	filesMap := map[string]*File{}
	functions := noteObj.FunctionNotes()
	for _, fn := range functions {
		if fn.Function == nil {
			continue
		}

		// 计算函数控制流图
		blocks := len(fn.Arcs) + 1
		if major >= 8 && fn.Blocks != nil && len(fn.Blocks.Flags) > 0 {
			blocks = int(fn.Blocks.Flags[0])
		}
		if blocks <= 0 {
			continue
		}
		graph, err := cfg.BuildCFG(blocks, fn.Arcs, counters[fn.Function.Ident])
		if err != nil {
			return ret, fmt.Errorf(
				"build function %q (file: %q) control flow graph error: %w",
				fn.Function.Name, fn.Function.Source, err,
			)
		}

		execBlocks := uint32(0)
		for i := uint32(2); i < uint32(blocks); i++ {
			if blk := graph.Get(i); blk != nil && blk.Count() > 0 {
				execBlocks++
			}
		}

		// 记录函数覆盖信息
		fileName := fn.Function.Source
		f := filesMap[fileName]
		if f == nil {
			ret.Files = append(ret.Files, File{Filename: fileName})
			f = &ret.Files[len(ret.Files)-1]
			filesMap[fileName] = f
		}

		f.Functions = append(f.Functions, Function{
			Name:           fn.Function.Name,
			StartLine:      fn.Function.StartLineNo,
			StartColumn:    fn.Function.StartColumn,
			EndLine:        fn.Function.EndLineNo,
			EndColumn:      fn.Function.EndColumn,
			ExecutionCount: graph.Get(0).Count(),
			ReturnCount:    graph.Get(1).Count(),
			Blocks:         uint32(blocks) - 2,
			BlocksExecuted: execBlocks,
			DemangledName:  fn.Function.Name, // TODO: 应该不总是与 Name 相同，具体取值来源不确定
		})

		// 记录行覆盖信息

		for _, blkLines := range fn.Lines {
			blk := graph.Get(blkLines.BlockNo)
			if blk == nil {
				continue
			}
			for i, item := range blkLines.Lines {
				if item.Filename != "" {
					if item.Filename == fileName {
						continue
					}
					// 切换文件
					fileName = item.Filename
					f = filesMap[fileName]
					if f == nil {
						ret.Files = append(ret.Files, File{Filename: fileName})
						f = &ret.Files[len(ret.Files)-1]
						filesMap[fileName] = f
					}
					continue
				}

				// 分支
				call := false
				branches := make([]Branch, 0)
				if blkOut := blk.Out(); i == len(blkLines.Lines)-1 && len(blkOut) > 1 {
					// 出边对应分支关联到块中最后一行
					for _, arc := range blkOut {
						dstBlkNo := arc.Destination().No()
						if dstBlkNo == 1 {
							call = true
							continue
						}
						branches = append(branches, Branch{
							Count:       arc.Count(),
							Fallthrough: arc.Flags()&raw.ArcFlagFallthrough != 0,
							Throw:       false, // TODO: ...
							blockNo:     dstBlkNo,
						})
					}
				}
				sort.Slice(branches, func(i, j int) bool {
					return branches[i].blockNo < branches[j].blockNo
				})
				var callBranches []Branch
				if call {
					callBranches = branches
					branches = make([]Branch, 0)
				}

				// 行
				f.Lines = append(f.Lines, Line{
					LineNumber:      item.LineNo,
					Count:           blk.Count(),
					Branches:        branches,
					CallBranches:    callBranches,
					UnexecutedBlock: blk.Count() == 0,
					FunctionName:    fn.Function.Name,
				})
			}
		}
	}

	// 排序、去重
	for fileI := range ret.Files {
		sort.Slice(ret.Files[fileI].Functions, func(i, j int) bool {
			return ret.Files[fileI].Functions[i].StartLine < ret.Files[fileI].Functions[j].StartLine
		})
		sort.SliceStable(ret.Files[fileI].Lines, func(i, j int) bool {
			return ret.Files[fileI].Lines[i].LineNumber < ret.Files[fileI].Lines[j].LineNumber
		})
		// 合并相同行
		newLines := make([]Line, 0, len(ret.Files[fileI].Lines))
		lastLineNo := uint32(0)
		for _, line := range ret.Files[fileI].Lines {
			if lastLineNo == 0 || line.LineNumber != lastLineNo {
				// 不同行跳过
				newLines = append(newLines, line)
				lastLineNo = line.LineNumber
				continue
			}

			// 相同行合并
			lastLine := &newLines[len(newLines)-1]
			if line.Count > lastLine.Count {
				lastLine.Count = line.Count
			}
			lastLine.UnexecutedBlock = lastLine.UnexecutedBlock || line.UnexecutedBlock
			lastLine.Branches = append(lastLine.Branches, line.Branches...)
			lastLine.CallBranches = append(lastLine.CallBranches, line.CallBranches...)
		}
		ret.Files[fileI].Lines = newLines
	}

	return ret, nil
}
