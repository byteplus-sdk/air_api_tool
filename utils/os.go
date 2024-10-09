package utils

import "os"

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func EnsureDirExists(dirPath string) error {
	// 检查目录是否存在
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// 如果目录不存在，尝试创建
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
