package handler

import (
	"fmt"
	"path/filepath"
	"tmrn/internal/service"
	"tmrn/pkg/util"

	"github.com/fatih/color"
)

var (
    successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
)

type RenameHandler struct {
	finder  service.Finder
	renamer service.Renamer
}

func NewRenameHandler(f service.Finder, r service.Renamer) *RenameHandler {
	return &RenameHandler{
		finder:  f,
		renamer: r,
	}
}

func (s *RenameHandler) Run(dirPath, fileExt string, reverseSort bool) error {
	// 1. 查找文件
	files, err := s.finder.FindFiles(dirPath, fileExt)
	if err != nil {
		return fmt.Errorf("查找文件时出错: %w", err)
	}

	if len(files) == 0 {
		warnColor.Printf("没有找到匹配的文件\n")
		return nil
	}

	// 2. 向用户确认
	if !util.AskForConfirmation(fmt.Sprintf("是否重命名 %d 个文件?", len(files))) {
		warnColor.Printf("操作已取消\n")
		return nil
	}

	// 3. 重命名文件
	results, err := s.renamer.RenameFiles(dirPath, files, reverseSort)
	if err != nil {
		return err
	}

	// 4. 打印成功的结果
	for _, result := range results {
		successColor.Printf("%s -> %s\n", filepath.Base(result.OriginalPath), filepath.Base(result.FinalPath))
	}

	if len(results) > 0 {
		fmt.Printf("\n一共完成 %d/%d 个文件\n", len(results), len(files))
	}

	return nil
}
