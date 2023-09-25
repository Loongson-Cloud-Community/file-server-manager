package main

import (
	"errors"
	"io"
	"log"
	"os"
	"path"
)

/* test/
   - file.txt -> hash

   test/
   - file.txt
   - file.txt.md5 -> hash
*/

var HashType = "md5"

func IncSyncToFileServer(op *Operation) error {
	if !op.Valid {
		return nil
	}
	switch op.Type {
	case CREATE:
		if err := AddToFileServer(op); err != nil {
			return err
		}
		log.Printf("[file-server]:\tCREATE: %s/%s %s", op.Directory, op.File, op.Hash)
	case DELETE:
		if err := DeleteFromFileServer(op); err != nil {
			return err
		}
		log.Printf("[file-server]:\tDELETE: %s/%s %s", op.Directory, op.File, op.Hash)
	}
	return nil
}

func DeleteFromFileServer(op *Operation) error {
	if err := deleteFile(op); err != nil {
		return err
	}
	if err := deleteHashFile(op); err != nil {
		return err
	}
	return nil
}

func AddToFileServer(op *Operation) error {
	// 如果目录不存在，创建目录
	dir := path.Join(FileServerDir, op.Directory)
	stat, err := os.Stat(dir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if err != nil && !stat.IsDir() {
		return errors.New("FileSources: Path already exist and is not directory: " + dir)
	}
	if err := addFile(op); err != nil {
		return err
	}
	if err := addHashFile(op); err != nil {
		return err
	}
	return nil
}

func addFile(op *Operation) error {
	file := path.Join(FileServerDir, op.Directory, op.File)
	old := path.Join(FileHistoryDir, op.Hash)
	if err := CopyFile(old, file); err != nil {
		return nil
	}
	return nil
}

func deleteFile(op *Operation) error {
	file := path.Join(FileServerDir, op.Directory, op.File)
	if err := RemoveFile(file); err != nil {
		return err
	}
	return nil
}

func CopyFile(oldpath string, newpath string) error {
	of, err := os.Open(oldpath)
	if err != nil {
		return err
	}
	defer of.Close()

	nf, err := os.Create(newpath)
	if err != nil {
		return err
	}
	defer nf.Close()

	io.Copy(nf, of)
	return nil
}

func addHashFile(op *Operation) error {
	file := path.Join(FileServerDir, op.Directory, op.File+"."+HashType)
	if err := CreateFile(file, []byte(op.Hash)); err != nil {
		return err
	}
	return nil
}

func deleteHashFile(op *Operation) error {
	file := path.Join(FileServerDir, op.Directory, op.File+"."+HashType)
	if err := RemoveFile(file); err != nil {
		return err
	}
	return nil
}
