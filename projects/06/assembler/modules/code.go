package modules

import (
	"github.com/terashin777/assembler/table"
)

type Code struct{}

func (c *Code) Dest(s string) (byte, error) {
	return table.Dest.ToBinary(s)
}

func (c *Code) Comp(s string) (byte, error) {
	return table.Comp.ToBinary(s)
}

func (c *Code) Jump(s string) (byte, error) {
	return table.Jump.ToBinary(s)
}
