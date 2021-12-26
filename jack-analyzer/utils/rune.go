package utils

var newlineCodes map[rune]struct{} = map[rune]struct{}{
	0x000A: {},
	0x000B: {},
	0x000C: {},
	0x000D: {},
	0x0085: {},
	0x2028: {},
	0x2029: {},
}

func IsNewline(r rune) bool {
	_, ok := newlineCodes[r]
	return ok
}
