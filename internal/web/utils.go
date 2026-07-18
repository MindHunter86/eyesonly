package web

import (
	"fmt"
	"io"
	"io/fs"
)

func openFile(path string) (_ []byte, e error) {
	var fd fs.File
	if fd, e = Static.Open(path); e != nil {
		return
	}
	defer fd.Close()

	var fifo fs.FileInfo
	if fifo, e = fd.Stat(); e != nil {
		return
	}

	if fifo.IsDir() {
		return nil, fmt.Errorf("is file %s a directory?", path)
	}

	var buf []byte
	if buf, e = io.ReadAll(fd); e != nil {
		return
	}

	return buf, nil
}
