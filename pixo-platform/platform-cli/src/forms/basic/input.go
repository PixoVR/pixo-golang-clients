package basic

import (
	"bufio"
	"github.com/kyokomi/emoji"
	"gitlab.com/david_mbuvi/go_asterisks"
	"io"
	"os"
	"strings"
)

func (f *FormHandler) GetResponseFromUser(prompt string) (string, error) {
	prompt = emoji.Sprintf(":fountain_pen: Enter %s: ", prompt)
	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return "", err
	}

	reader := bufio.NewReader(io.MultiReader(f.buffer, f.readerOrStdin()))
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	remaining, err := io.ReadAll(reader)
	if err != nil && err != io.EOF {
		return "", err
	}
	f.buffer.Reset()
	f.buffer.Write(remaining)

	return strings.Trim(response, "\r\n"), nil
}

func (f *FormHandler) GetSensitiveResponseFromUser(prompt string) (string, error) {
	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return "", err
	}

	prompt = emoji.Sprintf(":lock: Enter %s: ", prompt)
	var reader io.Reader
	var fd uintptr

	if f.readerOrStdin() != os.Stdin {
		reader = io.MultiReader(f.buffer, f.reader)
	} else {
		reader = os.Stdin
	}

	if reader == os.Stdin {
		fd = os.Stdin.Fd()
	} else if file, ok := reader.(*os.File); ok {
		fd = file.Fd()
	} else {
		fd = 0
	}

	customReader := NewCustomFieldReader(reader, fd)
	val, err := go_asterisks.GetUsersPassword(prompt, true, customReader, f.writer)
	if err != nil {
		return "", err
	}

	return strings.Trim(string(val), "\r\n"), err
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
