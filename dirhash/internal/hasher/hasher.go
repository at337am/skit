package hasher

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// hashFile 计算文件的 SHA256 哈希
func hashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// result 表示一个文件的哈希计算结果或错误
type result struct {
	path string
	hash string
	err  error
}

// GenerateHashMap 为指定目录或文件生成哈希映射
func GenerateHashMap(rootPath string) (map[string]string, error) {
	numWorkers := runtime.NumCPU()
	jobs := make(chan string, numWorkers)
	results := make(chan result, numWorkers)
	var wg sync.WaitGroup

	// 关键执行逻辑：1. 先启动消费者（worker goroutine），它们会尝试从空的 jobs 通道读取，因此会立即阻塞等待任务。
	for range numWorkers {
		go func() {
			for path := range jobs {
				hash, err := hashFile(path)
				relativePath, relErr := filepath.Rel(rootPath, path)
				if relErr != nil {
					results <- result{err: fmt.Errorf("无法获取相对路径 '%s': %w", path, relErr)}
					wg.Done()
					continue
				}

				if err != nil {
					results <- result{err: fmt.Errorf("计算文件哈希失败 '%s': %w", path, err)}
					wg.Done()
					continue
				}

				if path == rootPath {
					relativePath = filepath.Base(rootPath)
				}

				results <- result{path: relativePath, hash: hash}
				wg.Done() // 使 wg 内部的计数器减一
			}
		}()
	}

	// 准备结果收集器, 统一接收所有结果
	hashMap := make(map[string]string)
	var firstErr error
	collectorDone := make(chan struct{})

	go func() {
		// 当 results 信道被关闭, 且所有数据也被读取完毕. 也就是循环结束后才会被 close
		defer close(collectorDone)
		for res := range results {
			if res.err != nil {
				if firstErr == nil {
					firstErr = res.err
				}
				continue
			}
			hashMap[res.path] = res.hash
		}
	}()

	// 2. 然后由生产者（filepath.WalkDir）开始向 jobs 通道发送任务，唤醒阻塞的 worker goroutine 进行处理。
	walkErr := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			// wg.Wait() 是一个阻塞调用，它会等待 wg 的内部计数器减到零。
			// 因此，wg.Add(1) 的作用就是告诉 wg.Wait()：
			// “嘿，你现在需要等待的任务数量又多了一个。”
			wg.Add(1)
			jobs <- path
		}
		return nil
	})

	// 任务发送完毕，关闭 jobs 通道
	close(jobs)

	// 等待所有 worker 完成
	wg.Wait()

	// 所有结果已生成，关闭 results 通道
	close(results)

	// 阻塞, 直到 collectorDone 这个信号通道被关闭时, 才会解除阻塞
	<-collectorDone

	// 返回时检查是否有错误
	if walkErr != nil {
		return nil, walkErr
	}
	if firstErr != nil {
		return nil, firstErr
	}

	return hashMap, nil
}
