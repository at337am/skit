package processor

import (
	"vid2mp4/internal/service/converter"
)

// ProcessResult 存储目录处理过程中的统计信息
type ProcessResult struct {
	SuccessJobs  map[string]*converter.ConvertResult // 存储成功的结果结构体
	FailedJobs   map[string]error                    // 存储失败路径和具体错误
	AccessErrors map[string]error                    // 存储遍历时访问路径的错误
}

// Processor 定义了处理行为
type Processor interface {
	ProcessVideoDir(directory, extension, outputDir string) (*ProcessResult, error)
}
