package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"time"
)

func StoreToFileHistory(op *Operation) error {
	if op.Type == DELETE {
		return nil
	}
	// store tmp file
	tmpName := path.Join(FileHistoryDir, tmpFileName())
	f, err := os.Create(tmpName)
	if err != nil {
		return err
	}
	io.Copy(f, op.Body)
	f.Sync()
	// calc hash
	// set offset to head first
	f.Seek(0, 0)
	d := md5.New()
	io.Copy(d, f)
	sum := d.Sum(nil)
	hx := hex.EncodeToString(sum[:])
	// close file descriptor now
	f.Close()
	// rename tmp file to real file
	realName := path.Join(FileHistoryDir, hx)
	os.Rename(tmpName, realName)
	// don't forget to assign to hash
	op.Hash = hx

	log.Printf("[file-history]:\t%s/%s %s", op.Directory, op.File, op.Hash)
	return nil
}

func tmpFileName() string {
	return fmt.Sprintf("tmp.%d-%d", rand.Intn(1024), time.Now().UnixNano())
}
