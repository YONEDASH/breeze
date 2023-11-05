package common

import (
	"os"
	"path/filepath"
)

type Position struct {
	Index  int
	Line   int
	Column int
}

func InitPosition() Position {
	return Position{Line: 1, Column: 1, Index: 0}
}

type SourceFile struct {
	Path string
}

func (sf *SourceFile) Validate() error {
	absPath, err := filepath.Abs(sf.Path)
	if err != nil {
		return err
	}
	sf.Path = absPath
	return nil
}

func (sf *SourceFile) GetContent() (string, error) {
	content, err := os.ReadFile(sf.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func InitSource(path string) SourceFile {
	return SourceFile{Path: path}
}
