package hls

import (
	"log"
	"os"
)

// write file
type fileWriter struct {
	file string
	f    *os.File
}

func newFileWriter() *fileWriter {
	return &fileWriter{}
}

func (fw *fileWriter) open(file string) (err error) {

	if nil != fw.f {
		// already opened
		return
	}

	fw.f, err = os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Println("open file error, file=", file)
		return
	}

	fw.file = file

	return
}

func (fw *fileWriter) close() {

	if nil == fw.f {
		return
	}

	fw.f.Close()

	// after close, rest the file write to nil
	fw.f = nil

}

func (fw *fileWriter) isOpen() bool {
	return nil == fw.f
}

// write data to file
func (fw *fileWriter) write(buf []byte) (err error) {

	if _, err = fw.f.Write(buf); err != nil {
		log.Println("write to file failed, file=", fw.file, ",err=", err)
		return
	}

	return
}
