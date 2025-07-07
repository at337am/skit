package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

type FileRenamer struct{}

func NewFileRenamer() Renamer {
	return &FileRenamer{}
}

// renameOp 定义了重命名计划中的一个操作单元
type renameOp struct {
	originalPath string
	tmpPath      string
	finalPath    string
}

const tmpSuffix = ".tmrn-tmp"

func (r *FileRenamer) RenameFiles(dirPath string, files []FileInfo, reverseSort bool) ([]RenameResult, error) {
	if len(files) == 0 {
		return nil, errors.New("没有找到匹配的文件")
	}

	// 依据文件修改时间进行排序
	slices.SortFunc(files, func(a, b FileInfo) int {
		if reverseSort {
			return b.ModTime.Compare(a.ModTime) // 降序
		}
		return a.ModTime.Compare(b.ModTime) // 升序
	})

	// 步骤 1: 构建重命名计划
	plan := createRenamePlan(dirPath, files)

	// 步骤 2: 执行第一阶段重命名（加锁）
	if err := renameToTemp(plan); err != nil {
		// 如果第一阶段失败, 返回 nil 结果和错误
		return nil, err
	}

	// 步骤 3: 执行第二阶段重命名（解锁）, 并收集结果
	results := renameToFinal(plan)

	return results, nil
}

// createRenamePlan 负责根据文件列表生成重命名计划
func createRenamePlan(dirPath string, files []FileInfo) []renameOp {
	numFiles := len(files)
	digits := len(fmt.Sprintf("%d", numFiles))
	formatTemplate := fmt.Sprintf("%%0%dd", digits)
	plan := make([]renameOp, numFiles)

	for i, file := range files {
		finalName := fmt.Sprintf(formatTemplate, i+1) + file.Ext
		plan[i] = renameOp{
			originalPath: file.Path,
			tmpPath:      file.Path + tmpSuffix,
			finalPath:    filepath.Join(dirPath, finalName),
		}
	}
	return plan
}

// renameToTemp 执行第一阶段重命名 (原始文件 -> 临时文件)
// 这一阶段是原子性的，如果中途失败，会尝试回滚所有已成功的操作。
func renameToTemp(plan []renameOp) error {
	for i, op := range plan {
		if err := os.Rename(op.originalPath, op.tmpPath); err != nil {
			// 尝试回滚已成功的重命名操作
			for j := range i {
				// 尽力而为，忽略回滚错误
				_ = os.Rename(plan[j].tmpPath, plan[j].originalPath)
			}
			return fmt.Errorf("操作已中断: 文件 '%s' 重命名失败: %w. 已尝试回滚", filepath.Base(op.originalPath), err)
		}
	}
	return nil
}

// renameToFinal 执行第二阶段重命名 (临时文件 -> 最终文件)
func renameToFinal(plan []renameOp) []RenameResult {
	results := make([]RenameResult, 0, len(plan))
	for _, op := range plan {
		if err := os.Rename(op.tmpPath, op.finalPath); err != nil {
			fmt.Fprintf(os.Stderr, "错误: 无法将 %s 重命名为 %s: %v\n", filepath.Base(op.tmpPath), filepath.Base(op.finalPath), err)
			continue // 继续处理下一个文件
		}
		// 将成功的结果添加到切片中
		results = append(results, RenameResult{
			OriginalPath: op.originalPath,
			FinalPath:    op.finalPath,
		})
	}
	return results
}
