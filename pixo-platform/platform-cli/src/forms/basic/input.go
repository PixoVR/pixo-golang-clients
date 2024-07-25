package basic

import (
	"bufio"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
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

func (f *Handler) GetResponseFromUser(prompt string, response *string) error {
	prompt = strings.ReplaceAll(prompt, "-", " ")
	prompt = emoji.Sprintf(":fountain_pen: Enter %s: ", prompt)
	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return err
	}

	answer, err := f.ReadLine()
	if err != nil {
		return err
	}

	if response != nil {
		*response = answer
	}
	return nil
}

func (f *Handler) GetSensitiveResponseFromUser(prompt string, response *string) error {
	prompt = strings.ReplaceAll(prompt, "-", " ")
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

	log.Debug().Str("response", strings.Trim(string(val), "\r\n")).Msg("User response")
	if response != nil {
		*response = trim(string(val))
	}
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
