package processor

import (
	"vid2mp4/internal/converter"
)

// IProcessor 定义了处理行为
type IProcessor interface {
	ProcessVideoDir(directory, extension, outputDir string) (*ProcessResult, error)
	DeleteOriginalVideo(ps *ProcessResult) (*DeletionResult, error)
}

type Processor struct {
	conv converter.IConverter // 依赖 IConverter 接口
}

func NewProcessor(c converter.IConverter) IProcessor {
	return &Processor{
		conv: c,
	}
}
