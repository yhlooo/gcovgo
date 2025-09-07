package raw

// Data gcov data 原始数据（ gcda ）
type Data interface {
	// IsData 判断是否 data （ gcda ）
	IsData() bool
	// FunctionsData 按 data 结构整理返回函数相关记录
	FunctionsData() []FunctionDataRecords
}

// FunctionDataRecords data 中函数相关记录
type FunctionDataRecords struct {
	Function *RecordFunction
	Counter  *RecordCounter
}

var _ Data = (*Raw)(nil)

// IsData 判断是否 data （ gcda ）
func (raw *Raw) IsData() bool {
	return raw.Magic == MagicData
}

// FunctionsData 按 data 结构整理返回函数相关记录
func (raw *Raw) FunctionsData() []FunctionDataRecords {
	var functions []FunctionDataRecords

	var (
		funcRecord *RecordFunction
		counter    *RecordCounter
	)
	for _, record := range raw.Records {
		switch record.Tag {
		case TagFunction:
			if funcRecord != nil {
				functions = append(functions, FunctionDataRecords{
					Function: funcRecord,
					Counter:  counter,
				})
			}
			funcRecord = record.Function
			counter = nil
		case TagCounter:
			counter = record.Counter
		}
	}
	if funcRecord != nil {
		functions = append(functions, FunctionDataRecords{
			Function: funcRecord,
			Counter:  counter,
		})
	}

	return functions
}

// FunctionCounters 获取每个函数计数器，键为函数 Ident ，值为计数器
func (raw *Raw) FunctionCounters() map[uint32][]uint64 {
	functions := raw.FunctionsData()
	counters := make(map[uint32][]uint64, len(functions))
	for _, fn := range functions {
		if fn.Function == nil || fn.Counter == nil {
			continue
		}
		counters[fn.Function.Ident] = fn.Counter.Counts
	}
	return counters
}
