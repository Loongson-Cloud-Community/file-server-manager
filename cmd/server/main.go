package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
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
	Host = "127.0.0.1:8080"
	FileHistoryDir = "/tmp/file-uploader/file-history"
	FileSourcesDir = "/tmp/file-uploader/file-sources"
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(Host, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var op Operation

	check := func(r *http.Request) error {
		return nil
	}
	if err = check(r); err != nil {
		goto fail
	}
	if op, err = ParseOperation(r); err != nil {
		goto fail
	}
	if err = StoreToFileHistory(op); err != nil {
		goto fail
	}
	if err = CommitToFileSources(op); err != nil {
		goto fail
	}

	w.Write([]byte("ok"))
	return

fail:
	w.Write([]byte(err.Error()))
	return
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
}

func StoreToFileHistory(op Operation) error {
	var err error
	var fileInfo os.FileInfo

	key := GenerateKey(op.Bytes)
	file := path.Join(FileHistoryDir, key)
	fileInfo, err = os.Stat(file)
	if err == nil && fileInfo.IsDir() {
		return errors.New("FileHistory: File already exist and is directory: " + file)
	}
	if err = WriteFile(file, op.Bytes); err != nil {
		return err
	}
	return nil
}

func CommitToFileSources(op Operation) error {
	// 判断目录
	var err error
	var fileInfo os.FileInfo

	dir := path.Join(FileSourcesDir, op.Directory)
	fileInfo, err = os.Stat(dir)
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
	key := GenerateKey(op.Bytes)
	if err = WriteFile(file, []byte(key)); err != nil {
		return err
	}

	// Git Commit
	return nil
}

func GenerateKey(bytes []byte) string {
	var key string
	md5sum := md5.Sum(bytes)
	key = hex.EncodeToString(md5sum[:])
	return key
}
