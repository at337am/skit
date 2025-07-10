package action

import "fmt"

// NewRunner 构造函数 (也可以在这里设置参数默认值)
func NewRunner() *Runner {
	return &Runner{}
}

// Validate 校验参数
func (r *Runner) Validate() error {
	if r.Path == "" {
		return fmt.Errorf("The path is empty -> '%s'", r.Path)
	}

	// 在这里实现校验参数的逻辑

	return nil
}

// Run 执行核心逻辑
func (r *Runner) Run() error {

	// 在这里实现核心逻辑
	funcName()

	return nil
}
