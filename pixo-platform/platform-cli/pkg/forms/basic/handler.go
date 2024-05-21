package basic

import "io"

type FormHandler struct {
	reader io.Reader
	writer io.Writer
}

func NewFormHandler(input io.ReadWriter, output io.ReadWriter) *FormHandler {
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
