package modules

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/terashin777/jack-analyzer/models"
	"github.com/terashin777/jack-analyzer/utils"
)

type CompilationEngine struct {
	t          *JackTokenizer
	w          models.WriteSeekCloser
	f          string
	indent     int
	indentUnit int
	classNames []string
}

func NewCompilationEngine(src, dest string, t *JackTokenizer) (*CompilationEngine, error) {
	return &CompilationEngine{
		t:          t,
		indent:     0,
		indentUnit: 2,
		classNames: []string{},
	}, nil
}

func (e *CompilationEngine) SetFile(f string) {
	e.f = f
}

func (e *CompilationEngine) Open() error {
	err := e.t.Open(e.f)
	if err != nil {
		return err
	}

	nf, err := utils.MakeSameFileName(e.f, "xml")
	if err != nil {
		return err
	}
	w, err := os.Create(nf)
	if err != nil {
		return err
	}
	e.w = w

	return nil
}

func (e *CompilationEngine) Close() error {
	var errs []error
	err := e.w.Close()
	if err != nil {
		errs = append(errs, err)
	}

	err = e.t.Close()
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf(fmt.Sprintf("%s", errs))
	}

	return nil
}

func (e *CompilationEngine) CompileClass() error {
	defer func() {
		e.Close()
	}()

	k := e.t.Keyword()
	if !k.IsClass() {
		return e.invalidKeywordError(e.t.CurToken().Value, k)
	}

	return e.wrap("class", e.compileClass)
}

