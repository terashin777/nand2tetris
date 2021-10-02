package utils

import "strconv"

var StringUtil = stringUtil{}

type stringUtil struct{}

func (u stringUtil) ToBinary(s string) (byte, error) {
	i, err := strconv.ParseUint(s, 2, 8)
	if err != nil {
		return 0, err
	}

	return uint8(i), nil
}

func (u stringUtil) ToBinaryNoError(s string) byte {
	i, _ := strconv.ParseUint(s, 2, 8)
	return uint8(i)
}
