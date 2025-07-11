package cli

import "fmt"

func (r *Runner) compareFile() error {
	hash1, err := r.hash.HashFile(r.Path1)
	if err != nil {
		return fmt.Errorf("路径: '%s' 计算哈希时出错: %w", r.Path1, err)
	}

	hash2, err := r.hash.HashFile(r.Path2)
	if err != nil {
		return fmt.Errorf("路径: '%s' 计算哈希时出错: %w", r.Path2, err)
	}

	if hash1 == hash2 {
		sameColor.Printf("\n两个文件内容完全一致!\n")
		fmt.Printf("\nSHA-256: %s\n", hash1)
	} else {
		diffColor.Printf("\n两个文件内容不一致!\n")
		fmt.Printf("\n文件: %s\n", r.Path1)
		diffColor.Printf("  └─ SHA-256: %s\n", hash1)
		fmt.Printf("\n文件: %s\n", r.Path2)
		diffColor.Printf("  └─ SHA-256: %s\n", hash2)
	}

	return nil
}
