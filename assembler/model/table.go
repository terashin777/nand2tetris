package model

type Table interface {
	Name() string
	ToBinary(s string) string
}
