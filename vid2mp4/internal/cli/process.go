package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// processResult 存储目录处理过程中的统计信息
type processResult struct {
	successJobs map[string]*ConvertResult // 存储成功的原始路径和结果
	failedJobs  map[string]error          // 存储失败路径和具体错误
}

// processingInfo 在 goroutine 之间传递, 处理结果
type processingInfo struct {
	vidPath       string
	convertResult *ConvertResult
	convertError  error
}

// processBatch 遍历指定目录及其子目录, 转换所有指定扩展名的文件
func (r *Runner) processBatch() (*processResult, error) {
	stats := &processResult{
		successJobs: make(map[string]*ConvertResult),
		failedJobs:  make(map[string]error),
	}

	filesToProcess, err := collectFilesToProcess(r.InputPaths)
	if err != nil {
		return nil, err
	}

	if len(filesToProcess) == 0 {
		return stats, nil
	}

	// 同步计数器
	var wg sync.WaitGroup

	// 传递结果
	procInfoChan := make(chan processingInfo)

	// 传递工作任务
	jobs := make(chan string)

	// 启动消费者等待任务
	// 根据 CPU 数量启动固定数量的 worker goroutine
	// 因为是 CPU 密集型, 所以少一点, 如果视频遇到转码, 可能会占满CPU
	numWorkers := min(runtime.NumCPU(), 4)
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			r.processWorker(jobs, procInfoChan)
		}()
	}

	// 启动生产者, 负责查找文件
	go func() {
		defer close(jobs)
		for _, path := range filesToProcess {
			jobs <- path
		}
	}()

	// 协调者, 等待所有任务完成后关闭 procInfoChan
	go func() {
		wg.Wait()
		close(procInfoChan)
	}()

	// 收集者, 在主 goroutine 中收集所有结果
	for procInfo := range procInfoChan {
		if procInfo.convertError != nil {
			stats.failedJobs[procInfo.vidPath] = procInfo.convertError
		} else {
			stats.successJobs[procInfo.vidPath] = procInfo.convertResult
		}
	}

	return stats, nil
}

// processWorker 从 jobs 通道接收文件路径, 进行转换, 并将结果发送到 procInfoChan。
func (r *Runner) processWorker(jobs <-chan string, procInfoChan chan<- processingInfo) {
	for filePath := range jobs {
		result, err := r.conv.ConvertToMP4(filePath, r.OutputDir)
		procInfoChan <- processingInfo{
			vidPath:       filePath,
			convertResult: result,
			convertError:  err,
		}
	}
}

// collectFilesToProcess 收集文件, 自动跳过目录和 mp4, 遇到无法访问的文件直接报错
func collectFilesToProcess(paths []string) ([]string, error) {
	var validFiles []string
	for _, path := range paths {
		// 1. 检查文件状态
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("无法访问文件/路径: %s, 错误: %w", path, err)
		}

		// 2. 跳过目录
		if info.IsDir() {
			warnColor.Printf("跳过目录: %s\n", path)
			continue
		}

		// 3. 跳过 .mp4 文件
		if strings.EqualFold(filepath.Ext(path), ".mp4") {
			warnColor.Printf("跳过 MP4 文件: %s\n", path)
			continue
		}

		validFiles = append(validFiles, path)
	}
	return validFiles, nil
}
