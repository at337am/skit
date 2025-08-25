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
	Path      string // 待处理的路径
	Decrypt   bool   // 解密模式
	OutputDir string // 指定输出目录
	password  string // 输入的密码
}

func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数, 协调执行各个校验步骤
func (r *Runner) Validate() error {
	// 校验待处理的路径
	info, err := r.validateInputPath()
	if err != nil {
		return err
	}

	// 获取并设置密码
	if err := r.acquirePassword(); err != nil {
		return err
	}

	// 设置并准备输出目录
	if err := r.setupAndPrepareOutputDir(info.IsDir()); err != nil {
		return err
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 1. 依赖注入
	c, err := cryptor.NewPasswordCryptor(r.password)
	if err != nil {
		return fmt.Errorf("初始化对称加密结构时出错: %w", err)
	}
	h := handler.NewHandler(r.Path, r.OutputDir, c)

	// 2. 执行操作
	if r.Decrypt {
		return h.HandleDecrypt()
	}
	return h.HandleEncrypt()
}

// validateInputPath 校验待处理的路径是否存在且可访问
func (r *Runner) validateInputPath() (os.FileInfo, error) {
	if r.Path == "" {
		return nil, errors.New("待处理的路径为空")
	}
	info, err := os.Stat(r.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("路径不存在: %s", r.Path)
		}
		return nil, fmt.Errorf("无法访问路径 %s: %w", r.Path, err)
	}
	return info, nil
}

// acquirePassword 提示用户输入并设置密码
func (r *Runner) acquirePassword() error {
	if r.Decrypt {
		fmt.Print("请输入解密所需的密码:")
	} else {
		fmt.Print("请设定一个密码以用于加密:")
	}

	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("读取密码失败: %w", err)
	}
	fmt.Println() // 换行

	password := string(passwordBytes)
	if password == "" {
		return errors.New("密码不能为空")
	}

	// 解密模式不需要二次确认
	if r.Decrypt {
		r.password = password
		return nil
	}

	// 加密模式需要二次确认
	fmt.Print("请再次确认密码:")
	confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("读取确认密码失败: %w", err)
	}
	fmt.Println()

	if password != string(confirmBytes) {
		return errors.New("两次输入的密码不一致")
	}

	r.password = password
	return nil
}

// setupAndPrepareOutputDir 设置并准备输出目录
func (r *Runner) setupAndPrepareOutputDir(isSourceDir bool) error {
	// 1. 如果输出目录未指定, 则设置默认值
	if r.OutputDir == "" {
		if isSourceDir {
			if r.Decrypt {
				r.OutputDir = "decrypted_result"
			} else {
				r.OutputDir = "encrypted_result"
			}
		} else {
			r.OutputDir = "."
		}
	}

	// 2. 确保输出目录存在且是一个目录
	info, err := os.Stat(r.OutputDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if mkErr := os.MkdirAll(r.OutputDir, 0755); mkErr != nil {
				return fmt.Errorf("创建输出目录失败: %w", mkErr)
			}
			return nil // 创建成功
		}
		return fmt.Errorf("检查输出路径失败: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("输出路径存在但不是目录: %s", r.OutputDir)
	}
	return nil
}
