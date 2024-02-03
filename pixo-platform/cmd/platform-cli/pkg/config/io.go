package config

import (
	"io"
)

func (f *fileManagerImpl) SetReader(r io.Reader) {
	f.reader = r
}

func (f *fileManagerImpl) Reader() io.Reader {
	return f.reader
}

func (f *fileManagerImpl) SetWriter(w io.Writer) {
	f.writer = w
}

func (f *fileManagerImpl) Writer() io.Writer {
	return f.writer
}
