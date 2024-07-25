package config

import "strconv"

func ToInt(val string) (int, bool) {
	i, err := strconv.Atoi(val)
	return i, err == nil && i != 0
}

func ToBool(val string) (bool, bool) {
	if val == "yes" {
		return true, false
	} else if val == "no" {
		return false, false
	}

	b, err := strconv.ParseBool(val)
	return b, err == nil
}
