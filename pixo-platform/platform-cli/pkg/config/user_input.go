package config

import (
	"bufio"
	"fmt"
	"github.com/kyokomi/emoji"
	"gitlab.com/david_mbuvi/go_asterisks"
	"io"
	"os"
	"strings"
)

func (f *fileManagerImpl) ReadFromUser(prompt string) string {
	if f.writer == nil || f.reader == nil {
		return ""
	}

	if _, err := f.writer.Write([]byte(emoji.Sprintf(":fountain_pen:Enter %s: ", prompt))); err != nil {
		return ""
	}

	bytesReader := bufio.NewReader(f.reader)
	message, _ := bytesReader.ReadString('\n')

	return strings.Trim(message, "\r\n")
}

func (f *fileManagerImpl) ReadSensitiveFromUser(prompt string) string {
	if f.writer == nil {
		return ""
	}

	prompt = emoji.Sprintf(":lock: Enter %s: ", prompt)
	var fieldReader *os.File
	bytes, err := io.ReadAll(f.reader)
	if err != nil {
		return ""
	}

	tmpFilePath := fmt.Sprintf("%s/%s", os.TempDir(), "fieldReader")
	if err = os.WriteFile(tmpFilePath, bytes, 0644); err != nil {
		return ""
	}

	fieldReader, err = os.Open(tmpFilePath)
	if err != nil {
		return ""
	}

	val, err := go_asterisks.GetUsersPassword(prompt, true, fieldReader, f.writer)
	if err != nil {
		return ""
	}

	return strings.Trim(string(val), "\r\n")
}

func (f *fileManagerImpl) ReadFromUserOrReturn(prompt, defaultValue string) string {
	val := f.ReadFromUser(prompt)
	if val == "" {
		return defaultValue
	}

	return strings.Trim(val, "\r\n")
}
