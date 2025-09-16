package gcovdata

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"
)

func init() {
	if err := findAllData(); err != nil {
		panic(fmt.Errorf("find all data error: %w", err))
	}
}

//go:embed gcc_*
var FS embed.FS

// Data gcov 测试数据
var Data []Item

// Item gcov 测试数据项
type Item struct {
	// .gcno 文件路径
	NoteFile string
	// .gcda 文件路径
	DataFile string
	// 人类可读的 .gcov 文件路径
	HumanReadableOutputFile string
	// 文本或 JSON 中间格式文件路径
	IntermediateOutputFile string
	// IntermediateOutputFile 是否 JSON 中间格式
	IntermediaJSON bool
}

// findAllData 查找所有测试数据
func findAllData() error {
	items, err := FS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("list data root error: %w", err)
	}
	for i, item := range items {
		dataItems, err := walkDataDir(item.Name())
		if err != nil {
			return fmt.Errorf("walk data dir %d error: %w", i, err)
		}
		Data = append(Data, dataItems...)
	}
	return nil
}

// walkDataDir 查找一个数据目录
func walkDataDir(root string) ([]Item, error) {
	var ret []Item

	roots := []string{root}
	for len(roots) > 0 {
		r := roots[0]
		roots = roots[1:]

		items, err := FS.ReadDir(r)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			path := filepath.Join(r, item.Name())
			if path == "human_readable" || path == "intermediate" {
				continue
			}
			if item.IsDir() {
				roots = append(roots, path)
				continue
			}

			if filepath.Ext(item.Name()) != ".gcno" {
				continue
			}

			data := Item{
				NoteFile:                path,
				DataFile:                "",
				HumanReadableOutputFile: "",
				IntermediateOutputFile:  "",
				IntermediaJSON:          false,
			}

			gcdaPath := strings.TrimSuffix(path, ".gcno") + ".gcda"
			if f, err := FS.Open(gcdaPath); err == nil {
				_ = f.Close()
				data.DataFile = gcdaPath
			}

			humanReadablePath := filepath.Join(
				root, "human_readable",
				strings.TrimSuffix(item.Name(), ".gcno")+".gcov",
			)
			if f, err := FS.Open(humanReadablePath); err == nil {
				_ = f.Close()
				data.HumanReadableOutputFile = humanReadablePath
			}

			intermediatePaths := []string{
				filepath.Join(root, "intermediate", item.Name()+".gcov"),
				filepath.Join(root, "intermediate", strings.TrimSuffix(item.Name(), ".gcno")+".gcov"),
				filepath.Join(root, "intermediate", item.Name()+".gcov.json"),
				filepath.Join(root, "intermediate", strings.TrimSuffix(item.Name(), ".gcno")+".gcov.json"),
			}
			for _, p := range intermediatePaths {
				if f, err := FS.Open(p); err == nil {
					_ = f.Close()
					data.IntermediateOutputFile = p
					data.IntermediaJSON = filepath.Ext(p) == ".json"
					break
				}
			}

			ret = append(ret, data)
		}
	}

	return ret, nil
}
