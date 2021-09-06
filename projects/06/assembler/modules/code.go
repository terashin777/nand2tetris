package modules

import (
	"fmt"
	"strconv"

	"github.com/terashin777/assembler/model"
	"github.com/terashin777/assembler/table"
)

type Code struct{}

func (c *Code) Dest(s string) (byte, error) {
	return c.toBinary(table.Dest, s)
}

func (c *Code) Comp(s string) (byte, error) {
	return c.toBinary(table.Comp, s)
}

func (c *Code) Jump(s string) (byte, error) {
	return c.toBinary(table.Jump, s)
}

func (c *Code) toBinary(t model.Table, s string) (byte, error) {
	bi := t.ToBinary(s)
	if bi == "" {
		return 0, fmt.Errorf("%s is invalid", t.Name())
	}

	by, err := strconv.ParseUint(bi, 2, 8)
	if err != nil {
		return 0, fmt.Errorf("%s binary is invalid", t.Name())
	}

	return uint8(by), nil
}
