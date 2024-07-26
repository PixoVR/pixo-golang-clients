package basic

import (
	"bufio"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/kyokomi/emoji"
	"gitlab.com/david_mbuvi/go_asterisks"
	"io"
	"strings"
)

func (f *Handler) ReadLine() (string, error) {
	if f.buffer == nil {
		f.buffer = bufio.NewReader(f.readerOrStdin())
	}

	line, err := f.buffer.ReadString('\n')
	if err != nil {
		return "", err
	}

	return trim(line), nil
}

func (f *Handler) GetResponseFromUser(question *forms.Question) error {
	prompt := strings.ReplaceAll(question.Prompt, "-", " ")
	prompt = emoji.Sprintf(":fountain_pen: Enter %s: ", prompt)

	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return err
	}

	line, err := f.ReadLine()
	if err != nil {
		return err
	}

	question.Answer = forms.String(line)
	return nil
}

func (f *Handler) GetSensitiveResponseFromUser(question *forms.Question) error {
	prompt := strings.ReplaceAll(question.Prompt, "-", " ")
	prompt = emoji.Sprintf(":lock: Enter %s: ", prompt)
	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return err
	}

	if f.buffer == nil {
		f.buffer = bufio.NewReader(f.readerOrStdin())
	}

	var fd uintptr
	customReader := NewCustomFieldReader(f.buffer, fd)
	val, err := go_asterisks.GetUsersPassword(prompt, true, customReader, f.writer)
	if err != nil {
		return err
	}

	question.Answer = trim(string(val))
	return nil
}

type CustomFieldReader struct {
	*bufio.Reader
	fd uintptr
}

func (c *CustomFieldReader) Fd() uintptr {
	return c.fd
}

func NewCustomFieldReader(reader io.Reader, fd uintptr) *CustomFieldReader {
	return &CustomFieldReader{
		Reader: bufio.NewReader(reader),
		fd:     fd,
	}
}

func trim(line string) string {
	return strings.Trim(line, "\r\n")
}
