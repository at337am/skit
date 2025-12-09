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
	successJobs  map[string]*ConvertResult // 存储成功的原始路径和结果
	failedJobs   map[string]error          // 存储失败路径和具体错误
	accessErrors map[string]error          // 存储遍历时访问路径的错误
}

// processingInfo 在 goroutine 之间传递, 处理结果
type processingInfo struct {
	vidPath       string
	convertResult *ConvertResult
	accessError   error
	convertError  error
}

// processBatch 遍历指定目录及其子目录, 转换所有指定扩展名的文件
func (r *Runner) processBatch() (*processResult, error) {
	stats := &processResult{
		successJobs:  make(map[string]*ConvertResult),
		failedJobs:   make(map[string]error),
		accessErrors: make(map[string]error),
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
		defer close(jobs) // 生产者完成工作后, 关闭 jobs 通道
		r.processProducer(jobs, procInfoChan)
	}()

	// 协调者, 等待所有任务完成后关闭 procInfoChan
	go func() {
		wg.Wait()
		close(procInfoChan)
	}()

	// 收集者, 在主 goroutine 中收集所有结果
	for procInfo := range procInfoChan {
		if procInfo.accessError != nil {
			stats.accessErrors[procInfo.vidPath] = procInfo.accessError
		} else if procInfo.convertError != nil {
			stats.failedJobs[procInfo.vidPath] = procInfo.convertError
		} else {
			stats.successJobs[procInfo.vidPath] = procInfo.convertResult
		}
	}

	if len(stats.accessErrors) > 0 || len(stats.failedJobs) > 0 {
		return stats, fmt.Errorf("视频转换失败/路径访问失败")
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

// processProducer 遍历 InputPaths, 过滤掉目录和 MP4 文件, 将有效文件发送到 jobs 通道
func (r *Runner) processProducer(jobs chan<- string, procInfoChan chan<- processingInfo) {
	for _, path := range r.InputPaths {
		// 1. 检查路径状态
		info, err := os.Stat(path)
		if err != nil {
			procInfoChan <- processingInfo{vidPath: path, accessError: err}
			continue
		}

		// 2. 如果是目录，提示跳过
		if info.IsDir() {
			warnColor.Printf("跳过目录: %s\n", path)
			continue
		}

		// 3. 如果是 mp4 后缀，提示跳过
		if strings.EqualFold(filepath.Ext(path), ".mp4") {
			warnColor.Printf("跳过 MP4 文件: %s\n", path)
			continue
		}

		// 4. 发送任务
		jobs <- path
	}
}
