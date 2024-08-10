package shell

import "runtime"

func OSFeature[T any](winFeature T, linuxFeature T) T {
	if runtime.GOOS == "windows" {
		return winFeature
	}

	return linuxFeature
}
