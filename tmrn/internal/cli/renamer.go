package cli

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
)

// renameResult 单个文件重命名后的结果
type renameResult struct {
	originalPath string
	finalPath    string
}

// renameFiles 时间排序重命名
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

// randomizeFiles 随机前缀
func (r *Runner) randomizeFiles(files []fileInfo) ([]renameResult, error) {
	if len(files) == 0 {
		return nil, errors.New("没有找到匹配的文件")
	}

	results := make([]renameResult, 0, len(files))

	// 定义字符集 (大小写英文字母)
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for _, file := range files {
		// 生成 4 位随机英文字符
		randBytes := make([]byte, 4)
		for i := range randBytes {
			randBytes[i] = charset[rand.N(len(charset))]
		}
		prefix := string(randBytes)

		// 构造新文件名: asdf_filename.ext
		originalName := filepath.Base(file.path)
		finalName := fmt.Sprintf("%s_%s", prefix, originalName)
		finalPath := filepath.Join(r.DirPath, finalName)

		// 执行重命名
		if err := os.Rename(file.path, finalPath); err != nil {
			warnColor.Fprintf(os.Stderr, "注意: 无法将 %s 重命名为 %s: %v\n", originalName, finalName, err)
			continue
		}

		results = append(results, renameResult{
			originalPath: file.path,
			finalPath:    finalPath,
		})
	}

	return results, nil
}
