package nacos

import (
	"os"
	"path/filepath"
)

type FileSource struct {
	BaseDir string
}

func NewFileSource(baseDir string) *FileSource {
	return &FileSource{BaseDir: baseDir}
}

func (f *FileSource) Source() Source {
	return func(group string, dataId string) (string, error) {
		var path string
		if group == "" {
			path = filepath.Join(f.BaseDir, dataId)
		} else {
			path = filepath.Join(f.BaseDir, group, dataId)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
}
