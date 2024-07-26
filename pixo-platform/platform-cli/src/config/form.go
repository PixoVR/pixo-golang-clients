package config

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform-cli/src/forms"
	"github.com/spf13/cobra"
	"reflect"
	"strconv"
	"strings"
)

type Value struct {
	forms.Question
}

func (c *ConfigManager) GetValuesOrSubmitForm(values []Value, cmd *cobra.Command) (map[string]interface{}, error) {
	vals := make(map[string]interface{})
	questions := make([]forms.Question, 0)

	for _, value := range values {
		switch value.Type {
		case forms.Input, forms.SensitiveInput, forms.Select:
			val, ok := c.GetFlagOrConfigValue(value.Question.Key, cmd)
			if ok {
				vals[value.Question.Key] = forms.String(val)
			} else {
				questions = append(questions, value.Question)
			}
		case forms.Confirm:
			val, ok := c.GetBoolFlagOrConfigValue(value.Question.Key, cmd)
			if ok {
				vals[value.Question.Key] = forms.Bool(val)
			} else {
				questions = append(questions, value.Question)
			}
		case forms.SelectID:
			val, ok := c.GetFlagOrConfigValue(value.Question.Key, cmd)
			if ok {
				var id int
				options, err := value.GetOptions()
				if err != nil {
					return nil, err
				}
				for _, option := range options {
					if val == option.Label {
						id, _ = strconv.Atoi(option.Value)
						break
					}
				}
				vals[value.Question.Key] = id
			} else {
				questions = append(questions, value.Question)
			}
		case forms.MultiSelect:
			val, ok := c.GetFlagOrConfigValue(value.Question.Key, cmd)
			if ok {
				vals[value.Question.Key] = forms.StringSlice(strings.Split(val, ","))
			} else {
				questions = append(questions, value.Question)
			}
		case forms.MultiSelectIDs:
			val, ok := c.GetFlagOrConfigValue(value.Question.Key, cmd)
			if ok {
				strVals := strings.Split(val, ",")
				ids := make([]int, len(strVals))
				for i, strVal := range strVals {
					var id int
					options, err := value.GetOptions()
					if err != nil {
						return nil, err
					}
					for _, option := range options {
						if strVal == option.Label {
							id, _ = strconv.Atoi(option.Value)
							ids[i] = id
							break
						}
					}
				}
				vals[value.Question.Key] = forms.IntSlice(ids)
			} else {
				questions = append(questions, value.Question)
			}
		default:
			return nil, fmt.Errorf("unsupported question type: %s", reflect.ValueOf(value.Question.Type).String())
		}
	}

	answers, err := c.formHandler.AskQuestions(questions)
	if err != nil {
		return nil, err
	}

	for key, val := range answers {
		vals[key] = val
	}

	return vals, nil
}
