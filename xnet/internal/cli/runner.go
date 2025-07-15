package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
)

var (
	green = color.New(color.FgGreen)
	blue  = color.New(color.FgBlue)
	cyan  = color.New(color.FgCyan)
)

// Runner 存储选项参数
type Runner struct {
	Proxy string
	Size  uint64
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.Size == 0 {
		return fmt.Errorf("参数 --size 必须大于 0")
	}

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	const bytesInMiB uint64 = 1024 * 1024

	// 将 MiB 转换为 Bytes
	bytesToDownload := r.Size * bytesInMiB

	downloadURL := fmt.Sprintf("https://speed.cloudflare.com/__down?bytes=%d", bytesToDownload)

	fmt.Printf("Speed Test Configuration:\n")
	fmt.Printf("  Download URL: %s\n", downloadURL)
	fmt.Printf("  Download Size: ")
	cyan.Printf("%d MiB (%d bytes)\n", r.Size, bytesToDownload)

	useProxy := r.Proxy != ""
	if useProxy {
		fmt.Printf("  Using Proxy: %s\n", r.Proxy)
	} else {
		fmt.Printf("  Proxy: Disabled\n")
	}

	blue.Println("------ Start ------")

	var outputDevice string
	if runtime.GOOS == "windows" {
		outputDevice = "NUL"
	} else {
		outputDevice = "/dev/null"
	}

	cmdArgs := []string{"-O", outputDevice, downloadURL}
	wgetCmd := exec.Command("wget", cmdArgs...)

	currentEnv := os.Environ()
	if useProxy {
		currentEnv = append(currentEnv, "http_proxy="+r.Proxy)
		currentEnv = append(currentEnv, "https_proxy="+r.Proxy)
	}
	wgetCmd.Env = currentEnv

	wgetCmd.Stdout = os.Stdout
	wgetCmd.Stderr = os.Stderr

	fmt.Printf("Starting download with wget...\n")

	startTime := time.Now()
	err := wgetCmd.Run()
	endTime := time.Now()

	blue.Println("------ Complete ------")

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("wget command failed with exit code %d.\n", exitError.ExitCode())
		} else {
			return fmt.Errorf("Error executing wget: %w", err)
		}
	}

	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()

	fmt.Printf("Download finished.\n")
	fmt.Printf("Time taken: ")
	green.Printf("%.2f seconds\n", seconds)

	speedBytesPerSec := float64(bytesToDownload) / seconds
	speedMbps := (speedBytesPerSec * 8) / (1000 * 1000)
	fmt.Printf("Average download speed: ")
	green.Printf("%.2f Mbps\n", speedMbps)

	return nil
}
