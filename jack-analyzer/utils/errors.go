package utils

import (
	"bytes"
	"runtime/debug"
)

func Stack(skip int) string {
	s := debug.Stack()
	ss := bytes.Split(s, []byte("\n"))

	return string(bytes.Join(ss[1+skip*2:], []byte("\n")))
}
