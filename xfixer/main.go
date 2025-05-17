package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/h2non/filetype"
)

// 不匹配文件信息结构体
type MismatchedFile struct {
	Path        string
	DetectedExt string
}

// 检查单个文件
func checkFile(path string) (*MismatchedFile, error) {
	// 读取文件头部用于检测类型
	buf := make([]byte, 261) // 读取前261字节用于文件类型识别
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %v", path, err)
	}
	defer file.Close()

	_, err = file.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("无法读取文件 %s: %v", path, err)
	}

	// 检测文件类型
	kind, err := filetype.Match(buf)
	if err != nil {
		return nil, fmt.Errorf("检测文件类型出错 %s: %v", path, err)
	}

	// 如果文件类型无法识别，则返回nil
	if kind == filetype.Unknown {
		return nil, nil
	}

	// 获取文件当前扩展名
	currentExt := strings.ToLower(filepath.Ext(path))
	if currentExt != "" {
		// 去掉前导点号
		currentExt = currentExt[1:]
	}

	// 获取检测到的扩展名
	detectedExt := kind.Extension

	// 如果扩展名不匹配
	if currentExt != detectedExt {
		return &MismatchedFile{
			Path:        path,
			DetectedExt: detectedExt,
		}, nil
	}

	// 文件类型与扩展名匹配，返回nil
	return nil, nil
}

// 生成不冲突的文件名
func getUniqueFilePath(originalPath string, detectedExt string) string {
	dir := filepath.Dir(originalPath)
	filenameWithoutExt := strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(originalPath))
	
	suggestedPath := filepath.Join(dir, filenameWithoutExt+"."+detectedExt)
	uniquePath := suggestedPath
	counter := 1
	
	// 检查文件是否已存在，如果存在则添加计数器
	for {
		_, err := os.Stat(uniquePath)
		if os.IsNotExist(err) {
			break // 文件不存在，可以使用这个名称
		}
		
		// 生成新的文件名
		uniquePath = filepath.Join(dir, fmt.Sprintf("%s_%d.%s", filenameWithoutExt, counter, detectedExt))
		counter++
	}
	
	return uniquePath
}

// 扫描目录
func scanDirectory(dirPath string) ([]MismatchedFile, error) {
	var mismatchedFiles []MismatchedFile

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		mismatch, err := checkFile(path)
		if err != nil {
			fmt.Printf("💡 处理文件 %s 时异常: %v\n", path, err)
			return nil // 继续处理下一个文件
		}

		if mismatch != nil {
			mismatchedFiles = append(mismatchedFiles, *mismatch)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("❌ 扫描目录出错: %v", err)
	}

	return mismatchedFiles, nil
}

// 修正文件名
func correctFiles(files []MismatchedFile) error {
	for _, file := range files {
		// 生成唯一的目标路径
		uniquePath := getUniqueFilePath(file.Path, file.DetectedExt)
		
		// 执行重命名
		err := os.Rename(file.Path, uniquePath)
		if err != nil {
			return fmt.Errorf("❌ 重命名文件 %s 失败: %v", file.Path, err)
		}
		fmt.Printf(" - %s -> %s\n", file.Path, uniquePath)
	}
	
	return nil
}


// 使用方法信息
func customUsage() {
	fmt.Println("使用方法:")
	fmt.Println("  vidfixer <文件或目录路径>")
}

func main() {
    flag.Usage = customUsage
    flag.Parse()

    args := flag.Args()
    if len(args) != 1 {
        fmt.Println("❌ 请提供一个文件或目录路径")
		os.Exit(1)
    }

    path := args[0]
	
	// 检查路径是文件还是目录
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("❌ 无法访问路径 %s: %v\n", path, err)
		os.Exit(1)
	}

	var mismatchedFiles []MismatchedFile

	if fileInfo.IsDir() {
		// 目录模式
		mismatchedFiles, err = scanDirectory(path)
		if err != nil {
			fmt.Printf("❌ 扫描目录出错: %v\n", err)
			os.Exit(1)
		}
		
		if len(mismatchedFiles) == 0 {
			fmt.Println("\n✅ 所有文件扩展名与检测到的类型匹配")
			os.Exit(0)
		}
		
		fmt.Println("\n🧐 发现以下扩展名不匹配的文件：")
		for _, file := range mismatchedFiles {
			fmt.Printf(" - %s -> %s\n", file.Path, file.DetectedExt)
		}
	} else {
		// 单文件模式
		mismatch, err := checkFile(path)
		if err != nil {
			fmt.Printf("❌ 检查文件出错: %v\n", err)
			os.Exit(1)
		}
		
		if mismatch == nil {
			fmt.Println("\n✅ 文件扩展名与检测到的类型匹配")
			os.Exit(0)
		}

		fmt.Println("\n🧐 发现以下扩展名不匹配的文件：")
		fmt.Printf(" - %s -> %s\n", mismatch.Path, mismatch.DetectedExt)
		
		mismatchedFiles = append(mismatchedFiles, *mismatch)
	}
	
	// 询问用户是否要修正文件
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n✨ 是否修正这些文件? (y/n): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	
	if response == "y" || response == "yes" {
		err = correctFiles(mismatchedFiles)
		if err != nil {
			fmt.Printf("❌ 修正文件时出错: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\n✅ 所有文件已修正")
	} else {
		fmt.Println("❌ 操作已取消，未修改任何文件")
	}
}
