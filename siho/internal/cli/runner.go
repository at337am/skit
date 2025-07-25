package cli

import (
	"errors"
	"fmt"
	"os"
	"siho/internal/cryptor"
	"siho/internal/handler"

	"golang.org/x/term"
)

// Runner 存储选项参数
type Runner struct {
	Path      string // 传入的路径
	Decrypt   bool   // 解密
	OutputDir string // 指定输出目录
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	// 检查路径参数
	if r.Path == "" {
		return errors.New("待处理的路径为空")
	}

	if _, err := os.Stat(r.Path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("路径不存在: %s", r.Path)
		}
		return fmt.Errorf("无法访问路径 %s: %w", r.Path, err)
	}

	// 根据加密还是解密模式决定输出目录
	if r.OutputDir == "" {
		if r.Decrypt {
			r.OutputDir = "result_decrypt"
		} else {
			r.OutputDir = "result_encrypt"
		}
	}

	// 检查输出路径
	if info, err := os.Stat(r.OutputDir); err == nil {
		// 输出路径存在的时候
		if !info.IsDir() {
			return fmt.Errorf("输出路径存在但不是目录: %s", r.OutputDir)
		}
	} else {
		if errors.Is(err, os.ErrNotExist) {
			if mkErr := os.MkdirAll(r.OutputDir, 0755); mkErr != nil {
				return fmt.Errorf("创建输出目录失败: %w", mkErr)
			}
		} else {
			return fmt.Errorf("检查输出路径失败: %w", err)
		}
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	if r.Decrypt {
		fmt.Printf("准备开始解密\n")
	} else {
		fmt.Printf("准备开始加密\n")
	}

	// 1. 获取密码
	password, err := getPassword(!r.Decrypt)
	if err != nil {
		return err
	}

	// 2. 依赖注入
	c, err := cryptor.NewPasswordCryptor(password)
	if err != nil {
		return fmt.Errorf("初始化对称加密结构时出错: %w", err)
	}
	h := handler.NewHandler(r.Path, r.OutputDir, c)

	// 3. 执行操作
	if r.Decrypt {
		return h.HandleDecrypt()
	}
	return h.HandleEncrypt()
}

// getPassword 负责与用户交互以获取密码, 并根据需要进行二次确认
func getPassword(withConfirmation bool) (string, error) {
	fmt.Print("请输入密码: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("读取密码失败: %w", err)
	}
	fmt.Println() // 读取后换行，保持终端整洁

	password := string(passwordBytes)
	if password == "" {
		return "", errors.New("密码不能为空")
	}

	// 如果不需要二次确认 (例如解密模式), 则直接返回
	if !withConfirmation {
		return password, nil
	}

	// --- 执行二次确认 ---
	fmt.Print("请再次输入密码进行确认: ")
	confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("读取确认密码失败: %w", err)
	}
	fmt.Println()

	if password != string(confirmBytes) {
		return "", errors.New("两次输入的密码不匹配")
	}

	return password, nil
}
