package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"vid2mp4/internal/converter"
)

// ProcessResult 存储目录处理过程中的统计信息
type ProcessResult struct {
	SuccessJobs  map[string]*converter.ConvertResult // 存储成功的结果结构体
	FailedJobs   map[string]error                    // 存储失败路径和具体错误
	AccessErrors map[string]error                    // 存储遍历时访问路径的错误
}

// 在 goroutine 之间传递, 处理结果
type processingInfo struct {
	vidPath       string
	convertResult *converter.ConvertResult
	accessError   error
	convertError  error
}

// ProcessVideoDir 遍历指定目录及其子目录, 转换所有指定扩展名的文件
func (p *Processor) ProcessVideoDir(directory, extension, outputDir string) (*ProcessResult, error) {
	stats := &ProcessResult{
		SuccessJobs:  make(map[string]*converter.ConvertResult),
		FailedJobs:   make(map[string]error),
		AccessErrors: make(map[string]error),
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
			p.processWorker(jobs, procInfoChan, outputDir)
		}()
	}

	// 启动生产者, 负责查找文件
	go processProducer(jobs, procInfoChan, directory, extension)

	// 协调者, 等待所有任务完成后关闭 procInfoChan
	go func() {
		wg.Wait()
		close(procInfoChan)
	}()

	// 收集者, 在主 goroutine 中收集所有结果
	for procInfo := range procInfoChan {
		if procInfo.accessError != nil {
			stats.AccessErrors[procInfo.vidPath] = procInfo.accessError
		} else if procInfo.convertError != nil {
			stats.FailedJobs[procInfo.vidPath] = procInfo.convertError
		} else {
			stats.SuccessJobs[procInfo.vidPath] = procInfo.convertResult
		}
	}

	if len(stats.AccessErrors) > 0 || len(stats.FailedJobs) > 0 {
		return stats, fmt.Errorf("视频转换失败/路径访问失败")
	}

	return stats, nil
}

// processWorker 从 jobs 通道接收文件路径, 进行转换, 并将结果发送到 procInfoChan。
func (p *Processor) processWorker(jobs <-chan string, procInfoChan chan<- processingInfo, outputDir string) {
	for filePath := range jobs {
		result, err := p.conv.ConvertToMP4(filePath, outputDir)
		procInfoChan <- processingInfo{
			vidPath:       filePath,
			convertResult: result,
			convertError:  err,
		}
	}
}

// processProducer 遍历目录, 将找到的指定扩展名文件路径发送到 jobs 通道, 并通过 procInfoChan 报告访问错误
func processProducer(jobs chan<- string, procInfoChan chan<- processingInfo, directory, extension string) {
	defer close(jobs) // 生产者完成工作后, 关闭 jobs 通道

	filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			procInfoChan <- processingInfo{vidPath: path, accessError: err}
			return nil
		}
		if !d.IsDir() && strings.ToLower(filepath.Ext(path)) == extension {
			jobs <- path // 将任务路径发送给 Worker
		}
		return nil
	})
}
