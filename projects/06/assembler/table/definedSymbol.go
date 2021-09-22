package table

import "fmt"

var DefinedSymbolTable = map[string]uint16{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"SCREEN": 16384,
	"KBD":    24576,
}

func init() {
	for i := 0; i < 16; i++ {
		DefinedSymbolTable[fmt.Sprintf("R%d", i)] = uint16(i)
	}
}
