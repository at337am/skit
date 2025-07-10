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
	// 在这里实现校验参数的逻辑

	// if r.Path == "" {
	// 	return fmt.Errorf("The path is empty -> '%s'", r.Path)
	// }

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {
	// 在这里实现核心逻辑

	// funcName()
	// askForConfirmation("Hello")

	return nil
}
