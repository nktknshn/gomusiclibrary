package tests

import (
	"path"
	"runtime"
)

func CurrentTestFolder() string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Dir(filename)
}

func CurrentTestPath(subpath string) string {
	_, filename, _, _ := runtime.Caller(1)

	return path.Join(path.Dir(filename), subpath)
}
