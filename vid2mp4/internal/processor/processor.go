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

// DeletionStats 存储删除操作的结果
type DeletionStats struct {
	SuccessfullyDeleted []string         // 成功删除的文件路径列表
	FailedDeletions     map[string]error // 删除失败的文件路径及其错误
}

// ProcessingStats 存储目录处理过程中的统计信息
type ProcessingStats struct {
	SuccessJobs  map[string]*converter.ConvertResult // 存储成功的结果结构体
	FailedJobs   map[string]error                    // 存储失败路径和具体错误
	AccessErrors map[string]error                    // 存储遍历时访问路径的错误
}

// 在 goroutine 之间传递, 处理结果
type processingResult struct {
	vidPath       string
	convertResult *converter.ConvertResult
	accessError   error
	convertError  error
}

// 在 goroutine 之间传递, 删除结果
type deleteResult struct {
	path string
	err  error
}

// --------------- 处理逻辑 ---------------

// findFiles 遍历目录，将找到的 .mov 文件路径发送到 jobs 通道，并通过 resultsChan 报告访问错误。
func findFiles(jobs chan<- string, resultsChan chan<- processingResult, directory string) {
	defer close(jobs) // 生产者完成工作后，关闭 jobs 通道

	filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			resultsChan <- processingResult{vidPath: path, accessError: err}
			return nil
		}
		if !d.IsDir() && strings.ToLower(filepath.Ext(path)) == ".mov" {
			jobs <- path // 将任务路径发送给 Worker
		}
		return nil
	})
}

// worker 从 jobs 通道接收文件路径，进行转换，并将结果发送到 resultsChan。
func worker(jobs <-chan string, resultsChan chan<- processingResult) {
	for filePath := range jobs {
		result, err := converter.ConvertToMP4(filePath)
		resultsChan <- processingResult{
			vidPath:       filePath,
			convertResult: result,
			convertError:  err,
		}
	}
}

// ProcessDirectory 遍历指定目录及其子目录, 转换所有 .mov 文件
func ProcessDirectory(directory string) (*ProcessingStats, error) {
	stats := &ProcessingStats{
		SuccessJobs:  make(map[string]*converter.ConvertResult),
		FailedJobs:   make(map[string]error),
		AccessErrors: make(map[string]error),
	}

	// 同步计数器
	var wg sync.WaitGroup

	// 传递结果
	resultsChan := make(chan processingResult)

	// 传递工作任务
	jobs := make(chan string)

	// 1. 先启动消费者 (Workers)
	// 根据 CPU 数量启动固定数量的 worker goroutine
	numWorkers := min(runtime.NumCPU(), 4)
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			worker(jobs, resultsChan)
		}()
	}

	// 2. 再启动生产者 (Producer)
	// 启动生产者 goroutine，负责查找文件
	go findFiles(jobs, resultsChan, directory)

	// 启动一个 goroutine，等待所有任务完成后关闭 resultsChan
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集者, 在主 goroutine 中收集所有结果
	for res := range resultsChan {
		if res.accessError != nil {
			stats.AccessErrors[res.vidPath] = res.accessError
		} else if res.convertError != nil {
			stats.FailedJobs[res.vidPath] = res.convertError
		} else {
			stats.SuccessJobs[res.vidPath] = res.convertResult
		}
	}

	if len(stats.AccessErrors) > 0 || len(stats.FailedJobs) > 0 {
		return stats, fmt.Errorf("视频转换失败/路径访问失败")
	}

	return stats, nil
}

// --------------- 删除逻辑 ---------------

// producerD 将文件路径发送到 jobs 通道
func producerD(files []string, jobs chan<- string) {
	defer close(jobs)
	for _, path := range files {
		jobs <- path
	}
}

func workerD(jobs <-chan string, results chan<- deleteResult) {
	for path := range jobs {
		err := os.Remove(path)
		results <- deleteResult{path: path, err: err}
	}
}

// DeleteOriginals 尝试删除指定的原始文件, 并返回详细的成功与失败报告
// 如果有任何文件删除失败, 它会额外返回一个概括性的错误
func DeleteOriginals(files []string) (*DeletionStats, error) {
	stats := &DeletionStats{
		SuccessfullyDeleted: make([]string, 0),
		FailedDeletions:     make(map[string]error),
	}

	// 创建通道
	jobs := make(chan string)
	resultsChan := make(chan deleteResult)

	// 任务计数器
	var wg sync.WaitGroup

	// 消费者
	numWorkers := min(runtime.NumCPU(), 4)
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			workerD(jobs, resultsChan)
		}()
	}

	// 生产者
	go producerD(files, jobs)

	// 经理 goroutine 等待结束的 wg.Wait()
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 收集者
	for asda := range resultsChan {
		if asda.err != nil {
			stats.FailedDeletions[asda.path] = asda.err
		} else {
			stats.SuccessfullyDeleted = append(stats.SuccessfullyDeleted, asda.path)
		}
	}

	if len(stats.FailedDeletions) > 0 {
		return stats, fmt.Errorf("无法删除 %d 个视频文件 (共 %d 个)", len(stats.FailedDeletions), len(files))
	}

	return stats, nil
}
