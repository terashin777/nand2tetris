package modules

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/terashin777/jack-analyzer/models"
	"github.com/terashin777/jack-analyzer/utils"
)

var (
	predefinedClasses = []string{
		"Array",
	}
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
		classNames: predefinedClasses,
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
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeClassName()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("{"))
	if err != nil {
		return err
	}

	for {
		k := e.t.Keyword()
		if k.IsSubRoutine() {
			err = e.CompileSubroutine()
			if err != nil {
				return err
			}
		} else if k.IsClassVarDec() {
			err = e.CompileClassVarDec()
			if err != nil {
				return err
			}
		} else {
			break
		}
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("}"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) writeClassName() error {
	e.classNames = append(e.classNames, e.t.CurToken().Value)
	err := e.writeTokenWithValidate(e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileClassVarDec() error {
	return e.wrap("classVarDec", e.compileClassVarDec)
}

func (e *CompilationEngine) compileClassVarDec() error {
	err := e.writeToken()
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
	err := e.writeTokenWithValidate(e.makeSymbolValidator(","))
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileSubroutine() error {
	return e.wrap("subrouneDec", e.compileSubroutine)
}

func (e *CompilationEngine) compileSubroutine() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.compileType()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("("))
	if err != nil {
		return err
	}

	err = e.CompileParameterList()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator(")"))
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
	err := e.writeTokenWithValidate(e.makeSymbolValidator("{"))
	if err != nil {
		return err
	}

	for e.t.CurToken().IsKeywordOf(models.VAR) {
		err = e.CompileVarDec()
		if err != nil {
			return err
		}
	}

	err = e.CompileStatements()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("}"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileVarDec() error {
	return e.wrap("varDec", e.compileVarDec)
}

func (e *CompilationEngine) compileVarDec() error {
	err := e.writeTokenWithValidate(e.makeKeywordValidator(models.VAR))
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

	err = e.writeTokenWithValidate(e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	for !e.t.CurToken().IsSymbolOf(";") {
		err := e.writeVarNames()
		if err != nil {
			return err
		}
	}

	err = e.writeToken()
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileStatements() error {
	return e.wrap("statements", e.compileStatements)
}

func (e *CompilationEngine) compileStatements() error {
	for {
		switch e.t.Keyword() {
		case models.LET:
			err := e.CompileLet()
			if err != nil {
				return err
			}
			continue
		case models.IF:
			err := e.CompileIf()
			if err != nil {
				return err
			}
			continue
		case models.WHILE:
			err := e.CompileWhile()
			if err != nil {
				return err
			}
			continue
		case models.DO:
			err := e.CompileDo()
			if err != nil {
				return err
			}
			continue
		case models.RETURN:
			err := e.CompileReturn()
			if err != nil {
				return err
			}
			continue
		}

		return nil
	}
}

func (e *CompilationEngine) CompileLet() error {
	return e.wrap("letStatement", e.compileLet)
}

func (e *CompilationEngine) compileLet() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	err = e.compileLetArrayExpression()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("="))
	if err != nil {
		return err
	}

	err = e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator(";"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) compileLetArrayExpression() error {
	if e.t.CurToken().IsSymbolOf("[") {
		err := e.writeToken()
		if err != nil {
			return err
		}

		// TODO: expression
		err = e.writeToken()
		if err != nil {
			return err
		}

		err = e.writeTokenWithValidate(e.makeSymbolValidator("]"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *CompilationEngine) CompileIf() error {
	return e.wrap("ifStatement", e.compileIf)
}

func (e *CompilationEngine) compileIf() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("("))
	if err != nil {
		return err
	}

	// TODO: expression
	err = e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator(")"))
	if err != nil {
		return err
	}

	err = e.compileIfBody()
	if err != nil {
		return err
	}

	if e.t.CurToken().IsKeywordOf(models.ELSE) {
		err = e.compileIfElse()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *CompilationEngine) compileIfElse() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.compileIfBody()
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) compileIfBody() error {
	err := e.writeTokenWithValidate(e.makeSymbolValidator("{"))
	if err != nil {
		return err
	}
	err = e.CompileStatements()
	if err != nil {
		return err
	}
	err = e.writeTokenWithValidate(e.makeSymbolValidator("}"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileWhile() error {
	return e.wrap("whileStatement", e.compileWhile)
}

func (e *CompilationEngine) compileWhile() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("("))
	if err != nil {
		return err
	}

	// TODO: expression
	err = e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator(")"))
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("{"))
	if err != nil {
		return err
	}

	err = e.compileStatements()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeSymbolValidator("}"))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileDo() error {
	return e.wrap("doStatement", e.compileDo)
}

func (e *CompilationEngine) compileDo() error {
	// TODO: subroutineCall
	err := e.writeToken()
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) CompileReturn() error {
	return e.wrap("returnStatement", e.compileReturn)
}

func (e *CompilationEngine) compileReturn() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	if !e.t.CurToken().IsSymbolOf(";") {
		// TODO: expression
		err = e.writeToken()
		if err != nil {
			return err
		}
	}

	err = e.writeToken()
	if err != nil {
		return err
	}
	return nil
}

func (e *CompilationEngine) compileType() error {
	return e.writeToken()

	t := e.t.CurToken()
	// TODO: 同一フォルダ内でクラスが定義されているかチェックできるようにしたい
	if t.IsPrimitiveType() || e.isClassDefined(t.Value) {
		return e.writeToken()
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
	for !e.t.CurToken().IsSymbolOf(")") {
		err := e.writeParameterList()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *CompilationEngine) writeParameterList() error {
	err := e.writeToken()
	if err != nil {
		return err
	}

	err = e.writeTokenWithValidate(e.makeTokenTypeValidator(models.IDENTIFIER))
	if err != nil {
		return err
	}

	return nil
}

func (e *CompilationEngine) makeSymbolValidator(s string) func(t *models.Token) error {
	return func(t *models.Token) error {
		if t.IsSymbolOf(s) {
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

func (e *CompilationEngine) makeTokenTypeValidator(tt models.TokenType) func(t *models.Token) error {
	return func(t *models.Token) error {
		if t.Type == tt {
			return nil
		}

		return e.invalidTokenTypeError(t.Type)
	}
}

func (e *CompilationEngine) invalidTypeError(v string) error {
	return fmt.Errorf("'%s' is not type or not defined\n %s", v, utils.Stack(2))
}

func (e *CompilationEngine) invalidTokenTypeError(k models.TokenType) error {
	return fmt.Errorf("token is not token type '%s'\n %s", models.TokenTypeString[k], utils.Stack(2))
}

func (e *CompilationEngine) invalidSymbolError(v string) error {
	return fmt.Errorf("token is not symbol '%s'\n %s", v, utils.Stack(2))
}

func (e *CompilationEngine) invalidKeywordError(v string, k models.KeywordType) error {
	return fmt.Errorf("token '%s' is not keyword '%s'\n %s", v, models.KeywordTypeString[k], utils.Stack(2))
}

func (e *CompilationEngine) writeTokenWithValidate(validate func(t *models.Token) error) error {
	if err := validate(e.t.CurToken()); err != nil {
		return err
	}

	if err := e.writeToken(); err != nil {
		return err
	}
	return nil
}

func (e *CompilationEngine) writeToken() error {
	_, err := e.w.Write([]byte(e.indentString() + models.WrapTokenTag(e.t.CurToken()) + "\n"))
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
