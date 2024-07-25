package basic

import (
	"bytes"
	"io"
	"os"
)

type FormHandler struct {
	reader io.Reader
	writer io.Writer
	buffer *bytes.Buffer
}

func NewFormHandler(input io.Reader, output io.Writer) *FormHandler {
	return &FormHandler{
		reader: input,
		writer: output,
		buffer: new(bytes.Buffer),
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
