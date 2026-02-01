package cli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"xla/internal/handler"
	"xla/internal/translator"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var (
	successColor = color.New(color.FgGreen)
	warnColor    = color.New(color.FgCyan)
	errorColor   = color.New(color.FgRed)
)

// Runner 存储选项参数
type Runner struct {
	SourceLang string // 源语言
	TargetLang string // 目标语言
	Proxy      string // 代理地址
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	return nil
}

func (r *Runner) Run() error {
	client, err := createHTTPClient(r.Proxy)
	if err != nil {
		return err
	}

	tran := translator.NewGoogleTranslate(client)
	h := handler.NewHandler(tran)

	// 定义提示符, 并使用 color 包生成带 ANSI escape code 的彩色字符串
	modePrompt := fmt.Sprintf("[%s -> %s] > ", r.SourceLang, r.TargetLang)
	coloredPrompt := warnColor.Sprint(modePrompt)

	// 创建 readline 实例, 并将彩色提示符传入配置
	rl, err := readline.NewEx(&readline.Config{
		Prompt: coloredPrompt,
	})
	if err != nil {
		return fmt.Errorf("创建 readline 实例失败: %w", err)
	}
	// 确保程序退出时关闭 readline, 这会恢复终端的原始状态
	defer rl.Close()

	for {
		// Readline 会自动显示提示符并等待用户输入
		line, err := rl.Readline()

		// 在循环内部处理错误
		if err == readline.ErrInterrupt {
			// 用户按下了 Ctrl+C
			// 如果当前行没有内容, 通常意味着用户想退出程序
			if len(line) == 0 {
				break
			}
			// 如果当前行有内容, 则只清空当前行, 继续循环
			continue
		} else if err == io.EOF {
			// 用户按下了 Ctrl+D (End Of File)
			break
		} else if err != nil {
			// 其他类型的读取错误
			errorColor.Fprintf(os.Stderr, "读取输入时出错: %v\n", err)
			break
		}

		text := strings.TrimSpace(line)

		// 检查退出命令 (不区分大小写)
		if strings.EqualFold(text, "q") {
			break
		}

		// 如果输入为空, 则继续等待下一次输入
		if text == "" {
			continue
		}

		// 调用翻译处理器
		result, err := h.Process(text, r.SourceLang, r.TargetLang)
		if err != nil {
			errorColor.Fprintf(os.Stderr, "翻译出错: %v\n", err)
			continue
		}

		successColor.Println(result)
	}

	fmt.Printf("\nbye~\n")
	return nil
}

// createHTTPClient 是一个辅助函数, 用于根据是否提供代理地址来创建 http.Client
func createHTTPClient(proxyAddr string) (*http.Client, error) {
	if proxyAddr == "" {
		return http.DefaultClient, nil
	}

	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("无效的代理 URL '%s': %w", proxyAddr, err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	return &http.Client{Transport: transport}, nil
}
