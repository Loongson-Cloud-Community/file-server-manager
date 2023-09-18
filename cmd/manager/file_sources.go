package main

import (
	"errors"
	"os/exec"
	"path"
  "os"
  "fmt"
)

func CommitToFileSources(op Operation) error {
	switch op.Type {
	case CREATE:
		if err := AddFileToFileSources(op); err != nil {
			return err
		}
	case DELETE:
		if err := DeleteFileFromFileSources(op); err != nil {
			return err
		}
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
	return nil
}

func DeleteFileFromFileSources(op Operation) error {
	file := path.Join(FileSourcesDir, op.Directory, op.File)
	if err := RemoveFile(file); err != nil {
		return err
	}
	return nil
}

func AddFileToFileSources(op Operation) error {
	dir := path.Join(FileSourcesDir, op.Directory)
	fileInfo, err := os.Stat(dir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	if err != nil && !fileInfo.IsDir() {
		return errors.New("FileSources: Path already exist and is not directory: " + dir)
	}
	file := path.Join(dir, op.File)
	fileInfo, err = os.Stat(file)
	if err == nil && fileInfo.IsDir() {
		return errors.New("FileSources: File already exist and is directory: " + file)
	}
	if err = CreateFile(file, []byte(op.Key)); err != nil {
		return err
	}
	return nil

}

// CREATE: dir/file
func GenerateCommitMessage(op Operation) string {
	var t string
	switch op.Type {
	case CREATE:
		t = "CREATE"
	case DELETE:
		t = "DELETE"
	}
	return fmt.Sprintf("%s: %s/%s", t, op.Directory, op.File)
}
