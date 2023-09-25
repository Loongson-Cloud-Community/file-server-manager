package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

var ErrInvalidOperation error = errors.New("Invalid operation")

func CommitToFileSources(op *Operation) error {
	switch op.Type {
	case CREATE:
		if err := AddToFileSources(op); err != nil {
			return err
		}
	case DELETE:
		if err := DeleteFromFileSources(op); err != nil {
			return err
		}
	}
	if !op.Valid {
		return nil
	}
	// Git Commit
	var cmd *exec.Cmd
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = FileSourcesDir
	if err := cmd.Run(); err != nil {
		return err
	}

	msg := GenerateCommitMessage(op)
	cmd = exec.Command("git", "commit", "-m", msg)
	cmd.Dir = FileSourcesDir
	if err := cmd.Run(); err != nil {
		return err
	}
	log.Printf("[file-sources]:\t%s", strings.Replace(msg, "\n", " ", -1))
	return nil
}

func DeleteFromFileSources(op *Operation) error {
	file := path.Join(FileSourcesDir, op.Directory, op.File)
	// 判断有效删除
	_, err := os.Stat(file)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		op.Valid = false
		return nil
	}
	// 读取删除文件 hash
	f, _ := os.Open(file)
	b, _ := io.ReadAll(f)
	op.Hash = string(b)
	if err := RemoveFile(file); err != nil {
		return err
	}
	return nil
}

func AddToFileSources(op *Operation) error {
	dir := path.Join(FileSourcesDir, op.Directory)
	stat, err := os.Stat(dir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if err != nil && !stat.IsDir() {
		return errors.New("FileSources: Path already exist and is not directory: " + dir)
	}
	file := path.Join(dir, op.File)
	stat, err = os.Stat(file)
	// 文件不存在
	if err != nil && errors.Is(err, os.ErrNotExist) {
		CreateFile(file, []byte(op.Hash))
		return nil
	}
	// 文件存在
	if err == nil {
		if stat.IsDir() {
			return errors.New("FileSources: File already exist and is directory: " + file)
		}
		f, _ := os.Open(file)
		defer f.Close()
		b, _ := io.ReadAll(f)
		if string(b) == op.Hash {
			op.Valid = false
		}
	}
	return nil
}

// CREATE: dir/file
// $HASH
func GenerateCommitMessage(op *Operation) string {
	var t string
	switch op.Type {
	case CREATE:
		t = "CREATE"
	case DELETE:
		t = "DELETE"
	}
	return fmt.Sprintf("%s: %s/%s\n%s", t, op.Directory, op.File, op.Hash)
}
