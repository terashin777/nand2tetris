package table

import (
	"fmt"

	"github.com/terashin777/assembler/utils"
)

var Comp comp = comp{}

type comp struct{}

func (t comp) Name() string {
	return "comp"
}

func (t comp) ToBinary(mn string) (byte, error) {
	switch mn {
	case "0":
		return utils.StringUtil.ToBinaryNoError("0101010"), nil
	case "1":
		return utils.StringUtil.ToBinaryNoError("0111111"), nil
	case "-1":
		return utils.StringUtil.ToBinaryNoError("0111010"), nil
	case "D":
		return utils.StringUtil.ToBinaryNoError("0001100"), nil
	case "A":
		return utils.StringUtil.ToBinaryNoError("0110000"), nil
	case "M":
		return utils.StringUtil.ToBinaryNoError("1110000"), nil
	case "!D":
		return utils.StringUtil.ToBinaryNoError("0001101"), nil
	case "!A":
		return utils.StringUtil.ToBinaryNoError("0110001"), nil
	case "!M":
		return utils.StringUtil.ToBinaryNoError("1110001"), nil
	case "-D":
		return utils.StringUtil.ToBinaryNoError("0001111"), nil
	case "-A":
		return utils.StringUtil.ToBinaryNoError("0110011"), nil
	case "-M":
		return utils.StringUtil.ToBinaryNoError("1110011"), nil
	case "D+1":
		return utils.StringUtil.ToBinaryNoError("0011111"), nil
	case "A+1":
		return utils.StringUtil.ToBinaryNoError("0110111"), nil
	case "M+1":
		return utils.StringUtil.ToBinaryNoError("1110111"), nil
	case "D-1":
		return utils.StringUtil.ToBinaryNoError("0001110"), nil
	case "A-1":
		return utils.StringUtil.ToBinaryNoError("0110010"), nil
	case "M-1":
		return utils.StringUtil.ToBinaryNoError("1110010"), nil
	case "D+A":
		return utils.StringUtil.ToBinaryNoError("0000010"), nil
	case "D+M":
		return utils.StringUtil.ToBinaryNoError("1000010"), nil
	case "D-A":
		return utils.StringUtil.ToBinaryNoError("0010011"), nil
	case "D-M":
		return utils.StringUtil.ToBinaryNoError("1010011"), nil
	case "A-D":
		return utils.StringUtil.ToBinaryNoError("0000111"), nil
	case "M-D":
		return utils.StringUtil.ToBinaryNoError("1000111"), nil
	case "D&A":
		return utils.StringUtil.ToBinaryNoError("0000000"), nil
	case "D&M":
		return utils.StringUtil.ToBinaryNoError("1000000"), nil
	case "D|A":
		return utils.StringUtil.ToBinaryNoError("0010101"), nil
	case "D|M":
		return utils.StringUtil.ToBinaryNoError("1010101"), nil
	}

	return 0, fmt.Errorf("not match")
}
