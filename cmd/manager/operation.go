package main

import (
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
	Body      io.Reader
	Hash      string
	Valid     bool
}

// ParseOperation 解析请求
// 这里仅解析 type, dir, file，不解析 body
func ParseOperation(r *http.Request) (*Operation, error) {
	var dir, file string
	var err error
	var t int
	op := &Operation{}

	if dir, file, err = parseUrl(r.URL.Path); err != nil {
		return nil, err
	}

	switch r.Method {
	case http.MethodPost:
		t = CREATE
	case http.MethodDelete:
		t = DELETE
	default:
		return nil, errors.New("Invalid Request Method: " + r.Method)
	}

	op.Type = t
	op.Directory = dir
	op.File = file
	op.Body = r.Body
	op.Valid = true

	return op, nil
}

func parseUrl(url string) (dir string, file string, err error) {
	u := strings.Trim(url, "/")
	if !strings.ContainsRune(u, '/') {
		return "", "", errors.New("Invalid url " + url)
	}
	dir = path.Dir(u)
	file = path.Base(u)
	return dir, file, nil
}
