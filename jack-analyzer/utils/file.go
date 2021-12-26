package utils

import (
	"os"
	"path/filepath"
)

func ReplaceExt(fn, newExt string) string {
	ext := filepath.Ext(fn)
	return fn[0:len(fn)-len(ext)] + newExt
}

func MakeSameDirFileName(path, ext string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if fi.IsDir() {
		return MakeSameDirName(path, ext)
	} else {
		return MakeSameFileName(path, ext)
	}
}

func MakeSameDirName(path, ext string) (string, error) {
	f := filepath.Base(path)
	return filepath.Join(path, ReplaceExt(f, ext)), nil
}

func MakeSameFileName(path, ext string) (string, error) {
	dir := filepath.Dir(path)
	f := filepath.Base(path)
	return filepath.Join(dir, ReplaceExt(f, ext)), nil
}
