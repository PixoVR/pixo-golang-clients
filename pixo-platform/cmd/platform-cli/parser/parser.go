package parser

import (
	"errors"
	"fmt"
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

func (i *IniParser) ParseSemanticVersion(iniInfo ...string) (string, error) {

	var sectionName string
	var key string

	if len(iniInfo) == 0 {
		sectionName = DefaultVersionSectionName
		key = DefaultVersionKey
	}

	if len(iniInfo) > 0 {
		sectionName = iniInfo[0]
	}

	if len(iniInfo) > 1 {
		key = iniInfo[1]
	}

	if section := i.iniFile.Section(sectionName); section != nil {
		return section.Key(key).String(), nil
	} else {
		return "", errors.New(fmt.Sprintf("could not find %s section in ini file", sectionName))
	}
}