func (e *CompilationEngine) compileClass() error {
	err := e.writeToken(e.t.CurToken())
	if err != nil {
		return err
	}

	err = e.writeClassName()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeSymbolValidator("{"))
	if err != nil {
		return err
	}

	for {
		k := e.t.Keyword()
		if k.IsSubRoutine() {
			e.CompileSubroutine()
		} else if k.IsClassVarDec() {
			e.CompileClassVarDec()
		} else {
			break
		}
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeSymbolValidator("}"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) writeClassName() error {
	t := e.t.CurToken()
	e.classNames = append(e.classNames, t.Value)
	err := e.writeTokenWithValidate(t, e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileClassVarDec() error {
	return e.wrap("classVarDec", e.compileClassVarDec)
}

func (e *CompilationEngine) compileClassVarDec() error {
	err := e.writeToken(e.t.CurToken())
	if err != nil {
		return err
	}

	err = e.compileVarDecBody()
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) writeVarNames() error {
	err := e.writeTokenWithValidate(e.t.CurToken(), e.makeSymbolValidator(","))
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileSubroutine() error {
	return e.wrap("subrouneDec", e.compileSubroutine)
}

func (e *CompilationEngine) compileSubroutine() error {
	err := e.writeToken(e.t.CurToken())
	if err != nil {
		return err
	}

	err = e.compileType()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeSymbolValidator("("))
	if err != nil {
		return err
	}

	err = e.CompileParameterList()
	if err != nil {
		return err
	}

	err = e.CompileSubroutineBody()
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileSubroutineBody() error {
	return e.wrap("subroutineBody", e.compileSubroutineBody)
}

func (e *CompilationEngine) compileSubroutineBody() error {
	err := e.writeTokenWithValidate(e.t.CurToken(), e.makeSymbolValidator("{"))
	if err != nil {
		return err
	}

	err = e.CompileVarDec()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeSymbolValidator("}"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileVarDec() error {
	return e.wrap("varDec", e.compileVarDec)
}

func (e *CompilationEngine) compileVarDec() error {
	err := e.writeTokenWithValidate(e.t.CurToken(), e.makeKeywordValidator(models.VAR))
	if err != nil {
		return err
	}

	err = e.compileVarDecBody()
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) compileVarDecBody() error {
	err := e.compileType()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}
	for !e.isSymbolOf(e.t.CurToken(), ";") {
		err := e.writeVarNames()
		if err != nil {
			return err
		}
	}

	err = e.writeToken(e.t.CurToken())
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileStatements() error {
	return e.wrap("statements", e.compileVarDec)
}

func (e *CompilationEngine) compileType() error {
	t := e.t.CurToken()
	if t.IsPrimitiveType() || e.isClassDefined(t.Value) {
		return nil
	}

	return e.invalidTypeError(t.Value)
}

func (e *CompilationEngine) isClassDefined(n string) bool {
	for _, cn := range e.classNames {
		if cn == n {
			return true
		}
	}

	return false
}

func (e *CompilationEngine) CompileParameterList() error {
	return e.wrap("parameterList", e.compileParameterList)
}

func (e *CompilationEngine) compileParameterList() error {
	for !e.isSymbolOf(e.t.CurToken(), ")") {
		err := e.writeParameterList()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *CompilationEngine) writeParameterList() error {
	err := e.writeToken(e.t.CurToken())
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.t.CurToken(), e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) makeSymbolValidator(s string) func(t *models.Token) error {
	return func(t *models.Token) error {
		if e.isSymbolOf(t, s) {
			return nil
		}

		return e.invalidSymbolError(t.Value)
	}
}

func (e *CompilationEngine) makeKeywordValidator(k models.KeywordType) func(t *models.Token) error {
	return func(t *models.Token) error {
		if t.IsKeywordOf(k) {
			return nil
		}

		return e.invalidKeywordError(t.Value, k)
	}
}

func (e *CompilationEngine) isSymbolOf(t *models.Token, s string) bool {
	_, ok := models.Symbols[t.Value]
	if ok && t.Type == models.SYMBOL && t.Value == s {
		return true
	}

	return false
}

func (e *CompilationEngine) makeTokenTypeValidator(tt models.TokenType) func(t *models.Token) error {
	return func(t *models.Token) error {
		if t.Type == tt {
			return nil
		}

		return e.invalidTokenTypeError(t.Type)
	}
}

func (e *CompilationEngine) invalidTypeError(v string) error {
	return fmt.Errorf("'%s' is not type or not defined", v)
}

func (e *CompilationEngine) invalidTokenTypeError(k models.TokenType) error {
	return fmt.Errorf("token is not token type '%s'", models.TokenTypeString[k])
}

func (e *CompilationEngine) invalidSymbolError(v string) error {
	return fmt.Errorf("token is not symbol '%s'", v)
}

func (e *CompilationEngine) invalidKeywordError(v string, k models.KeywordType) error {
	return fmt.Errorf("token '%s' is not keyword '%s'", v, models.KeywordTypeString[k])
}

func (e *CompilationEngine) isClass(t models.KeywordType) bool {
	return t == models.CLASS
}

func (e *CompilationEngine) isSubRoutine(t models.KeywordType) bool {
	return t == models.CONSTRUCTOR || t == models.FUNCTION || t == models.METHOD
}

func (e *CompilationEngine) isClassVarDec(t models.KeywordType) bool {
	return t == models.STATIC || t == models.FIELD
}

func (e *CompilationEngine) writeTokenWithValidate(t *models.Token, validate func(t *models.Token) error) error {
	if err := validate(t); err != nil {
		return err
	}

	if err := e.writeToken(t); err != nil {
		return err
	}
	return nil
}

func (e *CompilationEngine) writeToken(t *models.Token) error {
	_, err := e.w.Write([]byte(e.indentString() + models.WrapTokenTag(t) + "\n"))
	if err != nil {
		return err
	}

	if err := e.t.Advance(); err != nil {
		return err
	}
	return nil
}

func (e *CompilationEngine) indentString() string {
	var b strings.Builder
	totalIndent := e.indent * e.indentUnit
	for i := 0; i < totalIndent; i++ {
		b.WriteString(" ")
	}

	return b.String()
}

func (e *CompilationEngine) wrap(tag string, write func() error) error {
	_, err := e.w.Write([]byte(fmt.Sprintf("<%s>\n", tag)))
	if err != nil {
		return err
	}
	e.indent++
	err = write()
	if err != nil && err != io.EOF {
		return err
	}
	_, err = e.w.Write([]byte(fmt.Sprintf("</%s>\n", tag)))
	if err != nil {
		return err
	}
	e.indent--

	return nil
}
