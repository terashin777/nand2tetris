package utils

type char struct{}

var Char = char{}

func (u char) IsAlphabet(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z'
}

func (u char) IsNumber(r rune) bool {
	return '0' <= r && r <= '9'
}
