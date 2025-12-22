package handler

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

type Cryptor interface {
	Encrypt(inputPath, outputPath string) error
	Decrypt(inputPath, outputPath string) error
}

type Handler struct {
	FilePaths []string
	OutputDir string
	crypt     Cryptor
}

func NewHandler(paths []string, outputDir string, c Cryptor) *Handler {
	return &Handler{
		FilePaths: paths,
		OutputDir: outputDir,
		crypt:     c,
	}
}

// HandleEncrypt 统一处理文件和目录的加密逻辑
func (h *Handler) HandleEncrypt() error {
	// 定义加密文件的具体操作
	encryptFile := func(inputPath string) (string, error) {
		baseName := filepath.Base(inputPath)
		outputPath := filepath.Join(h.OutputDir, fmt.Sprintf("%s_enc", baseName))
		err := h.crypt.Encrypt(inputPath, outputPath)
		return outputPath, err
	}

	return h.processFiles("Encrypted", encryptFile)
}

// HandleDecrypt 统一处理文件和目录的解密逻辑
func (h *Handler) HandleDecrypt() error {
	// 定义解密文件的具体操作
	decryptFile := func(inputPath string) (string, error) {
		baseName := filepath.Base(inputPath)
		var outputBaseName string
		// 根据文件名是否以 "_enc" 结尾, 决定输出文件名
		if strings.HasSuffix(baseName, "_enc") {
			outputBaseName = strings.TrimSuffix(baseName, "_enc")
		} else {
			outputBaseName = fmt.Sprintf("%s_dec", baseName)
		}
		outputPath := filepath.Join(h.OutputDir, outputBaseName)
		err := h.crypt.Decrypt(inputPath, outputPath)
		return outputPath, err
	}
	return h.processFiles("Decrypted", decryptFile)
}

// processFiles 使用 worker pool 并发处理文件
func (h *Handler) processFiles(opName string, processFunc func(string) (string, error)) error {
	// jobResult 用于在 goroutine 之间传递处理结果
	type jobResult struct {
		inputPath  string
		outputPath string
		err        error
	}

	files, err := collectFilesToProcess(h.FilePaths)
	if err != nil {
		return fmt.Errorf("无法获取待%s的文件: %w", opName, err)
	}
	if len(files) == 0 {
		warnColor.Printf("未找到可%s的文件\n", opName)
		return nil
	}

	// 1. 设置 worker pool
	numWorkers := runtime.NumCPU()
	jobs := make(chan string, numWorkers*2)
	results := make(chan jobResult)
	var wg sync.WaitGroup

	// 2. 启动 workers
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			// 从 jobs 通道接收任务, 直到通道关闭
			for inputPath := range jobs {
				outputPath, err := processFunc(inputPath)
				results <- jobResult{inputPath: inputPath, outputPath: outputPath, err: err}
			}
		}()
	}

	// 3. 等待所有 worker 完成任务后关闭 results 通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. 分发任务
	go func() {
		defer close(jobs)
		for _, inputPath := range files {
			jobs <- inputPath
		}
	}()

	// 5. 收集并处理结果
	var errs []error
	for result := range results {
		if result.err != nil {
			errs = append(errs, fmt.Errorf("文件%s失败: %s, 错误: %w", opName, filepath.Base(result.inputPath), result.err))
			errorColor.Printf("Failed -> %s\n", result.inputPath)
			continue
		}
		successColor.Printf("%s -> %s\n", opName, result.outputPath)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// collectFilesToProcess 收集待处理的文件路径切片
func collectFilesToProcess(paths []string) ([]string, error) {
	var files []string

	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("路径不存在: %s", path)
			}
			return nil, fmt.Errorf("无法访问路径 %s: %w", path, err)
		}

		if info.IsDir() {
			warnColor.Printf("Skipping directory: %s\n", path)
			continue
		}
		files = append(files, path)
	}

	return files, nil
}
