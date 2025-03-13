package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/urfave/cli/v2"
)

// verifyMediaFileTypeCarefully 检查文件是否为 MP4 或 MOV 视频文件
// 返回值: "mp4"、"mov"、"unknown" 或错误信息
func verifyMediaFileTypeCarefully(filePath string) string {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "error: 文件不存在"
	}

	// 使用 MIME 类型检测
	mime, err := mimetype.DetectFile(filePath)
	if err != nil {
		return fmt.Sprintf("error: magic number 检测失败: %v", err)
	}

	// 检查 MIME 类型
	if mime.String() == "video/mp4" {
		return "mp4"
	} else if mime.String() == "video/quicktime" {
		return "mov"
	}

	return "unknown"
}

// 处理单个文件
func processFile(filePath string, fileName string) {
	result := verifyMediaFileTypeCarefully(filePath)
	
	// 处理非媒体文件或错误
	if handleNonMediaResult(result, fileName) {
		return // 如果文件不是媒体文件或发生错误，直接返回
	}
	
	// 处理媒体文件重命名
	handleMediaFileRename(filePath, fileName, result)
}

// 处理非媒体文件结果
// 返回 true 表示无需进一步处理
func handleNonMediaResult(result, fileName string) bool {
	if result == "unknown" {
		fmt.Printf("文件 '%s' 不是 MP4 或 MOV 文件，或者无法确定类型。\n", fileName)
		return true
	}
	
	if strings.HasPrefix(result, "error") {
		fmt.Printf("验证文件 '%s' 过程中出现错误: %s\n", fileName, strings.TrimPrefix(result, "error: "))
		return true
	}
	
	if result != "mp4" && result != "mov" {
		return true
	}
	
	return false
}

// 处理媒体文件重命名
func handleMediaFileRename(filePath, fileName, fileType string) {
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)
	targetExt := "." + fileType
	
	// 若扩展名已正确，无需处理
	if strings.ToLower(ext) == targetExt {
		return
	}
	
	dirPath := filepath.Dir(filePath)
	newFilename := nameWithoutExt + targetExt
	newFilepath := filepath.Join(dirPath, newFilename)
	
	// 处理文件名冲突
	if needNewName, uniquePath, uniqueName := resolveNameConflict(dirPath, newFilename); needNewName {
		newFilepath = uniquePath
		newFilename = uniqueName
	}
	
	// 执行重命名
	performRename(filePath, newFilepath, fileName, newFilename)
}

// 解决文件名冲突，返回是否需要新名称、新路径和新名称
func resolveNameConflict(dirPath, newFilename string) (bool, string, string) {
	newFilepath := filepath.Join(dirPath, newFilename)
	
	// 如果文件不存在，无需新名称
	if _, err := os.Stat(newFilepath); os.IsNotExist(err) {
		return false, "", ""
	}
	
	// 文件存在，需要生成新名称
	nameWithoutExt := strings.TrimSuffix(newFilename, filepath.Ext(newFilename))
	targetExt := filepath.Ext(newFilename)
	
	counter := 1
	for {
		uniqueFilename := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, targetExt)
		uniqueFilepath := filepath.Join(dirPath, uniqueFilename)
		
		if _, err := os.Stat(uniqueFilepath); os.IsNotExist(err) {
			return true, uniqueFilepath, uniqueFilename
		}
		counter++
	}
}

// 执行文件重命名
func performRename(oldPath, newPath, oldName, newName string) {
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Printf("重命名文件 '%s' 出错: %v\n", oldName, err)
	} else {
		fmt.Printf("✅ 后缀名修正: '%s' -> '%s'\n", oldName, newName)
	}
}

// 处理目录
func processDirectory(dirPath string) {
	// 检查目录是否存在
	dirInfo, err := os.Stat(dirPath)
	if err != nil || !dirInfo.IsDir() {
		fmt.Printf("错误: '%s' 不是一个有效的目录。\n", dirPath)
		return
	}

	// 读取目录内容
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Printf("读取目录 '%s' 出错: %v\n", dirPath, err)
		return
	}

	// 处理每个条目
	for _, entry := range entries {
		entryPath := filepath.Join(dirPath, entry.Name())
		
		if entry.IsDir() {
			// 处理子目录
			fmt.Printf("\n--- 发现子目录: '%s' ---\n", entry.Name())
			processDirectory(entryPath)
		} else {
			// 处理文件
			processFile(entryPath, entry.Name())
		}
	}

	fmt.Printf("\n--- MIME 类型清洗完成 ---\n")
}

func main() {
	app := &cli.App{
		Name:  "video-checker",
		Usage: "验证指定路径下文件的视频格式，并根据真实格式重命名文件后缀 (仅限 MP4 和 MOV)",
		Action: func(c *cli.Context) error {
			if c.NArg() < 1 {
				return fmt.Errorf("请提供目录路径")
			}
			dirPath := c.Args().First()
			processDirectory(dirPath)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("程序运行出错: %v\n", err)
	}
}
