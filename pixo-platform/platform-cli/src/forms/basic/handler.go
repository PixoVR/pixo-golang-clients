package basic

import (
	"bufio"
	"io"
	"os"
)

type FormHandler struct {
	reader io.Reader
	writer io.Writer
	buffer *bufio.Reader
}

func NewFormHandler(input io.Reader, output io.Writer) *FormHandler {
	return &FormHandler{
		reader: input,
		writer: output,
	}
}

func (f *FormHandler) SetReader(r io.Reader) {
	f.reader = r
}

func (f *FormHandler) SetWriter(w io.Writer) {
	f.writer = w
}

func (f *FormHandler) writerOrStdout() io.Writer {
	if f.writer == nil {
		return os.Stdout
	}
	return f.writer
}

func (f *FormHandler) readerOrStdin() io.Reader {
	if f.reader == nil {
		return os.Stdin
	}
	return f.reader
}
