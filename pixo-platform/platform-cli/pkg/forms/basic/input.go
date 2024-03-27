package basic

import (
	"bufio"
	"strings"
)

func (f *FormHandler) GetResponseFromUser(prompt string) (string, error) {
	if _, err := f.writer.Write([]byte(prompt)); err != nil {
		return "", err
	}

	bytesReader := bufio.NewReader(f.reader)
	response, err := bytesReader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.Trim(response, "\r\n"), nil
}
