package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileInfo struct {
	Path    string
	ModTime time.Time
	Name    string
	Ext     string
}

func main() {
	// 定义命令行参数
	dirPath := flag.String("d", "", "要处理的目录路径")
	fileExt := flag.String("e", "", "要处理的文件格式 (例如: .jpg, .txt)")
	reverseSort := flag.Bool("r", false, "启用从晚到早排序 (默认为从早到晚)")

	// 解析命令行参数
	flag.Parse()

	// 检查必要参数
	if *dirPath == "" {
		fmt.Println("错误: 必须指定目录路径，使用 -d 参数")
		fmt.Println("使用方法: tmrn -d <目录路径> [-e <文件格式>] [-r]")
		return
	}

	// 确保文件格式以点号开头
	if *fileExt != "" && !strings.HasPrefix(*fileExt, ".") {
		*fileExt = "." + *fileExt
	}

	// 获取目录信息 - 使用现代 API
	entries, err := os.ReadDir(*dirPath)
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		return
	}

	// 准备文件信息列表，只包含常规文件且符合指定格式
	var files []FileInfo
	for _, entry := range entries {
		// 只处理文件，跳过目录
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := filepath.Ext(name)

		// 如果指定了文件格式，则只处理匹配的文件
		if *fileExt != "" && !strings.EqualFold(ext, *fileExt) {
			continue
		}

		// 获取文件信息以访问修改时间
		filePath := filepath.Join(*dirPath, name)
		info, err := entry.Info()
		if err != nil {
			fmt.Printf("获取文件信息失败 %s: %v\n", filePath, err)
			continue
		}

		baseName := strings.TrimSuffix(name, ext)

		files = append(files, FileInfo{
			Path:    filePath,
			ModTime: info.ModTime(),
			Name:    baseName,
			Ext:     ext,
		})
	}

	// 检查是否有匹配的文件
	if len(files) == 0 {
		fmt.Println("没有找到匹配的文件")
		return
	}

	// 依据文件修改时间进行排序
	sort.Slice(files, func(i, j int) bool {
		if *reverseSort {
			// 如果使用 -r 参数，则从晚到早排序
			return files[i].ModTime.After(files[j].ModTime)
		}
		// 默认从早到晚排序
		return files[i].ModTime.Before(files[j].ModTime)
	})

	numFiles := len(files)

	fmt.Printf("即将对目录 \"%s\" 中的 %d 个文件进行批量重命名操作。\n", *dirPath, numFiles)
	fmt.Print("是否继续？ (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		fmt.Println("操作已取消。")
		return
	}

	// 根据文件数量确定格式化模板
	var formatTemplate string

	if numFiles <= 9 {
		formatTemplate = "%d"
	} else {
		// 计算需要的位数
		digits := int(math.Log10(float64(numFiles))) + 1
		formatTemplate = "%0" + fmt.Sprintf("%dd", digits)
	}

	// 创建一个跟踪已使用文件名的映射
	usedNames := make(map[string]bool)

	// 批量重命名文件
	for i, file := range files {
		baseNewName := fmt.Sprintf(formatTemplate, i+1) + file.Ext
		newName := baseNewName
		newPath := filepath.Join(*dirPath, newName)

		// 检查文件名是否已存在，如果存在则添加递增的后缀
		counter := 1
		for {
			// 检查该文件是否是我们自己要重命名的源文件 - 如果是，则可以继续重命名
			if filepath.Base(file.Path) == newName {
				break
			}

			// 检查新的文件名是否已经存在或已经被本次操作使用过
			_, fileExists := os.Stat(newPath)
			if (fileExists == nil || usedNames[newName]) && filepath.Base(file.Path) != newName {
				// 生成新的文件名，添加后缀
				baseName := fmt.Sprintf(formatTemplate+"_%d", i+1, counter)
				newName = baseName + file.Ext
				newPath = filepath.Join(*dirPath, newName)
				counter++
			} else {
				// 文件名不存在，可以使用
				break
			}
		}

		// 标记该文件名已被使用
		usedNames[newName] = true

		// 执行重命名
		err := os.Rename(file.Path, newPath)
		if err != nil {
			fmt.Printf("重命名文件 %s 失败: %v\n", file.Path, err)
		} else {
			fmt.Printf("重命名: %s -> %s\n", file.Path, newPath)
		}
	}

	fmt.Printf("完成！共重命名 %d 个文件\n", len(files))
}
