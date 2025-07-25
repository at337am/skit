package cryptor

import (
	"fmt"
	"io"
	"os"

	"filippo.io/age"
)

// PasswordCryptor 结构体中缓存可复用的 Recipient 和 Identity
type PasswordCryptor struct {
	recipient age.Recipient
	identity  age.Identity
}

// NewPasswordCryptor 在构造时就生成 Recipient 和 Identity, 并处理可能发生的错误
func NewPasswordCryptor(p string) (*PasswordCryptor, error) {
	recipient, err := age.NewScryptRecipient(p)
	if err != nil {
		return nil, fmt.Errorf("创建 age recipient 失败: %w", err)
	}

	identity, err := age.NewScryptIdentity(p)
	if err != nil {
		return nil, fmt.Errorf("创建 age identity 失败: %w", err)
	}

	return &PasswordCryptor{
		recipient: recipient,
		identity:  identity,
	}, nil
}

// Encrypt 直接使用预先创建好的 Recipient
func (c *PasswordCryptor) Encrypt(inputPath, outputPath string) (err error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开输入文件失败: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败 '%s': %w", outputPath, err)
	}
	// defer 配合具名返回值 err, 确保在函数退出时执行清理逻辑
	// 如果 err 不为 nil (即加密失败), 则在关闭文件后删除已创建的输出文件
	defer func() {
		outputFile.Close()
		if err != nil {
			os.Remove(outputPath)
		}
	}()

	// 直接复用 c.recipient, 避免重复的密钥派生计算
	wc, err := age.Encrypt(outputFile, c.recipient)
	if err != nil {
		return err
	}

	if _, err = io.Copy(wc, inputFile); err != nil {
		return fmt.Errorf("复制文件内容至加密流时出错: %w", err)
	}

	if err = wc.Close(); err != nil {
		return fmt.Errorf("加密过程中关闭 writer 时出错: %w", err)
	}

	return nil
}

// Decrypt 直接使用预先创建好的 Identity
func (c *PasswordCryptor) Decrypt(inputPath, outputPath string) (err error) {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("打开输入文件出错: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件失败 '%s': %w", outputPath, err)
	}
	// defer 配合具名返回值 err, 确保在函数退出时执行清理逻辑
	// 如果 err 不为 nil (即解密失败), 则在关闭文件后删除已创建的输出文件
	defer func() {
		outputFile.Close()
		if err != nil {
			os.Remove(outputPath)
		}
	}()

	// 直接复用 c.identity
	r, err := age.Decrypt(inputFile, c.identity)
	if err != nil {
		return err
	}

	if _, err = io.Copy(outputFile, r); err != nil {
		return fmt.Errorf("复制解密数据到输出文件时出错: %w", err)
	}

	return nil
}
