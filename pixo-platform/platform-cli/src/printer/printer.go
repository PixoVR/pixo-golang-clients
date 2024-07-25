package printer

import (
	"github.com/kyokomi/emoji"
	"io"
	"os"
)

var _ Printer = (*EmojiPrinter)(nil)

type Printer interface {
	io.Writer
	SetWriter(w io.Writer)
	Print(a ...interface{})
	Printf(format string, a ...interface{})
	Println(a ...interface{})
}

type EmojiPrinter struct {
	writer io.Writer
}

func NewEmojiPrinter(writer io.Writer) *EmojiPrinter {
	return &EmojiPrinter{writer: writer}
}

func (e *EmojiPrinter) Printf(format string, a ...interface{}) {
	msg := emoji.Sprintf(format, a...)
	_, _ = e.writerOrStdout().Write([]byte(msg))
}

func (e *EmojiPrinter) Print(a ...interface{}) {
	msg := emoji.Sprint(a...)
	_, _ = e.writerOrStdout().Write([]byte(msg))
}

func (e *EmojiPrinter) Println(a ...interface{}) {
	msg := emoji.Sprint(a...)
	_, _ = e.writerOrStdout().Write([]byte(msg))
	_, _ = e.writerOrStdout().Write([]byte("\n"))
}

func (e *EmojiPrinter) Write(p []byte) (n int, err error) {
	return e.writerOrStdout().Write(p)
}

func (e *EmojiPrinter) SetWriter(w io.Writer) {
	e.writer = w
}

func (e *EmojiPrinter) writerOrStdout() io.Writer {
	if e.writer == nil {
		return io.Writer(os.Stdout)
	}
	return e.writer
}
