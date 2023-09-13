package main

import (
	"os"
)

// WriteFile 写文件 本地文件系统
// 如果文件已存在，则覆盖
func WriteFile(file string, bytes []byte) error {
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
