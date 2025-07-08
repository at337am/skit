package converter

type ConvertResult struct {
	OutputPath    string // 最终输出的文件路径
	StatusMessage string // 描述转换过程中的关键信息, 如音频是否转码
}

// Converter 定义了转换的行为
type Converter interface {
	ConvertToMP4(inputPath, outputDir string) (*ConvertResult, error)
}
