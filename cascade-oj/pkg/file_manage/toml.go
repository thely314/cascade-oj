package filemanage

import (
	"os"
	"path/filepath"
	"strconv"
)

// returns the content of the toml file
func ReadTomlFile(problemId int64, basePath string) ([]byte, error) {
	tomlFile, err := os.ReadFile(basePath + "/" + strconv.FormatInt(problemId, 10) + "/testcase/config.toml")
	if err != nil {
		return nil, err
	}
	return tomlFile, nil
}

func WriteTomlFile(problemId int64, basePath string, content []byte) error {
	filePath := basePath + "/" + strconv.FormatInt(problemId, 10) + "/testcase/config.toml"
	dirPath := filepath.Dir(filePath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, content, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
