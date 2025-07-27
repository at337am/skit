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

func (r *Runner) handleDir() error {
	fmt.Printf("准备处理目录...\n")

	// 调用执行, 返回的错误是描述性的
	result, err := r.processDir()
	if err != nil {
		warnColor.Printf("处理目录时发生错误: %v\n", err)
	}

	// ========= 显示结果 =========
	successCount := len(result.successJobs)
	failedCount := len(result.failedJobs)
	accessErrCount := len(result.accessErrors)

	if successCount == 0 && failedCount == 0 && accessErrCount == 0 {
		warnColor.Printf("指定目录内没有找到 %s 视频文件\n", r.Extension)
	}

	if successCount > 0 {
		for inputPath, convertResult := range result.successJobs {
			successColor.Printf("转换成功: %s -> %s\n", inputPath, convertResult.OutputPath)
			warnColor.Printf("  └─ %s\n", convertResult.StatusMessage) // 显示转换状态详情
		}
	}

	if failedCount > 0 {
		for path, err := range result.failedJobs {
			errorColor.Printf("转换失败: %s -> %v\n", path, err)
		}
	}

	if accessErrCount > 0 {
		for path, err := range result.accessErrors {
			warnColor.Printf("访问错误: %s -> %v\n", path, err)
		}
	}

	fmt.Printf("\n处理完毕...\n")

	// ========= 删除逻辑 =========
	if len(result.successJobs) > 0 {
		if r.AutoRemove || askForConfirmation("是否删除已成功转换的原始文件?") {
			for filePath := range result.successJobs {
				if err := os.Remove(filePath); err != nil {
					errorColor.Printf("删除失败 -> %s 错误: %v\n", filePath, err)
				} else {
					warnColor.Printf("已删除 -> %s\n", filePath)
				}
			}
		} else {
			warnColor.Printf("\n操作取消, 保留所有原始文件\n")
		}
	}

	return nil
}

// processDir 遍历指定目录及其子目录, 转换所有指定扩展名的文件
func (r *Runner) processDir() (*processResult, error) {
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

// processProducer 遍历目录, 将找到的指定扩展名文件路径发送到 jobs 通道, 并通过 procInfoChan 报告访问错误
func (r *Runner) processProducer(jobs chan<- string, procInfoChan chan<- processingInfo) {
	filepath.WalkDir(r.InputPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			procInfoChan <- processingInfo{vidPath: path, accessError: err}
			return nil
		}
		if !d.IsDir() && strings.EqualFold(filepath.Ext(path), r.Extension) {
			jobs <- path // 将任务路径发送给 Worker
		}
		return nil
	})
}
