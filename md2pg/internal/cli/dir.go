package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func (r *Runner) processDir() error {
	// jobResult 包含单个文件转换的结果
	type jobResult struct {
		outputPath string // 输出路径
		sourcePath string // 原始路径
		convertErr error  // 转换失败的 err
		walkErr    error  // 路径访问错误的 err
	}

	numWorkers := runtime.NumCPU()

	jobs := make(chan string, numWorkers*2)
	results := make(chan jobResult)

	var wg sync.WaitGroup

	// 1. 启动 worker pool
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			for path := range jobs {
				// 计算源文件相对于输入根目录的相对路径
				relPath, err := filepath.Rel(r.Path, path)
				if err != nil {
					results <- jobResult{sourcePath: path, convertErr: fmt.Errorf("无法计算相对路径: %w", err)}
					continue
				}

				// 构建保持目录结构的输出路径
				htmlRelPath := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + ".html"
				outputPath := filepath.Join(r.OutputDir, htmlRelPath)

				err = convert(path, outputPath)
				results <- jobResult{
					outputPath: outputPath,
					sourcePath: path,
					convertErr: err,
				}
			}
		}()
	}

	// 2. 启动 goroutine 扫描目录并发送任务
	go func() {
		defer close(jobs) // 任务发送完毕后, 关闭 jobs 通道
		filepath.WalkDir(r.Path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// 发现路径访问错误, 发送到 results 通道并继续遍历
				results <- jobResult{sourcePath: path, walkErr: err}
				return nil
			}
			if !d.IsDir() && strings.EqualFold(filepath.Ext(d.Name()), ".md") {
				jobs <- path
			}
			return nil
		})
	}()

	// 3. 启动 goroutine 等待所有 worker 完成后关闭 results 通道
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. 收集所有结果
	var convertErrors []error
	var walkErrors []error

	for res := range results {
		if res.walkErr != nil {
			// 收集路径访问错误
			err := fmt.Errorf("路径错误: %s, 错误: %w", res.sourcePath, res.walkErr)
			walkErrors = append(walkErrors, err)
			errorColor.Printf("Path Error -> %s\n", res.sourcePath)
		} else if res.convertErr != nil {
			// 收集文件转换错误
			err := fmt.Errorf("失败文件: %s, 错误: %w", res.sourcePath, res.convertErr)
			convertErrors = append(convertErrors, err)
			errorColor.Printf("Failed -> %s\n", res.sourcePath)
		} else {
			successColor.Printf("Converted -> %s\n", res.outputPath)
		}
	}

	// 合并所有错误
	allErrors := append(walkErrors, convertErrors...)
	if len(allErrors) > 0 {
		return errors.Join(allErrors...)
	}

	return nil
}
