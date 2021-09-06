package table

var Dest dest = dest{}

type dest struct{}

func (t dest) Name() string {
	return "dest"
}

func (t dest) ToBinary(mn string) string {
	switch mn {
	case "":
		return "000"
	case "M":
		return "001"
	case "D":
		return "010"
	case "MD":
		return "011"
	case "A":
		return "100"
	case "AM":
		return "101"
	case "AD":
		return "110"
	case "AMD":
		return "111"
	}

	return ""
}
