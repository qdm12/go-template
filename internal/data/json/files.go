package json

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrFileIsDirectory = errors.New("file is a directory")
)

func fileExists(filePath string) (exists bool, err error) {
	stat, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("stat-ing file: %w", err)
	}

	if stat.IsDir() {
		return false, fmt.Errorf("%w: %s", ErrFileIsDirectory, filePath)
	}

	return true, nil
}
