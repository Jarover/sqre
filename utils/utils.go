package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Fail to get local dir")
		return ""
	}

	return dir
}

