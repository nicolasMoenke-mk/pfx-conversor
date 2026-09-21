package utils

import (
	"fmt"
	"os"
)

func ReplaceNumber(file []byte) ([]byte, error) {

	file, err := os.ReadFile(string(file))
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return file, nil
}
