package differ

import "sort"

// DiffResult 存储两个哈希图的比较结果
type DiffResult struct {
	// 相同路径但哈希值不同的文件
	Modified []string
	// 只存在于第一个路径中的文件
	OnlyIn1 []string
	// 只存在于第二个路径中的文件
	OnlyIn2 []string
}

// IsEmpty 检查是否有任何差异
func (r *DiffResult) IsEmpty() bool {
	return len(r.Modified) == 0 && len(r.OnlyIn1) == 0 && len(r.OnlyIn2) == 0
}

// Compare 比较两个哈希图并返回差异
func Compare(map1, map2 map[string]string) DiffResult {
	var result DiffResult

	// 遍历 map1 找出修改的和只在 map1 中的文件
	for path, hash1 := range map1 {
		hash2, ok := map2[path]
		if !ok {
			// 文件只在 map1 中
			result.OnlyIn1 = append(result.OnlyIn1, path)
		} else if hash1 != hash2 {
			// 文件被修改了
			result.Modified = append(result.Modified, path)
		}
	}

	// 遍历 map2 找出只在 map2 中的文件
	for path := range map2 {
		if _, ok := map1[path]; !ok {
			result.OnlyIn2 = append(result.OnlyIn2, path)
		}
	}

	// 为了输出稳定和美观，对结果进行排序
	sort.Strings(result.Modified)
	sort.Strings(result.OnlyIn1)
	sort.Strings(result.OnlyIn2)

	return result
}
