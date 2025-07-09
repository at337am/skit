package renamer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

func (r *Runner) findFiles() ([]fileInfo, error) {
	var files []fileInfo

	// 使用 os.ReadDir 只读取目录的第一层条目，不进行递归
	entries, err := os.ReadDir(r.DirPath)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败 %s: %w", r.DirPath, err)
	}

	for _, d := range entries {
		if d.IsDir() {
			continue
		}

		name := d.Name()
		ext := filepath.Ext(name)

		// 如果指定了文件格式，则只处理匹配的文件
		if r.FileExt != "" && !strings.EqualFold(ext, r.FileExt) {
			continue
		}

		info, err := d.Info()
		if err != nil {
			warnColor.Fprintf(os.Stderr, "注意: 获取文件信息失败 %s: %v\n", name, err)
			continue
		}

		files = append(files, fileInfo{
			path:    filepath.Join(r.DirPath, name), // 手动拼接完整路径
			modTime: info.ModTime(),
			ext:     ext,
		})
	}

	return files, nil
}

func (r *Runner) renameFiles(files []fileInfo) ([]renameResult, error) {
	if len(files) == 0 {
		return nil, errors.New("没有找到匹配的文件")
	}

	// 依据文件修改时间进行排序
	slices.SortFunc(files, func(a, b fileInfo) int {
		if r.ReverseSort {
			return b.modTime.Compare(a.modTime) // 降序
		}
		return a.modTime.Compare(b.modTime) // 升序
	})

	// renameOp 单个文件的重命名计划
	type renameOp struct {
		originalPath string
		tmpPath      string
		finalPath    string
	}

	const tmpSuffix = ".tmrn-tmp"

	// 步骤 1: 负责根据文件列表生成重命名计划
	numFiles := len(files)
	digits := len(fmt.Sprintf("%d", numFiles))
	formatTemplate := fmt.Sprintf("%%0%dd", digits)
	plan := make([]renameOp, numFiles)

	for i, file := range files {
		finalName := fmt.Sprintf(formatTemplate, i+1) + file.ext
		plan[i] = renameOp{
			originalPath: file.path,
			tmpPath:      file.path + tmpSuffix,
			finalPath:    filepath.Join(r.DirPath, finalName),
		}
	}

	// 步骤 2: 执行第一阶段重命名 (原始文件 -> 临时文件)
	// 这一阶段是原子性的，如果中途失败，会尝试回滚所有已成功的操作。
	for i, op := range plan {
		if err := os.Rename(op.originalPath, op.tmpPath); err != nil {
			// 尝试回滚已成功的重命名操作
			for j := range i {
				// 尽力而为，忽略回滚错误
				_ = os.Rename(plan[j].tmpPath, plan[j].originalPath)
			}
			return nil, fmt.Errorf("操作已中断: 文件 '%s' 重命名失败: %w. 已尝试回滚", filepath.Base(op.originalPath), err)
		}
	}

	// 步骤 3: 执行第二阶段重命名 (临时文件 -> 最终文件), 并收集结果
	results := make([]renameResult, 0, len(plan))
	for _, op := range plan {
		if err := os.Rename(op.tmpPath, op.finalPath); err != nil {
			warnColor.Fprintf(os.Stderr, "注意: 无法将 %s 重命名为 %s: %v\n", filepath.Base(op.tmpPath), filepath.Base(op.finalPath), err)
			continue // 继续处理下一个文件
		}
		// 将成功的结果添加到切片中
		results = append(results, renameResult{
			originalPath: op.originalPath,
			finalPath:    op.finalPath,
		})
	}

	return results, nil
}
