package raw

// Note gcov note 原始数据（ gcno ）
type Note interface {
	// IsNote 判断是否 note （ gcno ）
	IsNote() bool
	// FunctionNotes 按 note 结构整理返回函数相关记录
	FunctionNotes() []FunctionNoteRecords
}

// FunctionNoteRecords note 中函数相关记录
type FunctionNoteRecords struct {
	Function *RecordFunction
	Blocks   *RecordBlocks
	Arcs     []*RecordArcs
	Lines    []*RecordLines
}

var _ Note = (*Raw)(nil)

// IsNote 判断是否 note （ gcno ）
func (raw *Raw) IsNote() bool {
	return raw.Magic == MagicNote
}

// FunctionNotes 按 note 结构整理返回函数相关记录
func (raw *Raw) FunctionNotes() []FunctionNoteRecords {
	var functions []FunctionNoteRecords

	var (
		funcRecord *RecordFunction
		blocks     *RecordBlocks
		arcs       []*RecordArcs
		lines      []*RecordLines
	)
	for _, record := range raw.Records {
		switch record.Tag {
		case TagFunction:
			if funcRecord != nil {
				functions = append(functions, FunctionNoteRecords{
					Function: funcRecord,
					Blocks:   blocks,
					Arcs:     arcs,
					Lines:    lines,
				})
			}
			funcRecord = record.Function
			blocks = nil
			arcs = nil
			lines = nil
		case TagBlocks:
			blocks = record.Blocks
		case TagArcs:
			arcs = append(arcs, record.Arcs)
		case TagLines:
			lines = append(lines, record.Lines)
		}
	}
	if funcRecord != nil {
		functions = append(functions, FunctionNoteRecords{
			Function: funcRecord,
			Blocks:   blocks,
			Arcs:     arcs,
			Lines:    lines,
		})
	}

	return functions
}
