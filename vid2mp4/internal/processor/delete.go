package processor

import (
	"fmt"
	"os"
	"sync"
)

// DeletionResult 存储删除操作的结果
type DeletionResult struct {
	SuccessfullyDeleted []string         // 成功删除的文件路径列表
	FailedDeletions     map[string]error // 删除失败的文件路径及其错误
}

// 在 goroutine 之间传递, 删除结果
type deletionInfo struct {
	path string
	err  error
}

// DeleteOriginalVideo 尝试删除指定的原始文件, 并返回详细的成功与失败报告
func (p *Processor) DeleteOriginalVideo(ps *ProcessResult) (*DeletionResult, error) {
	// 从 map 中提取原始文件路径用于删除
	files := make([]string, 0, len(ps.SuccessJobs))
	for path := range ps.SuccessJobs {
		files = append(files, path)
	}

	stats := &DeletionResult{
		SuccessfullyDeleted: make([]string, 0),
		FailedDeletions:     make(map[string]error),
	}

	// 创建通道
	jobs := make(chan string)
	delInfoChan := make(chan deletionInfo)

	// 任务计数器
	var wg sync.WaitGroup

	// 消费者 因为是 I/O 密集, 所以可以多一点
	numWorkers := 100
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			deleteWorker(jobs, delInfoChan)
		}()
	}

	// 生产者
	go deleteProducer(files, jobs)

	// 协调者 等待结束
	go func() {
		wg.Wait()
		close(delInfoChan)
	}()

	// 收集者
	for asda := range delInfoChan {
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

func deleteWorker(jobs <-chan string, delInfoChan chan<- deletionInfo) {
	for path := range jobs {
		err := os.Remove(path)
		delInfoChan <- deletionInfo{path: path, err: err}
	}
}

// deleteProducer 将文件路径发送到 jobs 通道
func deleteProducer(files []string, jobs chan<- string) {
	defer close(jobs)
	for _, path := range files {
		jobs <- path
	}
}
