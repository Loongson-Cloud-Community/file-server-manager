package main

import (
	"os"
  "errors"
)

// WriteFile 写文件 本地文件系统
// 如果文件已存在，覆盖
// 如果文件不存在，创建
func CreateFile(file string, bytes []byte) error {
	var err error
	var f *os.File

	if f, err = os.Create(file); err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(bytes); err != nil {
		return err
	}
	return nil
}

// 删除文件
// 如果文件不存在，返回 nil
func RemoveFile(file string) error {
	err := os.Remove(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
