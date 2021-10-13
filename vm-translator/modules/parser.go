package modules

import (
	"bufio"
	"io"
	"strconv"
	"strings"

	"github.com/terashin777/vm-translator/models"
)

var arithmethicCmds = map[string]struct{}{
	"add": struct{}{},
	"sub": struct{}{},
	"neg": struct{}{},
	"eq":  struct{}{},
	"gt":  struct{}{},
	"lt":  struct{}{},
	"and": struct{}{},
	"or":  struct{}{},
	"not": struct{}{},
}

type Parser struct {
	s     *bufio.Scanner
	parts []string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		s: bufio.NewScanner(r),
	}
}

func (p *Parser) Advance() error {
	t, err := p.scan()
	if err != nil {
		return err
	}

	p.parts = strings.Split(t, " ")
	return nil
}

func (p *Parser) scan() (string, error) {
	for {
		next := p.s.Scan()
		if !next {
			return "", io.EOF
		}

		t := p.extractCommand(p.s.Text())
		if t == "" {
			continue
		}

		return t, nil
	}
}

func (p *Parser) extractCommand(rt string) string {
	t := strings.TrimSpace(rt)
	if p.isCommentRow(t) {
		return ""
	}
	return p.removeComment(t)
}

func (p *Parser) isCommentRow(t string) bool {
	return strings.HasPrefix(t, "//")
}

func (p *Parser) removeComment(t string) string {
	return strings.TrimSpace(strings.Split(t, "//")[0])
}

func (p *Parser) CommandType() models.CommandType {
	switch true {
	case p.isPush():
		return models.C_PUSH
	case p.isPop():
		return models.C_POP
	case p.isArithmetic():
		return models.C_ARITHMETIC
	case p.isLabel():
		return models.C_LABEL
	case p.isGoto():
		return models.C_GOTO
	case p.isIf():
		return models.C_IF
	case p.isCall():
		return models.C_CALL
	case p.isReturn():
		return models.C_RETURN
	case p.isFunction():
		return models.C_FUNCTION
	default:
		return models.C_NONE
	}
}

func (p *Parser) isPush() bool {
	return len(p.parts) == 3 && p.parts[0] == "push"
}

func (p *Parser) isPop() bool {
	return len(p.parts) == 3 && p.parts[0] == "pop"
}

func (p *Parser) isArithmetic() bool {
	if !(len(p.parts) == 1) {
		return false
	}

	_, ok := arithmethicCmds[p.parts[0]]
	return ok
}

func (p *Parser) isLabel() bool {
	return len(p.parts) == 2 && p.parts[0] == "label"
}

func (p *Parser) isGoto() bool {
	return len(p.parts) == 2 && p.parts[0] == "goto"
}

func (p *Parser) isIf() bool {
	return len(p.parts) == 2 && p.parts[0] == "if-goto"
}

func (p *Parser) isCall() bool {
	return len(p.parts) == 3 && p.parts[0] == "call"
}

func (p *Parser) isReturn() bool {
	return len(p.parts) == 1 && p.parts[0] == "return"
}

func (p *Parser) isFunction() bool {
	return len(p.parts) == 3 && p.parts[0] == "function"
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
