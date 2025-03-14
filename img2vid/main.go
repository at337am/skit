package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	outputVid = "output.mp4" // 生成的视频文件
	duration = 0.6 // 每张图片显示的时间（秒）
	targetHeight = 2160 // 统一的高度
)

func getMaxWidthAndResize(images []string, imageDir string) (int, error) {
	maxWidth := 0
	for _, img := range images {
		filePath := filepath.Join(imageDir, img)
		file, err := os.Open(filePath)
		if err != nil {
			return 0, err
		}
		imgConfig, _, err := image.DecodeConfig(file)
		file.Close()
		if err != nil {
			return 0, err
		}
		width := imgConfig.Width * targetHeight / imgConfig.Height
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth, nil
}

func main() {

	imageDir := flag.String("d", "", "PNG 图片所在的目录")
	output := flag.String("o", outputVid, "输出视频文件名")
	durationFlag := flag.Float64("s", duration, "每张图片显示的时间（秒）")
	height := flag.Int("height", targetHeight, "输出视频高度")
	
	flag.Parse()

	if *imageDir == "" {
		fmt.Println("❌ 必须指定 PNG 图片所在的目录 (-d 参数)!")
		fmt.Println("💡 用法: go run main.go -d images")
		return
	}
	
	files, err := os.ReadDir(*imageDir)
	if err != nil {
		fmt.Println("无法读取图片目录:", err)
		return
	}
	
	var images []string
	for _, file := range files {
		if !file.IsDir() && len(file.Name()) > 4 && file.Name()[len(file.Name())-4:] == ".png" {
			images = append(images, file.Name())
		}
	}
	
	if len(images) == 0 {
		fmt.Println("未找到 PNG 图片")
		return
	}
	
	sort.Slice(images, func(i, j int) bool {
		ni, _ := strconv.Atoi(images[i][:len(images[i])-4])
		nj, _ := strconv.Atoi(images[j][:len(images[j])-4])
		return ni < nj
	})
	
	maxWidth, err := getMaxWidthAndResize(images, *imageDir)
	if err != nil {
		fmt.Println("获取最大宽度失败:", err)
		return
	}
	
	tempList := "file_list.txt"

	defer func() {
		if err := os.Remove(tempList); err != nil {
			fmt.Println("删除临时文件失败:", err)
		} else {
			fmt.Println("临时文件已删除:", tempList)
		}
	}()
	
	fileList, err := os.Create(tempList)
	if err != nil {
		fmt.Println("无法创建文件列表:", err)
		return
	}
	defer fileList.Close()
	
	for _, img := range images {
		fileList.WriteString(fmt.Sprintf("file '%s/%s'\nduration %.1f\n", *imageDir, img, *durationFlag))
	}
	
	cmd := exec.Command("ffmpeg", 
		"-y",
		"-f", "concat", 
		"-safe", "0", 
		"-i", tempList, 
		"-vf", fmt.Sprintf("scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", 
			maxWidth, *height, maxWidth, *height),
		"-c:v", "libx264", 
		"-crf", "0", 
		*output,
	)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		fmt.Println("FFmpeg 执行失败:", err)
		return
	}
	
	fmt.Println("视频生成成功:", *output)
}
