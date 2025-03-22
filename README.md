# skit - Some scripts and toolkits

## 脚本介绍

##### `aout` 从视频文件提取音频，支持单个文件或整个目录

```bash
aout -v 01.mp4

aout -d . -e mp3
```

##### `clipvid` 视频切割

```go
func main() {
	inputFile := "data/4.mkv"
	outputFile := "output.mkv"
	startTime := "00:25:56"
	endTime := "00:26:00"

	err := clipVideo(inputFile, outputFile, startTime, endTime)
	if err != nil {
		fmt.Println("❌ ", err)
	}
}
```

##### `dirhash` 用于比较两个文件或目录的哈希值，以检查它们是否完全一致

```bash
dirhash 01.txt 02.txt

dirhash dir1/ dir2/
```

##### `fmp4` 格式化 .mp4 后缀名, 删除所有 .mov 视频

```bash
fmp4 ./
```

##### `gostart` 脚手架, 生成一个简易 Go 项目

```bash
gostart skit
```

##### `img2vid` 将多张图片合成视频

```bash
img2vid -d .

img2vid -d ./images -s 0.5
```

##### `openurl` 快速打开多个 URL

```bash
cd openurl
./main
```

```yaml
sites:
  - name: "Google"
    url: "https://www.google.com"
  - name: "GitHub"
    url: "https://www.github.com"
```

##### `repaudio` 给视频添加或更换音频

```bash
repaudio -v 01.mp4 -a 02.mp3
```

##### `tmrn` 格式化 filename 为 %02d

```bash
tmrn -d .

tmrn -d . -e png
```

##### `vcat` 横向拼接两个视频，并可选择性地添加音频

```bash
vcat -v1 01.mp4 -v2 02.mp4

vcat -v1 01.mp4 -v2 02.mp4 -a 01.aac -o result.mp4
```


##### `vid2img` 提取视频所有帧保存为 .png

```bash
vid2img ./01.mp4
```

##### `vid2mp4` 视频转换成 .mp4, 支持单个文件和目录下 mov 批量处理

```bash
vid2mp4 ./

vid2mp4 ./01.mov
```

##### `xfixer` 检查文件扩展名是否正确 (支持单个文件和目录)

```bash
xfixer ./
```


## 视频批量清洗

1. 对路径中所有文件进行清洗, 检查文件扩展名是否正确  

```bash
xfixer ./
```

2. 将所有 mov 转换成 mp4

```bash
vid2mp4 ./
```

3. 统一后缀名 .mp4 为小写, 并删除多余的 mov  

```bash
fmp4 ./
```

4. 最后按照修改日期顺序重命名所有文件  

```bash
tmrn -d ./
```
