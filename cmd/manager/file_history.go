package main

import (
	"errors"
	"os"
	"path"
)

func StoreToFileHistory(op Operation) error {
	if op.Type != CREATE {
		return nil
	}
	file := path.Join(FileHistoryDir, op.Key)
	fileInfo, err := os.Stat(file)
	if err == nil && fileInfo.IsDir() {
		return errors.New("FileHistory: File already exist and is directory: " + file)
	}
	if err = CreateFile(file, op.Bytes); err != nil {
		return err
	}
	return nil
}
