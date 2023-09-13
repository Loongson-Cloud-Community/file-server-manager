package main

import (
	"os"
	"path"
	"testing"
)

func TestWriteFile(t *testing.T) {
	/* 模拟一个文件结构
	- /tmp/file_uploader/
	  - dir/
	  - file   -- 1
	*/
	rootDir := "/tmp/file_uploader"
	subDir := path.Join(rootDir, "dir")
	subFile := path.Join(rootDir, "file")
	os.RemoveAll(rootDir)
	os.MkdirAll(rootDir, os.ModePerm)
	os.MkdirAll(subDir, os.ModePerm)
	f, _ := os.Create(subFile)
	f.Write([]byte("1"))
	f.Close()

	type args struct {
		file  string
		bytes []byte
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new file",
			args: args{
				file:  "/tmp/file_uploader/dir/file",
				bytes: []byte("Hello, World!"),
			},
		},
		{
			name: "replace file",
			args: args{
				file:  "/tmp/file_uploader/file",
				bytes: []byte("2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteFile(tt.args.file, tt.args.bytes)
			if err != nil != tt.wantErr {
				t.Errorf("WriteToLocalFileSystem(%s, ...) err: %v", tt.args.file, err)
			}
		})
	}
}
