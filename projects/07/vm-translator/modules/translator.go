package modules

import "fmt"

type Translator struct {
}

func (t *Translator) TranslateArithmetic(c string) string {
	switch c {
	case "add":
		return `@SP
AM=M-1
D=M
A=A-1
M=D+M
`
	case "sub":
		return `@SP
AM=M-1
D=M
A=A-1
D=-D
M=D+M
`
	case "neg":
		return `@SP
A=M-1
M=-M
`
	case "eq":
		return t.conditionStatement("JEQ")
	case "gt":
		return t.conditionStatement("JLT")
	case "lt":
		return t.conditionStatement("JGT")
	case "and":
		return `@SP
AM=M-1
D=M
A=A-1
M=D&M
`
	case "or":
		return `@SP
AM=M-1
D=M
A=A-1
M=D|M
`
	case "not":
		return `@SP
A=M-1
M=!M
`
	default:
		return ""
	}
}

func (t *Translator) conditionStatement(jmp string) string {
	return fmt.Sprintf(`@SP
AM=A-1
D=M
A=A-1
D=D-M
@TRUE
D;%s
@SP
A=M-1
M=0
(TRUE)
@SP
A=M-1
M=-1
`, jmp)
}
