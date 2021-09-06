package table

var Jump jump = jump{}

type jump struct{}

func (t jump) Name() string {
	return "jump"
}

func (t jump) ToBinary(mn string) string {
	switch mn {
	case "":
		return "000"
	case "JGT":
		return "001"
	case "JEQ":
		return "010"
	case "JGE":
		return "011"
	case "JLT":
		return "100"
	case "JNE":
		return "101"
	case "JLE":
		return "110"
	case "JMP":
		return "111"
	}

	return ""
}
