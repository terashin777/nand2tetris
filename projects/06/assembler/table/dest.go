package table

import (
	"fmt"

	"github.com/terashin777/assembler/utils"
)

var Dest dest = dest{}

type dest struct{}

func (t dest) Name() string {
	return "dest"
}

func (t dest) ToBinary(mn string) (byte, error) {
	switch mn {
	case "":
		return utils.StringUtil.ToBinaryNoError("000"), nil
	case "M":
		return utils.StringUtil.ToBinaryNoError("001"), nil
	case "D":
		return utils.StringUtil.ToBinaryNoError("010"), nil
	case "MD":
		return utils.StringUtil.ToBinaryNoError("011"), nil
	case "A":
		return utils.StringUtil.ToBinaryNoError("100"), nil
	case "AM":
		return utils.StringUtil.ToBinaryNoError("101"), nil
	case "AD":
		return utils.StringUtil.ToBinaryNoError("110"), nil
	case "AMD":
		return utils.StringUtil.ToBinaryNoError("111"), nil
	}

	return 0, fmt.Errorf("not match")
}
