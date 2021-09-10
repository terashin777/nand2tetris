package modules

import (
	"bufio"
	"fmt"
	"io"
	"math/big"
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
	next := p.s.Scan()
	if !next {
		err := p.s.Err()
		if err != nil {
			return "", err
		}

		return "", io.EOF
	}

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
	if ct == A_COMMAND || ct == L_COMMAND {
		sym = p.symbol()
		v, err := strconv.ParseInt(sym, 10, 16)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%016s", strconv.FormatInt(v, 2)), nil
	}

	if ct == C_COMMAND {
		return p.parseCCommand()
	}

	return "", nil
}

func (p *Parser) parseCCommand() (string, error) {
	d, err := p.c.Dest(p.dest())
	if err != nil {
		return "", err
	}
	id := big.NewInt(int64(d))
	id = id.Lsh(id, 3)

	c, err := p.c.Comp(p.comp())
	if err != nil {
		return "", err
	}
	ic := big.NewInt(int64(c))
	ic = ic.Lsh(ic, 6)

	j, err := p.c.Jump(p.jump())
	if err != nil {
		return "", err
	}
	ij := big.NewInt(int64(j))

	odc := ic.Or(ic, id)
	o := odc.Or(odc, ij)

	pre, _ := strconv.ParseUint("1110000000000000", 2, 16)
	ipre := big.NewInt(int64(pre))
	res := ipre.Or(ipre, o)
	return fmt.Sprintf("%016s", strconv.FormatInt(res.Int64(), 2)), nil
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
		si = i + 1
	}
	if i := strings.IndexByte(p.cur, ';'); i != -1 {
		ei = i
	}

	return p.cur[si:ei]
}

func (p *Parser) jump() string {
	if i := strings.IndexByte(p.cur, ';'); i != -1 {
		return p.cur[i+1:]
	}

	return ""
}
