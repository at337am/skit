package cli

import (
	"fmt"
	"sort"
)

// diffResult 存储两个哈希图的比较结果
type diffResult struct {
	modified []string // 相同路径但哈希值不同的文件
	onlyIn1  []string // 只存在于第一个路径中的文件
	onlyIn2  []string // 只存在于第二个路径中的文件
}

// diff 比较两个哈希图并返回差异
func diff(map1, map2 map[string]string) *diffResult {
	result := &diffResult{}

	// 遍历 map1 找出修改的和只在 map1 中的文件
	for path, hash1 := range map1 {
		hash2, ok := map2[path]
		if !ok {
			// 文件只在 map1 中
			result.onlyIn1 = append(result.onlyIn1, path)
		} else if hash1 != hash2 {
			// 文件被修改了
			result.modified = append(result.modified, path)
		}
	}

	// 遍历 map2 找出只在 map2 中的文件
	for path := range map2 {
		if _, ok := map1[path]; !ok {
			result.onlyIn2 = append(result.onlyIn2, path)
		}
	}

	// 为了输出稳定和美观, 对结果进行排序
	sort.Strings(result.modified)
	sort.Strings(result.onlyIn1)
	sort.Strings(result.onlyIn2)

	return result
}

func (r *Runner) compareDir() error {
	path1 := r.Path1
	path2 := r.Path2

	// 为第一个路径生成哈希图
	map1, err := r.hash.HashDir(path1)
	if err != nil {
		return fmt.Errorf("路径: '%s' 计算哈希时出错: %w", path1, err)
	}

	// 为第二个路径生成哈希图
	map2, err := r.hash.HashDir(path2)
	if err != nil {
		return fmt.Errorf("路径: '%s' 计算哈希时出错: %w", path2, err)
	}

	fmt.Printf("%s -> %d 个文件\n", path1, len(map1))
	fmt.Printf("%s -> %d 个文件\n", path2, len(map2))

	// 比较两个哈希图
	diffs := diff(map1, map2)

	if len(diffs.modified) == 0 && len(diffs.onlyIn1) == 0 && len(diffs.onlyIn2) == 0 {
		sameColor.Printf("\n两个路径完全一致!\n")
		return nil
	}

	diffColor.Printf("\n两个路径存在差异!\n")

	if len(diffs.modified) > 0 {
		diffColor.Printf("\n-> 哈希不一致的文件:\n")
		for _, file := range diffs.modified {
			fmt.Println(file)
		}
	}

	if len(diffs.onlyIn1) > 0 {
		diffColor.Printf("\n-> 仅存在于 '%s' 的文件:\n", path1)
		for _, file := range diffs.onlyIn1 {
			fmt.Println(file)
		}
	}

	if len(diffs.onlyIn2) > 0 {
		diffColor.Printf("\n-> 仅存在于 '%s' 的文件:\n", path2)
		for _, file := range diffs.onlyIn2 {
			fmt.Println(file)
		}
	}

	return nil
}
