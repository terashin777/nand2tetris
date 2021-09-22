package modules

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/terashin777/assembler/table"
)

type (
	Parser struct {
		c    *Code
		s    *bufio.Scanner
		t    map[string]uint16
		ra   uint16
		cur  string
		src  []byte
		size int64
	}

	CommandType int
)

const (
	A_COMMAND CommandType = iota
	C_COMMAND
	L_COMMAND
)

func NewParser(r io.ReadSeeker, size int64) (*Parser, error) {
	p := &Parser{
		s:    bufio.NewScanner(r),
		c:    &Code{},
		ra:   16,
		t:    table.DefinedSymbolTable,
		src:  []byte{},
		size: size,
	}

	err := p.prepare()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) prepare() (err error) {
	defer func() {
		err = p.donePrepare()
	}()

	w := bytes.NewBuffer(make([]byte, 0, p.size))
	var c uint16 = 0
	for p.s.Scan() {
		com := p.extractCommand()
		if com == "" {
			continue
		}
		p.cur = com

		if p.isLCommand() {
			p.t[p.symbol()] = c
			continue
		}

		w.Write([]byte(com + "\n"))
		c++
	}

	p.src = w.Bytes()
	return p.s.Err()
}

func (p *Parser) extractCommand() string {
	l := p.s.Text()
	i := strings.Index(l, "//")
	if i == -1 {
		return strings.TrimSpace(l)
	}

	return strings.TrimSpace(l[0:i])
}

func (p *Parser) donePrepare() error {
	p.cur = ""
	p.s = bufio.NewScanner(bytes.NewReader(p.src))
	return nil
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

	com := p.s.Text()
	if com == "" {
		return "", nil
	}
	p.cur = com

	return p.parse()
}

// func (p *Parser) readCommand() string {
// 	l := p.s.Text()
// 	com := strings.TrimSpace(l)
// 	if p.isComment(com) {
// 		return ""
// 	}

// 	return com
// }

// func (p *Parser) isComment(l string) bool {
// 	return strings.HasPrefix(l, "//")
// }

func (p *Parser) parse() (string, error) {
	ct := p.commandType()

	if ct == A_COMMAND || ct == L_COMMAND {
		return p.resolveSymbol()
	}

	if ct == C_COMMAND {
		return p.parseCCommand()
	}

	return "", nil
}

func (p *Parser) resolveSymbol() (string, error) {
	sym := p.symbol()
	if sym == "" {
		return "", fmt.Errorf("symbol is blank")
	}
	v, err := strconv.ParseInt(sym, 10, 16)
	if err == nil {
		return fmt.Sprintf("%016s", strconv.FormatInt(v, 2)), nil
	}

	l, ok := p.t[sym]
	if ok {
		return fmt.Sprintf("%016s", strconv.FormatInt(int64(l), 2)), nil
	}

	return p.setRamAddress(sym), nil
}

func (p *Parser) setRamAddress(sym string) string {
	p.t[sym] = p.ra
	p.ra++
	return fmt.Sprintf("%016s", strconv.FormatInt(int64(p.t[sym]), 2))
}

func (p *Parser) parseCCommand() (string, error) {
	d, err := p.c.Dest(p.dest())
	if err != nil {
		return "", err
	}
	d16 := uint16(d) << 3

	c, err := p.c.Comp(p.comp())
	if err != nil {
		return "", err
	}
	c16 := uint16(c) << 6

	j, err := p.c.Jump(p.jump())
	if err != nil {
		return "", err
	}
	j16 := uint16(j)

	pre, _ := strconv.ParseUint("1110000000000000", 2, 16)
	res := uint16(pre) | c16 | d16 | j16
	return fmt.Sprintf("%016s", strconv.FormatInt(int64(res), 2)), nil
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
