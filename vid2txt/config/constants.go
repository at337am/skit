package config

import "strings"

const (
	TaskStatusOngoing   = "ONGOING"
	TaskStatusCompleted = "COMPLETED"
	TaskStatusFailed    = "FAILED"
)

// ANSI 颜色代码
const (
    ColorReset  = "\033[0m"
    ColorCyan   = "\033[36m"
    ColorGreen  = "\033[32m"
    ColorYellow = "\033[33m"
    ColorRed    = "\033[31m"
)

var PromptMenu = strings.Join([]string{
    ColorCyan + "┌──────────────────────────┐" + ColorReset,
    ColorCyan + "│  " + ColorYellow + "📌 请选择操作：" + ColorCyan + "         │" + ColorReset,
    ColorCyan + "├──────────────────────────┤" + ColorReset,
    ColorCyan + "│ " + ColorGreen + "1 - 视频转文字" + ColorCyan + "           │" + ColorReset,
    ColorCyan + "│ " + ColorGreen + "2 - 音频转文字" + ColorCyan + "           │" + ColorReset,
    ColorCyan + "│ " + ColorGreen + "3 - 查询任务状态" + ColorCyan + "         │" + ColorReset,
    ColorCyan + "│ " + ColorGreen + "4 - 修改语言配置" + ColorCyan + "         │" + ColorReset,
    ColorCyan + "└──────────────────────────┘" + ColorReset,
}, "\n")
