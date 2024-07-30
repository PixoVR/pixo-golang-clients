package config

import "strconv"

func ToInt(val string) (int, bool) {
	i, err := strconv.Atoi(val)
	return i, err == nil && i != 0
}

func ToBool(val string) (bool, bool) {
	if val == "yes" {
		return true, true
	} else if val == "no" {
		return false, true
	}

	b, err := strconv.ParseBool(val)
	return b, err == nil
}
