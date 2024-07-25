package basic

import (
	"bufio"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"io"
	"os"
)

var _ forms.FormHandler = &Handler{}

type Handler struct {
	reader io.Reader
	writer io.Writer
	buffer *bufio.Reader
}

func NewFormHandler(input io.Reader, output io.Writer) *Handler {
	return &Handler{
		reader: input,
		writer: output,
	}
}

func (f *Handler) SetReader(r io.Reader) {
	f.reader = r
}

func (f *Handler) SetWriter(w io.Writer) {
	f.writer = w
}

func (f *Handler) writerOrStdout() io.Writer {
	if f.writer == nil {
		return os.Stdout
	}
	return f.writer
}

func (f *Handler) readerOrStdin() io.Reader {
	if f.reader == nil {
		return os.Stdin
	}
	return f.reader
}
