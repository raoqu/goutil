package types

import "strings"

func Split2(str string, split string) (string, string) {
	parts := strings.SplitN(str, split, 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
