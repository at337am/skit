package handler

import (
	"vid2mp4/internal/service/converter"
	"vid2mp4/internal/service/processor"
)

// AppHandler 实现了 Handler 接口, 包含了应用的核心依赖和逻辑
type AppHandler struct {
	conv converter.Converter
	proc processor.Processor
}

// NewAppHandler 创建一个新的 Handler 实例
func NewAppHandler(conv converter.Converter, proc processor.Processor) Handler {
	return &AppHandler{
		conv: conv,
		proc: proc,
	}
}

// Execute 根据配置决定是处理文件还是目录
func (h *AppHandler) Execute(cfg *Config) error {
	if cfg.IsDirectory {
		return h.handleDirectory(cfg)
	}
	return h.handleFile(cfg)
}
