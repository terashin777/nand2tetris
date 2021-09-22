package modules

import (
	"bufio"
	"github.com/terashin777/vm-translator/models"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	s     *bufio.Scanner
	ct    models.CommandType
	cur   string
	parts []string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: bufio.NewScanner(r),
	}
}

func (p *Parser) Advance() error {
	next := p.s.Scan()
	if !next {
		return io.EOF
	}

	p.cur = p.s.Text()
	p.parts = strings.Split(p.cur, " ")
	return nil
}

func (p *Parser) CommandType() models.CommandType {
	switch true {
	case p.isPush():
		return models.C_PUSH
	case p.isPop():
		return models.C_POP
	case p.isArithmetic():
		return models.C_ARITHMETIC
	default:
		return models.C_NONE
	}
}

func (p *Parser) isPush() bool {
	return len(p.parts) == 3 && p.parts[0] == "push"
}

func (p *Parser) isPop() bool {
	return len(p.parts) == 3 && p.parts[0] == "push"
}

func (p *Parser) isArithmetic() bool {
	return len(p.parts) == 1
}

func (p *Parser) Arg1() string {
	if len(p.parts) < 1 {
		return ""
	}
	if len(p.parts) == 1 {
		return p.parts[0]
	}

	return p.parts[1]
}

func (p *Parser) Arg2() int {
	if len(p.parts) < 3 {
		return -1
	}

	i, err := strconv.ParseInt(p.parts[2], 10, 16)
	if err != nil {
		return -1
	}

	return int(i)
}
