package basic

import (
	"bufio"
	"github.com/kyokomi/emoji"
	"github.com/rs/zerolog/log"
	"gitlab.com/david_mbuvi/go_asterisks"
	"io"
	"strings"
)

func (f *FormHandler) GetResponseFromUser(prompt string) (string, error) {
	prompt = emoji.Sprintf(":fountain_pen: Enter %s: ", prompt)
	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return "", err
	}

	if f.buffer == nil {
		f.buffer = bufio.NewReader(f.readerOrStdin())
	}

	response, err := f.buffer.ReadString('\n')
	if err != nil {
		return "", err
	}

	trimmedResponse := strings.Trim(response, "\r\n")
	log.Debug().Str("response", trimmedResponse).Msg("User response")
	return trimmedResponse, nil
}

func (f *FormHandler) GetSensitiveResponseFromUser(prompt string) (string, error) {
	prompt = emoji.Sprintf(":lock: Enter %s: ", prompt)
	if _, err := f.writerOrStdout().Write([]byte(prompt)); err != nil {
		return "", err
	}

	if f.buffer == nil {
		f.buffer = bufio.NewReader(f.readerOrStdin())
	}

	var fd uintptr
	customReader := NewCustomFieldReader(f.buffer, fd)
	val, err := go_asterisks.GetUsersPassword(prompt, true, customReader, f.writer)
	if err != nil {
		return "", err
	}

	log.Debug().Str("response", strings.Trim(string(val), "\r\n")).Msg("User response")
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
