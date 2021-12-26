package modules

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/terashin777/jack-analyzer/models"
)

type JackTokenizer struct {
	ti        int
	curTokens []*models.Token
	r         io.ReadCloser
	s         *bufio.Scanner
	idChecker *regexp.Regexp
}

func NewJackTokenizer() *JackTokenizer {
	return &JackTokenizer{
		ti:        -1,
		idChecker: regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`),
	}
}

func (t *JackTokenizer) Open(f string) error {
	err := t.setReader(f)
	if err != nil {
		return err
	}

	return nil
}

func (t *JackTokenizer) Close() error {
	return t.r.Close()
}

func (t *JackTokenizer) Advance() error {
	if 0 <= t.ti && t.ti < len(t.curTokens)-1 {
		t.ti++
		return nil
	}

	if !t.s.Scan() {
		err := t.r.Close()
		if err != nil {
			return err
		}

		return io.EOF
	}

	t.initTokens()
	if len(t.curTokens) == 0 {
		return t.Advance()
	}
	return nil
}

func (t *JackTokenizer) initTokens() {
	t.curTokens = t.tokenizeLine(t.s.Text())
	t.ti = 0
}

func (t *JackTokenizer) setReader(f string) error {
	var err error
	t.r, err = os.Open(f)
	if err != nil {
		return err
	}

	t.s = bufio.NewScanner(t.r)
	t.Advance()
	return nil
}

func (t *JackTokenizer) removeComment(l string) string {
	i := strings.Index(l, "//")
	if i >= 0 {
		return l[:i]
	}

	return t.removeClosedComment(l)
}

func (t *JackTokenizer) removeClosedComment(l string) string {
	i := strings.Index(l, "/*")
	if i == -1 {
		return l
	}

	start := i + 2
	if l[start] == '*' {
		start++
	}

	end := strings.Index(string(l[start:]), "*/")
	lastIdx := start + end + 2
	if len(l)-1 < lastIdx {
		return t.removeClosedComment(string(l[:i]))
	}
	return t.removeClosedComment(string(l[:i]) + string(l[lastIdx:]))
}

func (t *JackTokenizer) tokenizeLine(l string) []*models.Token {
	l = t.removeComment(strings.TrimSpace(l))
	if l == "" {
		return []*models.Token{}
	}

	toks := []*models.Token{}
	rL := []rune(l)

	i := 0
	tmp := []rune{}
	for {
		if i >= len(rL) {
			break
		}

		r := rL[i]
		// 空白で分割できないトークンは別処理する↓
		// 定数文字列をトークンとして追加
		if r == '"' {
			end := t.stringValEndIndex(i, rL)
			toks = append(toks, models.NewToken(string(rL[i+1:end]), models.STRING_CONST))

			tmp = []rune{}
			// インデックス位置をスキップする
			i = end + 1
			continue
		}
		// シンボルを他のトークンと切り離して追加
		if t.IsSymbol(string(r)) {
			s := string(tmp)
			if s != "" {
				// symbol直前までのtokenを追加
				toks = append(toks, t.newToken(s))
			}
			// symbolを追加
			toks = append(toks, t.newToken(string(r)))

			tmp = []rune{}
			i++
			continue
		}

		// 空白で分割できるトークン
		if r == ' ' {
			s := string(tmp)
			if s != "" {
				// 空白直前までのtokenを追加
				toks = append(toks, t.newToken(s))
			}

			tmp = []rune{}
			i++
			continue
		}

		tmp = append(tmp, r)
		i++
	}

	// 最後残った部分をトークンとして追加
	s := string(tmp)
	if s != "" {
		toks = append(toks, t.newToken(s))
	}

	return toks
}

func (t *JackTokenizer) newToken(s string) *models.Token {
	return models.NewToken(s, t.tokenType(s))
}

func (t *JackTokenizer) stringValEndIndex(i int, l []rune) int {
	// 次の"のインデックスを取得
	s := string(l[i+1:])
	end := strings.IndexRune(s, '"')

	// 元の文字列での"のインデックスに直す（最初の"の分と次のインデックスに進めるので、+2する）
	return i + end + 1
}

func (t *JackTokenizer) CurToken() *models.Token {
	return t.curTokens[t.ti]
}

func (t *JackTokenizer) TokenValue() string {
	return t.curTokens[t.ti].Value
}

func (t *JackTokenizer) TokenType() models.TokenType {
	return t.curTokens[t.ti].Type
}

func (t *JackTokenizer) tokenType(tok string) models.TokenType {
	if t.IsKeyword(tok) {
		return models.KEYWORD
	}
	if t.IsSymbol(tok) {
		return models.SYMBOL
	}
	if t.IsIntVal(tok) {
		return models.INT_CONST
	}
	if t.IsStringVal(tok) {
		return models.STRING_CONST
	}
	if t.IsIdentifier(tok) {
		return models.IDENTIFIER
	}

	return models.UNKNOWN_TOKEN
}

func (t *JackTokenizer) IsKeyword(tok string) bool {
	_, ok := models.Keywords[tok]
	return ok
}

func (t *JackTokenizer) IsSymbol(tok string) bool {
	_, ok := models.Symbols[tok]
	return ok
}

func (t *JackTokenizer) IsIdentifier(tok string) bool {
	return t.idChecker.MatchString(tok)
}

func (t *JackTokenizer) IsIntVal(tok string) bool {
	_, err := strconv.ParseInt(tok, 10, 16)
	return err == nil
}

func (t *JackTokenizer) IsStringVal(tok string) bool {
	return strings.HasPrefix(tok, "\"") && strings.HasSuffix(tok, "\"") && utf8.ValidString(tok)
}

func (t *JackTokenizer) Keyword() models.KeywordType {
	return models.Keywords[t.CurToken().Value]
}

func (t *JackTokenizer) Symbol() string {
	return t.CurToken().Value
}

func (t *JackTokenizer) Identifier() string {
	return t.CurToken().Value
}

func (t *JackTokenizer) IntVal() int {
	i, _ := strconv.ParseInt(t.CurToken().Value, 10, 16)
	return int(i)
}

func (t *JackTokenizer) StringVal() string {
	return t.CurToken().Value
}
