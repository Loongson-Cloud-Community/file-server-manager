package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path"
)

const (
	CI_TYPE_UPDATE int = iota
	CI_TYPE_DELETE
)

const (
	FileHistoryDirName = "file-history"
	FileSourcesDirName = "file-sources"
	FileServerDirName  = "file-server"
)

var host string
var data string
var FileHistoryDir string
var FileSourcesDir string
var FileServerDir string

func init() {
	flag.StringVar(&host, "host", "127.0.0.1:8080", "")
	flag.StringVar(&data, "data", "/data", "")

	FileHistoryDir = path.Join(data, FileHistoryDirName)
	FileSourcesDir = path.Join(data, FileSourcesDirName)
	FileServerDir = path.Join(data, FileServerDirName)
}

func main() {
	flag.Parse()
	log.Printf("Host: %s", host)
	log.Printf("FileHistory: %s", FileHistoryDir)
	log.Printf("FileSources: %s", FileSourcesDir)
	log.Printf("FileServer: %s", FileServerDir)
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(host, nil); err != nil {
		log.Print(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	op, err := ParseOperation(r)
	if err != nil {
		goto fail
	}
	if err = StoreToFileHistory(op); err != nil {
		goto fail
	}
	if err = CommitToFileSources(op); err != nil {
		goto fail
	}
	if err = IncSyncToFileServer(op); err != nil {
		goto fail
	}

	w.Write([]byte("ok"))
	log.Printf("[Result]:\tvalid=%v %s", op.Valid, "ok")
	return
fail:
	w.Write([]byte(err.Error()))
	log.Printf("[Result]:\tvalid=%v %s", op.Valid, err.Error())
	return
}
