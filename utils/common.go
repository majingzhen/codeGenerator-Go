package utils

import (
	"os"
)

func GetProjectPath() string {
	// 获取当前的工作目录
	projectPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return projectPath
}
