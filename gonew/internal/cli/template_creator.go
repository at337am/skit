package cli

import (
	"fmt"
	"gonew/assets"
	"io/fs"
	"os"
	"path/filepath"
)

// createProjectFromTemplate 从嵌入的模板目录创建项目文件
func (r *Runner) createProjectFromTemplate(templateDir string) error {
	projectDir := r.ProjectName

	// 创建项目根目录
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("创建项目目录 '%s' 时发生错误: %w", projectDir, err)
	}

	// 遍历嵌入文件系统中的模板文件
	return fs.WalkDir(assets.FS, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过模板根目录本身
		if path == templateDir {
			return nil
		}

		// 构建目标路径, 移除模板目录前缀
		relativePath, _ := filepath.Rel(templateDir, path)
		destPath := filepath.Join(projectDir, relativePath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("创建目录 '%s' 时发生错误: %w", destPath, err)
			}
		} else {
			content, err := assets.FS.ReadFile(path)
			if err != nil {
				return fmt.Errorf("读取模板文件 '%s' 时发生错误: %w", path, err)
			}
			if err := os.WriteFile(destPath, content, 0644); err != nil {
				return fmt.Errorf("创建文件 '%s' 时发生错误: %w", destPath, err)
			}
		}
		return nil
	})
}
