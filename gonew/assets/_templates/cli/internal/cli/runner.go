package cli

// Runner 存储选项参数
type Runner struct {
	// Path    string
	// Message string
	// Port    int
	// Yes     bool
}

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {

	// if r.Path == "" {
	// 	return errors.New("path is empty")
	// }

	// if _, err := os.Stat(r.Path); err != nil {
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		return fmt.Errorf("path does not exist: %s", r.Path)
	// 	}
	// 	return fmt.Errorf("could not access path %s: %w", r.Path, err)
	// }

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {

	// funcName()
	// askForConfirmation("Hello")

	return nil
}
