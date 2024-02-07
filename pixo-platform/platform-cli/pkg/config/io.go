package config

import (
	"github.com/kyokomi/emoji"
	"io"
	"os"
)

func (f *fileManagerImpl) SetReader(r io.Reader) {
	f.reader = r
}

func (f *fileManagerImpl) readerOrStdin() io.Reader {
	if f.reader == nil {
		return os.Stdin
	}
	return f.reader
}

func (f *fileManagerImpl) Read(p []byte) (n int, err error) {
	return f.readerOrStdin().Read(p)
}

func (f *fileManagerImpl) SetWriter(w io.Writer) {
	f.writer = w
}

func (f *fileManagerImpl) writerOrStdout() io.Writer {
	if f.writer == nil {
		return os.Stdout
	}
	return f.writer
}

func (f *fileManagerImpl) Write(p []byte) (n int, err error) {
	return f.writerOrStdout().Write(p)
}

func (f *fileManagerImpl) Print(a ...interface{}) {
	msg := emoji.Sprint(a...)
	_, _ = f.Write([]byte(msg))
}

func (f *fileManagerImpl) Println(a ...interface{}) {
	msg := emoji.Sprint(a...)
	_, _ = f.Write([]byte(msg))
	_, _ = f.Write([]byte("\n"))
}

func (f *fileManagerImpl) Printf(format string, a ...interface{}) {
	msg := emoji.Sprintf(format, a...)
	_, _ = f.Write([]byte(msg))
}
