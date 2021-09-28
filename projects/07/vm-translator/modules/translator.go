package modules

import "fmt"

var memoryMap = map[string]string{
	"local":    "LCL",
	"argument": "ARG",
	"this":     "THIS",
	"that":     "THAT",
}

var baseAddrMap = map[string]int{
	"pointer": 3,
	"temp":    5,
}

type Translator struct {
	fileName       string
	translateCount int
}

func NewTranslator(fn string) *Translator {
	return &Translator{
		fileName: fn,
	}
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
M=M-D
`
	case "neg":
		return `@SP
A=M-1
M=-M
`
	case "eq":
		return t.conditionStatement("JEQ")
	case "gt":
		return t.conditionStatement("JGT")
	case "lt":
		return t.conditionStatement("JLT")
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
M=D|M`
	case "not":
		return `@SP
A=M-1
M=!M
`
	default:
		return ""
	}
}

func (t *Translator) TranslatePush(seg string, i int) string {
	switch seg {
	case "local", "argument", "this", "that":
		return fmt.Sprintf(`@%d
D=A
@%s
A=M
A=D+A
D=M
@SP
AM=M+1
A=A-1
M=D
`, i, memoryMap[seg])
	case "constant":
		return fmt.Sprintf(`@%d
D=A
@SP
AM=M+1
A=A-1
M=D
`, i)
	case "pointer", "temp":
		return fmt.Sprintf(`@%d
D=M
@SP
AM=M+1
A=A-1
M=D
`, baseAddrMap[seg]+i)
	case "static":
		return fmt.Sprintf(`@%s.%d
D=M
@SP
AM=M+1
A=A-1
M=D
`, t.fileName, i)
	default:
		return ""
	}
}

func (t *Translator) TranslatePop(seg string, i int) string {
	switch seg {
	case "local", "argument", "this", "that":
		return fmt.Sprintf(`@%s
D=M
@%d
D=D+A
@R13
M=D
@SP
AM=M-1
D=M
@R13
A=M
M=D
`, memoryMap[seg], i)
	case "pointer", "temp":
		return fmt.Sprintf(`@SP
AM=M-1
D=M
@%d
M=D
`, baseAddrMap[seg]+i)
	case "static":
		return fmt.Sprintf(`@SP
AM=M-1
D=M
@%s.%d
M=D
`, t.fileName, i)
	default:
		return ""
	}
}

func (t *Translator) conditionStatement(jmp string) string {
	t.translateCount++
	lb := fmt.Sprintf("%s_TRUE_%d", jmp, t.translateCount)
	return fmt.Sprintf(`@SP
AM=M-1
D=M
A=A-1
D=M-D
M=-1
@%s
D;%s
@SP
A=M-1
M=0
(%s)
`, lb, jmp, lb)
}
