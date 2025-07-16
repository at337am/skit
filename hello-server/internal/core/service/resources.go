package service

import (
	"os"
	"path"
	"path/filepath"
)

// Resource 资源结构体
type Resource struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	URL      string     `json:"url,omitempty"`
	Children []Resource `json:"children,omitempty"`
}

// GetDirectoryStructure 递归获取目录结构
func GetDirectoryStructure(rootPath, baseURL string) ([]Resource, error) {
	var resources []Resource
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		resource := Resource{Name: entry.Name()}
		fullPath := filepath.Join(rootPath, entry.Name())

		if entry.IsDir() {
			resource.Type = "directory"
			newBaseURL := path.Join(baseURL, entry.Name())
			children, err := GetDirectoryStructure(fullPath, newBaseURL)
			if err != nil {
				return nil, err
			}
			resource.Children = children
		} else {
			resource.Type = "file"
			resource.URL = path.Join(baseURL, entry.Name())
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
