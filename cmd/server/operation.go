package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"path"
	"strings"
)

const (
	CREATE int = iota
	DELETE
)

// Operation 对资源的操作 例如 POST https://cloud.loongnix.cn/org/repo/ver/file
//
// Type 操作类型 POST -> CREATE，DELETE -> DELETE
// Directory org/repo/ver
// Name file
type Operation struct {
	Type      int
	Directory string
	File      string
	Bytes     []byte
	Key       string
}

func ParseOperation(r *http.Request) (Operation, error) {
	var dir, file string
	var err error
	var op Operation
	var t int
	var bs []byte

	if dir, file, err = parseUrl(r.URL.Path); err != nil {
		return op, err
	}

	switch r.Method {
	case http.MethodPost:
		t = CREATE
	case http.MethodDelete:
		t = DELETE
	default:
		return op, errors.New("Invalid Request Method: " + r.Method)
	}

	if op.Type == CREATE {
		bs, err = io.ReadAll(r.Body)
		if err != nil {
			return op, errors.New("Read from request failed:" + err.Error())
		}
	}
	op.Type = t
	op.Directory = dir
	op.File = file
	op.Bytes = bs
  op.Key = GenerateKey(bs)

	return op, nil
}

func GenerateKey(bytes []byte) string {
	var key string
	md5sum := md5.Sum(bytes)
	key = hex.EncodeToString(md5sum[:])
	return key
}

func parseUrl(url string) (dir string, file string, err error) {
	var turl string

	turl = strings.Trim(url, "/")
	if !strings.ContainsRune(turl, '/') {
		return "", "", errors.New("Invalid url " + url)
	}
	dir = path.Dir(turl)
	file = path.Base(turl)
	return dir, file, nil
}
