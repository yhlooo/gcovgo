package raw

import (
	"encoding"
	"encoding/binary"
	"fmt"
)

// RecordProgramSummary 程序摘要记录
type RecordProgramSummary struct {
	// 校验和
	Checksum HexUint32
	// 计数摘要
	CountSummaries []CountSummary
}

var _ encoding.BinaryUnmarshaler = (*RecordProgramSummary)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	summary: int32:checksum {count-summary}GCOV_COUNTERS_SUMMABLE
func (r *RecordProgramSummary) UnmarshalBinary(data []byte) error {
	if len(data) < 4 {
		return newDataTooShortError(len(data), 4, "checksum")
	}
	r.Checksum = HexUint32(binary.LittleEndian.Uint32(data[:4]))
	data = data[4:]

	for len(data) > 4 {
		summary := CountSummary{}
		if err := summary.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unmarshal count summary %d error: %w", len(r.CountSummaries), err)
		}
		r.CountSummaries = append(r.CountSummaries, summary)
		data = data[summary.Size():]
	}
	return nil
}

// CountSummary 计数摘要
type CountSummary struct {
	Num       uint32
	Runs      uint32
	Sum       uint64
	Max       uint64
	SumMax    uint64
	Histogram Histogram
}

var _ encoding.BinaryUnmarshaler = (*CountSummary)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	count-summary: int32:num int32:runs int64:sum int64:max int64:sum_max histogram
func (summary *CountSummary) UnmarshalBinary(data []byte) error {
	if len(data) < 24 {
		return newDataTooShortError(len(data), 32, "num, runs, sum and max")
	}
	summary.Num = binary.LittleEndian.Uint32(data[:4])
	summary.Runs = binary.LittleEndian.Uint32(data[4:8])
	summary.Sum = binary.LittleEndian.Uint64(data[8:16])
	summary.Max = binary.LittleEndian.Uint64(data[16:24])
	summary.SumMax = binary.LittleEndian.Uint64(data[24:32])
	data = data[32:]

	if err := summary.Histogram.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("unmarshal histogram error: %w", err)
	}

	return nil
}

// Size 返回数据大小
func (summary *CountSummary) Size() int {
	return 32 + summary.Histogram.Size()
}

// Histogram 直方图
type Histogram struct {
	BitVectors [8]HexUint32
	Buckets    []HistogramBucket

	bucketsSize int
}

var _ encoding.BinaryUnmarshaler = (*Histogram)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	{int32:bitvector}8 histogram-buckets*
func (h *Histogram) UnmarshalBinary(data []byte) error {
	if len(data) < 32 {
		return newDataTooShortError(len(data), 32, "bitvectors")
	}
	for i := range h.BitVectors {
		h.BitVectors[i] = HexUint32(binary.LittleEndian.Uint32(data[:4]))
		data = data[4:]
	}

	for len(data) >= bucketSize {
		bucket := HistogramBucket{}
		if err := bucket.UnmarshalBinary(data); err != nil {
			return fmt.Errorf("unmarshal histogram bucket %d error: %w", len(h.Buckets), err)
		}
		h.Buckets = append(h.Buckets, bucket)
		h.bucketsSize += bucketSize
		data = data[bucketSize:]
	}

	return nil
}

// Size 返回数据大小
func (h *Histogram) Size() int {
	return 32 + h.bucketsSize
}

const bucketSize = 20

// HistogramBucket 直方图桶
type HistogramBucket struct {
	Num uint32
	Min uint64
	Sum uint64
}

var _ encoding.BinaryUnmarshaler = (*HistogramBucket)(nil)

// UnmarshalBinary 从二进制反序列化
//
//	histogram-buckets: int32:num int64:min int64:sum
func (bucket *HistogramBucket) UnmarshalBinary(data []byte) error {
	if len(data) < 20 {
		return newDataTooShortError(len(data), 20, "num, min and sum")
	}
	bucket.Num = binary.LittleEndian.Uint32(data[:4])
	bucket.Min = binary.LittleEndian.Uint64(data[4:12])
	bucket.Sum = binary.LittleEndian.Uint64(data[12:20])
	return nil
}
