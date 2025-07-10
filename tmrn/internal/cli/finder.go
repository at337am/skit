package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// fileInfo 单个文件的信息
type fileInfo struct {
	path    string
	modTime time.Time
	ext     string
}

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
