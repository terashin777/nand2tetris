package modules

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type (
	Parser struct {
		r   io.ReadSeeker
		c   *Code
		s   *bufio.Scanner
		cur string
	}

	CommandType int
)

const (
	A_COMMAND CommandType = iota
	C_COMMAND
	L_COMMAND
)

func NewParser(r io.ReadSeeker) *Parser {
	return &Parser{
		r: r,
		s: bufio.NewScanner(r),
		c: &Code{},
	}
}

func (p *Parser) Read() (string, error) {
	p.s.Scan()
	l := p.s.Text()
	p.cur = strings.TrimSpace(l)
	if p.cur == "" {
		return "", nil
	}
	if p.isComment() {
		return "", nil
	}

	return p.parse()
}

func (p *Parser) isComment() bool {
	return strings.HasPrefix(p.cur, "//")
}

func (p *Parser) parse() (string, error) {
	ct := p.commandType()

	var sym string
	if ct == A_COMMAND || ct == C_COMMAND {
		sym = p.symbol()
		v, err := strconv.ParseInt(sym, 10, 16)
		if err != nil {
			return "", err
		}
		return strconv.FormatInt(v, 2), nil
	}

	var dest, comp, jump string
	if ct == C_COMMAND {
		dest = p.dest()
		comp = p.comp()
		jump = p.jump()
	}

	return "", nil
}

func (p *Parser) commandType() CommandType {
	if p.isACommand() {
		return A_COMMAND
	}

	if p.isLCommand() {
		return L_COMMAND
	}

	return C_COMMAND
}

func (p *Parser) isACommand() bool {
	return strings.HasPrefix(p.cur, "@")
}

func (p *Parser) isLCommand() bool {
	return strings.HasPrefix(p.cur, "(") && strings.HasSuffix(p.cur, ")")
}

func (p *Parser) symbol() string {
	if p.cur[len(p.cur)-1] == ')' {
		return p.cur[1 : len(p.cur)-1]
	}

	return p.cur[1:]
}

func (p *Parser) dest() string {
	if i := strings.IndexByte(p.cur, '='); i != -1 {
		return p.cur[:i]
	}

	return ""
}

func (p *Parser) comp() string {
	si := 0
	ei := len(p.cur)
	if i := strings.IndexByte(p.cur, '='); i != -1 {
		si = i
	}
	if i := strings.IndexByte(p.cur, ';'); i != -1 {
		ei = i
	}

	return p.cur[si:ei]
}

func (p *Parser) jump() string {
	if i := strings.IndexByte(p.cur, ';'); i != -1 {
		return p.cur[i:]
	}

	return ""
}
