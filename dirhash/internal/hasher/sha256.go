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

type SHA256Hasher struct{}

func NewSHA256Hasher() *SHA256Hasher {
	return &SHA256Hasher{}
}

func (s *SHA256Hasher) HashFile(filePath string) (string, error) {
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

func (s *SHA256Hasher) HashDir(dirPath string) (map[string]string, error) {
	// 定义一个用于在 channel 中传递结果的结构体
	type result struct {
		path    string // 子文件的相对路径
		hash    string // 子文件的哈希值
		walkErr error  // 遍历时发生的错误
		hashErr error  // 计算哈希时发生的错误
	}

	numWorkers := runtime.NumCPU()

	// 应优先扩大 jobs 的容量, 确保 Worker 有足够原料可处理
	// 机器不能因为缺货而停转, 优先保障原材料供应链
	jobs := make(chan string, numWorkers*2)

	// 无需为 results 设置额外缓冲, 哪怕完全无缓冲也不会拖慢 Worker
	// 主 Goroutine 仅需 append(allResults, res), 操作极快, 几乎零延迟
	results := make(chan result)

	var wg sync.WaitGroup

	// 启动消费者
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			for path := range jobs {
				hash, err := s.HashFile(path)
				results <- result{
					path:    path,
					hash:    hash,
					hashErr: err,
				}
			}
		}()
	}

	// 等待所有 worker 完成后关闭 results channel
	go func() {
		defer close(results)
		wg.Wait()
	}()

	// 启动生产者
	go func() {
		defer close(jobs)
		filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				results <- result{path: path, walkErr: err}
				return nil
			}

			// 只计算普通文件
			if d.Type().IsRegular() {
				jobs <- path
			}

			return nil
		})
	}()

	// 收集所有中间结果, 不立即处理错误
	// 防止 worker 还没完成, 这边就直接返回了, 所以先收集
	var allResults []result
	for res := range results {
		allResults = append(allResults, res)
	}

	// 收集所有结果
	hashMap := make(map[string]string)
	for _, res := range allResults {
		if res.walkErr != nil {
			return nil, fmt.Errorf("遍历目录 '%s' 时出错: %w", res.path, res.walkErr)
		}
		if res.hashErr != nil {
			return nil, fmt.Errorf("文件 '%s' 计算哈希失败: %w", res.path, res.hashErr)
		}

		relativePath, err := filepath.Rel(dirPath, res.path)
		if err != nil {
			return nil, fmt.Errorf("获取 '%s' 相对路径时出错: %w", res.path, err)
		}

		hashMap[relativePath] = res.hash
	}

	return hashMap, nil
}
