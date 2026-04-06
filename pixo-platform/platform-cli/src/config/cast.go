package config

import "strconv"

func ToInt(val string) (int, bool) {
	i, err := strconv.Atoi(val)
	return i, err == nil && i != 0
}

func ToBool(val string) (bool, bool) {
	switch val {
	case "yes":
		return true, true
	case "no":
		return false, true
	}

	b, err := strconv.ParseBool(val)
	return b, err == nil
}
