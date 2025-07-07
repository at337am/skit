package converter

// IConverter 定义了转换的行为
type IConverter interface {
	ConvertToMP4(inputPath, outputDir string) (*ConvertResult, error)
}

type Converter struct{}

func NewConverter() IConverter {
	return &Converter{}
}
