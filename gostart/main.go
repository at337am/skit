package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "💡 usage: gostart <projectName>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	projectDir := args[0]

	fmt.Printf("⚙️  项目 [%s] 生成中...\n", projectDir)

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建项目目录失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 .gitignore
	gitignoreContent := `# https://github.com/github/gitignore/blob/main/Go.gitignore
*.exe
*.exe~
*.dll
*.so
*.dylib

*.test

*.out

go.work
go.work.sum

.env

# mjj

# directory
data/
tmp/
images/
videos/

# archive compressed
*.zip
*.rar
*.7z

*.tar
*.tgz
*.tar.gz
*.tar.xz
*.tar.bz2
*.tar.zst
*.tar.lzma
*.tar.lz
*.tar.lz4

*.iso

# audio file
*.mp3
*.Mp3
*.MP3

*.flac
*.Flac
*.FLAC

*.aac
*.Aac
*.AAC

*.wav
*.Wav
*.WAV

*.wma
*.Wma
*.WMA

*.ogg
*.Ogg
*.OGG

*.alac
*.Alac
*.ALAC
*.aiff
*.Aiff
*.AIFF

# video file
*.mp4
*.Mp4
*.MP4

*.mov
*.Mov
*.MOV

*.mkv
*.Mkv
*.MKV

*.avi
*.Avi
*.AVI

*.webm
*.WebM
*.WEBM

# image file
*.jpg
*.Jpg
*.JPG

*.jpeg
*.Jpeg
*.JPEG

*.png
*.Png
*.PNG

*.bmp
*.Bmp
*.BMP

*.tif
*.Tif
*.TIF

*.tiff
*.Tiff
*.TIFF

*.webp
*.Webp
*.WEBP

*.heif
*.Heif
*.HEIF

*.heic
*.Heic
*.HEIC

*.svg
*.Svg
*.SVG

*.raw
*.Raw
*.RAW

*.cr2
*.Cr2
*.CR2

*.nef
*.Nef
*.NEF

*.gif
*.Gif
*.GIF

# misc
.out
`
	gitignorePath := filepath.Join(projectDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 创建 .gitignore 文件失败: %v\n", err)
		os.Exit(1)
	}

	// 创建 main.go
	mainGoContent := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
`
	mainGoPath := filepath.Join(projectDir, "main.go")
	if err := os.WriteFile(mainGoPath, []byte(mainGoContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 创建 main.go 文件失败: %v\n", err)
		os.Exit(1)
	}

	// 获取项目的绝对路径用于显示
	absPath, err := filepath.Abs(projectDir)
	if err != nil {
		absPath = projectDir // 如果获取绝对路径失败，使用相对路径
	}

	// 执行 go mod init 命令
	cmd := exec.Command("go", "mod", "init", projectDir)
	cmd.Dir = projectDir // 设置工作目录为项目目录

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "❌ 初始化 go mod 失败: %v\n%s", err, stderr.String())
		os.Exit(1)
	}

	fmt.Printf("\n🎉 项目 [%s] 创建成功! 🎉\n", projectDir)
	fmt.Printf("📂 项目路径: %s\n", absPath)
	fmt.Printf("📋 项目结构:\n")
	fmt.Printf("  ├─ 📄 .gitignore\n")
	fmt.Printf("  ├─ 📄 main.go\n")
	fmt.Printf("  └─ 📄 go.mod\n")
	fmt.Printf("\n🚀 开始吧!\n")

}
