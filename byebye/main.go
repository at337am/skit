package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// --- 常量配置 ---
const (
	// 等待进程响应 SIGTERM 的总时间
	gracefulShutdownTimeout = 500 * time.Millisecond
	// 需要终止的进程名列表 (逗号分隔)
	processNamesToKill = "nekoray,telegram,crow"
)

// findPIDs 使用 pgrep 查找指定名称的进程 ID
func findPIDs(processName string) ([]string, error) {
	// 使用 pgrep -f 可以匹配完整命令行，可能更精确，但这里保持和之前一致
	// cmd := exec.Command("pgrep", "-f", processName)
	cmd := exec.Command("pgrep", processName)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok && status.ExitStatus() == 1 {
				log.Printf("未找到进程 '%s'\n", processName) // 静默处理未找到的情况
				return []string{}, nil // 未找到不是错误
			}
		}
		log.Printf("执行 pgrep %s 时出错: %v, Stderr: %s\n", processName, err, stderr.String())
		return nil, fmt.Errorf("查找进程 %s 失败: %w", processName, err)
	}

	output := strings.TrimSpace(stdout.String())
	if output == "" {
		log.Printf("未找到正在运行的进程 '%s'\n", processName) // 静默处理未找到的情况
		return []string{}, nil
	}

	pids := strings.Split(output, "\n")
	// log.Printf("找到进程 '%s' 的 PIDs: %v\n", processName, pids) // 在主函数统一打印
	return pids, nil
}

// sendSignalToPID 向单个 PID 发送信号
func sendSignalToPID(pidStr string, signal syscall.Signal) error {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Printf("无效的 PID 字符串 '%s': %v\n", pidStr, err)
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("os.FindProcess 对于 PID %d 失败 (通常不应发生): %v\n", pid, err)
		return err
	}

	err = process.Signal(signal)
	if err != nil {
		// 检查是否是进程已退出的错误
		if errors.Is(err, os.ErrProcessDone) || (strings.Contains(err.Error(), "process already finished")) || (strings.Contains(err.Error(), "no such process")) {
            // log.Printf("发送信号 %v 到 PID %d 时发现进程已退出。\n", signal, pid)
			return nil // 进程已退出，不算发送失败
		}
		log.Printf("向 PID %d 发送信号 %v 失败: %v\n", pid, signal, err)
		return err
	}
	// log.Printf("已向 PID %d 发送信号 %v。\n", pid, signal)
	return nil
}

// isProcessRunning 检查 PID 是否仍在运行
func isProcessRunning(pidStr string) (bool, error) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Printf("无效的 PID 字符串 '%s': %v\n", pidStr, err)
		return false, err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("os.FindProcess 对于 PID %d 失败: %v\n", pid, err)
		return false, err
	}

	// 发送 Signal 0 检查进程是否存在且有权限操作
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil // 进程存在且可操作
	}
	if errors.Is(err, os.ErrProcessDone) || strings.Contains(err.Error(), "no such process") {
		return false, nil // 进程不存在
	}
    // 其他错误，可能是权限问题 EPERM 等
	log.Printf("检查 PID %d 状态时出错: %v\n", pid, err)
    // 对于权限错误等，我们可能仍认为它“在运行”（我们只是无法确认它已停止）
    // 或者可以根据具体错误类型决定，这里保守地返回错误
	return false, err
}


// shutdownSystem 执行 sudo shutdown -h now
func shutdownSystem() {
	log.Println("准备执行 sudo shutdown -h now...")
	cmd := exec.Command("sudo", "shutdown", "-h", "now")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out // 合并 stdout 和 stderr

	err := cmd.Run()
	if err != nil {
		// 使用 log.Printf 而不是 log.Fatalf，以便脚本可以继续执行到最后打印日志
		log.Printf("执行 'sudo shutdown -h now' 失败: %v\n输出/错误: %s\n请检查 sudo 权限是否配置正确 (可能需要无密码执行 shutdown)。", err, out.String())
	} else {
		log.Println("关机命令已成功发送。系统应该很快会关闭。")
		// 注意：关机命令发送成功后，脚本本身可能很快被终止
	}
}

func main() {
	targetNames := strings.Split(processNamesToKill, ",")
	allPIDs := []string{}

	log.Println("开始查找需要终止的进程...")
	foundAny := false
	for _, name := range targetNames {
		trimmedName := strings.TrimSpace(name)
		if trimmedName == "" {
			continue
		}
		pids, err := findPIDs(trimmedName)
		if err != nil {
			// findPIDs 内部已记录错误
			continue
		}
		if len(pids) > 0 {
			log.Printf("找到进程 '%s' 的 PIDs: %v\n", trimmedName, pids)
			allPIDs = append(allPIDs, pids...)
			foundAny = true
		}
	}

	if !foundAny {
		log.Println("没有找到任何目标进程。")
	} else {
		// --- 阶段 1: 发送 SIGTERM ---
		log.Printf("向所有找到的 PIDs 发送 SIGTERM...\n")
		for _, pidStr := range allPIDs {
			_ = sendSignalToPID(pidStr, syscall.SIGTERM) // 忽略发送错误，继续尝试下一个
		}

		// --- 阶段 2: 等待 ---
		log.Printf("等待 %v 给所有进程优雅退出的时间...\n", gracefulShutdownTimeout)
		time.Sleep(gracefulShutdownTimeout)

		// --- 阶段 3: 检查并发送 SIGKILL ---
		remainingPIDs := []string{}
		log.Println("检查哪些进程仍在运行...")
		for _, pidStr := range allPIDs {
			running, err := isProcessRunning(pidStr)
			if err != nil {
				// isProcessRunning 内部已记录错误
                // 如果检查出错，保守起见可能也需要尝试 kill？或者跳过？这里选择跳过
                log.Printf("检查 PID %s 状态出错，跳过后续处理。\n", pidStr)
				continue
			}
			if running {
				remainingPIDs = append(remainingPIDs, pidStr)
			}
		}

		if len(remainingPIDs) > 0 {
			log.Printf("以下 PIDs 仍在运行，将发送 SIGKILL: %v\n", remainingPIDs)
			for _, pidStr := range remainingPIDs {
				_ = sendSignalToPID(pidStr, syscall.SIGKILL) // 忽略发送错误
			}
			log.Println("SIGKILL 已发送给剩余进程。")
		} else {
			log.Println("所有目标进程似乎都已成功退出。")
		}
	}

	// --- 阶段 4: 执行关机 ---
	// fmt.Println("✅ 关机成功...")
	shutdownSystem()
}
