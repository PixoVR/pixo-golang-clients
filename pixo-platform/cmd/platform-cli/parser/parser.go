package parser

import (
	"errors"
	"gopkg.in/ini.v1"
)

type IniParser struct {
	filepath string
	iniFile  *ini.File
}

func NewIniParser(input *string) (*IniParser, error) {
	var filepath string

	if input != nil {
		filepath = *input
	} else {
		filepath = DefaultConfigFilepath
	}

	f, err := ini.Load(filepath)
	if err != nil {
		return nil, err
	}

	return &IniParser{
		iniFile:  f,
		filepath: filepath,
	}, nil
}

func (i *IniParser) ParseServerVersion() (string, error) {
	if section := i.iniFile.Section("Project"); section != nil {
		return section.Key("ProjectVersion").String(), nil
	} else {
		return "", errors.New("could not find Project section in ini file")
	}
}
