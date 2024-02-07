package editor

import (
	"errors"
	"os"
	"os/exec"
)

type FileOpener interface {
	OpenEditor(fileName string) error
}

type fileOpenerImpl struct {
	editor string
}

func NewFileOpener(editor string) FileOpener {
	return &fileOpenerImpl{editor: editor}
}

func (f *fileOpenerImpl) OpenEditor(fileName string) error {
	if fileName == "" {
		return errors.New("file name cannot be empty")
	}

	if f.editor == "" {
		f.editor = f.getEditor()
	}

	cmd := exec.Command(f.editor, fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (f *fileOpenerImpl) getEditor() string {
	for _, envVar := range []string{"VISUAL", "EDITOR"} {
		if editor, exists := os.LookupEnv(envVar); exists {
			return editor
		}
	}

	return "vim"
}
