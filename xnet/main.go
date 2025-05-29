package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	bytesInMB = 1024 * 1024
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
)

func customUsage (){
	fmt.Println("用法: xnet [可选]")
	fmt.Println("选项:")
	fmt.Println("  -s int")
	fmt.Println("      指定测试数据大小，单位为 MB（例如：100 表示 100MB），默认值为 50MB")
	fmt.Println("  -p string")
	fmt.Println("      指定代理地址，默认值为 http://127.0.0.1:2080")
	fmt.Println("      若不使用代理，请传入空字符串")
}

func main() {
	sizeMB := flag.Int64("s", 50, "下载大小，单位是 MB")
	proxyURL := flag.String("p", "http://127.0.0.1:2080", "HTTP/HTTPS 代理地址")

	flag.Usage = customUsage

	flag.Parse()

	if *sizeMB <= 0 {
		fmt.Fprintf(os.Stderr, "❌ 下载大小 (-s) 必须大于 0 MB\n")
		os.Exit(1)
	}

	// 将 MB 转换为 Bytes
	bytesToDownload := *sizeMB * bytesInMB

	// 检查 wget 是否安装
	if _, err := exec.LookPath("wget"); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 未找到 wget 命令\n")
		os.Exit(1)
	}

	downloadURL := fmt.Sprintf("https://speed.cloudflare.com/__down?bytes=%d", bytesToDownload)

	fmt.Printf("Speed Test Configuration:\n")
	fmt.Printf("  Download URL: %s\n", downloadURL)
	fmt.Printf("  Download Size: %s%d MiB%s (%d bytes)\n", colorGreen, *sizeMB, colorReset, bytesToDownload)

	useProxy := *proxyURL != ""
	if useProxy {
		fmt.Printf("  Using Proxy: %s\n", *proxyURL)
	} else {
		fmt.Printf("  Proxy: Disabled\n")
	}

	fmt.Println(colorBlue + "--------------------------------------------------" + colorReset)

	var outputDevice string
	if runtime.GOOS == "windows" {
		outputDevice = "NUL"
	} else {
		outputDevice = "/dev/null"
	}

	cmdArgs := []string{"-O", outputDevice, downloadURL}
	cmd := exec.Command("wget", cmdArgs...)

	currentEnv := os.Environ()
	if useProxy {
		currentEnv = append(currentEnv, "http_proxy="+*proxyURL)
		currentEnv = append(currentEnv, "https_proxy="+*proxyURL)
	}
	cmd.Env = currentEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Starting download with wget...")

	startTime := time.Now()
	err := cmd.Run()
	endTime := time.Now()

	fmt.Println(colorBlue + "--------------------------------------------------" + colorReset)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "❌ wget command failed with exit code %d.\n", exitError.ExitCode())
		} else {
			fmt.Fprintf(os.Stderr, "❌ Error executing wget: %v\n", err)
		}
		os.Exit(1)
	}

	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()

	fmt.Printf("Download finished.\n")
	fmt.Printf("Time taken: %s%.2f seconds%s\n", colorGreen, seconds, colorReset)

	speedBytesPerSec := float64(bytesToDownload) / seconds
	speedMbps := (speedBytesPerSec * 8) / (1000 * 1000)
	fmt.Printf("Average download speed: %s%.2f Mbps%s\n", colorGreen, speedMbps, colorReset)
}
