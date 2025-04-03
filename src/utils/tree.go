package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	branch string = "│   "
	tee    string = "├── "
	last   string = "└── "

	folder    string = "📁 "
	nonFolder string = "📝 "
)

func GenerateTree(path string, depth int, distance int, indent string, result *string) error {
	if distance > depth {
		return nil
	}

	dir, err := os.Open(path)
	if err != nil {
		return err
	}

	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for i, file := range files {
		var next string
		if i == len(files)-1 {
			next = last
		} else {
			next = tee
		}

		if file.IsDir() {
			*result += fmt.Sprintf("%s%s%s%s\n", indent, next, folder, file.Name())
		} else {
			*result += fmt.Sprintf("%s%s%s%s\n", indent, next, nonFolder, file.Name())
		}

		if file.IsDir() && distance < depth {
			newIndent := indent + branch
			err := GenerateTree(filepath.Join(path, file.Name()), depth, distance+1, newIndent, result)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
