package models

import "io"

type WriteSeekCloser interface {
	io.WriteSeeker
	io.Closer
}
