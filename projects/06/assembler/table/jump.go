package table

import (
	"fmt"

	"github.com/terashin777/assembler/utils"
)

var Jump jump = jump{}

type jump struct{}

func (t jump) Name() string {
	return "jump"
}

func (t jump) ToBinary(mn string) (byte, error) {
	switch mn {
	case "":
		return utils.StringUtil.ToBinaryNoError("000"), nil
	case "JGT":
		return utils.StringUtil.ToBinaryNoError("001"), nil
	case "JEQ":
		return utils.StringUtil.ToBinaryNoError("010"), nil
	case "JGE":
		return utils.StringUtil.ToBinaryNoError("011"), nil
	case "JLT":
		return utils.StringUtil.ToBinaryNoError("100"), nil
	case "JNE":
		return utils.StringUtil.ToBinaryNoError("101"), nil
	case "JLE":
		return utils.StringUtil.ToBinaryNoError("110"), nil
	case "JMP":
		return utils.StringUtil.ToBinaryNoError("111"), nil
	}

	return 0, fmt.Errorf("not match")
}
