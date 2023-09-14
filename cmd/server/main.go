package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
)

const (
	CI_TYPE_UPDATE int = iota
	CI_TYPE_DELETE
)

var Host string
var FileHistoryDir string
var FileSourcesDir string

func init() {
	flag.StringVar(&Host, "host", "127.0.0.1:8080", "")
	flag.StringVar(&FileHistoryDir, "file-history", "/data/file-history", "")
	flag.StringVar(&FileSourcesDir, "file-sources", "/data/file-sources", "")
	flag.Parse()
}

func main() {
	log.Printf("Host: %s", Host)
	log.Printf("FileHistory: %s", FileHistoryDir)
	log.Printf("FileSources: %s", FileSourcesDir)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(Host, nil); err != nil {
		log.Print(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	op, err := ParseOperation(r)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	switch op.Type {
	case CREATE:
		err = createOperation(op)
	case DELETE:
		err = deleteOperation(op)
	}
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte("ok"))
	return
}

func createOperation(op Operation) error {
	if err := StoreToFileHistory(op); err != nil {
		return err
	}
	if err := AddFileToFileSources(op); err != nil {
		return err
	}
	if err := CommitToFileSources(op); err != nil {
		return err
	}
	return nil
}

func deleteOperation(op Operation) error {
	if err := DeleteFileFromFileSources(op); err != nil {
		return err
	}
	if err := CommitToFileSources(op); err != nil {
		return err
	}
	return nil
}

func StoreToFileHistory(op Operation) error {
	file := path.Join(FileHistoryDir, op.Key)
	fileInfo, err := os.Stat(file)
	if err == nil && fileInfo.IsDir() {
		return errors.New("FileHistory: File already exist and is directory: " + file)
	}
	if err = WriteFile(file, op.Bytes); err != nil {
		return err
	}
	return nil
}

func CommitToFileSources(op Operation) error {
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
	err := os.Remove(file)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
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
	if err = WriteFile(file, []byte(op.Key)); err != nil {
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
