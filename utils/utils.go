package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func FloatToString(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 6, 64)
}

func GetDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Fail to get local dir")
		return ""
	}

	return dir
}

func GetBaseFile() string {
	filename := os.Args[0] // get command line first parameter
	return strings.Split(filepath.Base(filename), ".")[0]
}
