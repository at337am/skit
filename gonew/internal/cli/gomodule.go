package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// initializeGoModule 执行 go mod init 和 go mod tidy
func (r *Runner) initializeGoModule() error {
	projectDir := r.ProjectName

	// 执行 go mod init 命令
	initCmd := exec.Command("go", "mod", "init", projectDir)
	initCmd.Dir = projectDir
	var initStderr bytes.Buffer
	initCmd.Stderr = &initStderr
	if err := initCmd.Run(); err != nil {
		// 关键错误: 不在此处打印, 而是包装更详细的错误信息并返回
		return fmt.Errorf("执行 'go mod init %s' 时发生错误: %w. 详情: %s", projectDir, err, initStderr.String())
	}

	// 执行 go mod tidy 命令
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = projectDir
	var tidyStderr bytes.Buffer
	tidyCmd.Stderr = &tidyStderr
	if err := tidyCmd.Run(); err != nil {
		// 非致命错误: 打印警告信息, 但不返回 error, 不中断流程
		warnColor.Fprintf(os.Stderr, "执行 go mod tidy 时发生错误: %v. 详情: %s\n", err, tidyStderr.String())
	}

	return nil
}
