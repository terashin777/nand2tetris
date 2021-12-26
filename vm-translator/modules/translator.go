package modules

import (
	"fmt"
	"strings"

	"github.com/terashin777/vm-translator/utils"
)

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
	fileName    string
	symbolCount int
}

func NewTranslator() *Translator {
	return &Translator{
		symbolCount: 0,
	}
}

func (t *Translator) SetFunctionName(n string) {
	t.fileName = n
}

func (t *Translator) TranslateInit() string {
	return fmt.Sprintf(`@256
D=A
@SP
M=D
%s`, t.translateCall("Sys.init", 0))
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
		return fmt.Sprintf(`@%s
D=M
@%d
A=D+A
D=M
%s
`, memoryMap[seg], i, t.push())
	case "constant":
		return fmt.Sprintf(`@%d
D=A
%s
`, i, t.push())
	case "pointer", "temp":
		return fmt.Sprintf(`@%d
D=M
%s
`, baseAddrMap[seg]+i, t.push())
	case "static":
		return fmt.Sprintf(`@%s.%d
D=M
%s
`, t.fileName, i, t.push())
	default:
		return ""
	}
}

func (t *Translator) push() string {
	return `@SP
AM=M+1
A=A-1
M=D`
}

func (t *Translator) TranslatePop(seg string, i int) string {
	if i == -1 {
		fmt.Print("ok")
	}
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

func (t *Translator) TranslateLabel(l string) string {
	t.validateLabel(l)
	return fmt.Sprintf(`(%s)
`, l)
}

func (t *Translator) validateLabel(l string) {
	if t.isStartWithNumber(l) {
		panic(fmt.Errorf("invalid format label: label is not allowed to start with number"))
	}
	if !t.isValidLabelChar(l) {
		panic(fmt.Errorf("invalid format label: allowed label character is alphabet, number, _, ., : or $"))
	}
}

func (t *Translator) isStartWithNumber(l string) bool {
	return utils.Char.IsNumber(rune(l[0]))
}

func (t *Translator) isValidLabelChar(l string) bool {
	for _, r := range l {
		if !t.isIncludeValidLabelCharOnly(r) {
			return false
		}
	}

	return true
}

func (t *Translator) isIncludeValidLabelCharOnly(r rune) bool {
	return utils.Char.IsAlphabet(r) || utils.Char.IsNumber(r) || strings.ContainsRune("_.:$", r)
}

func (t *Translator) TranslateGoto(l string) string {
	return fmt.Sprintf(`@%s
0;JMP
`, l)
}

func (t *Translator) TranslateIf(l string) string {
	return fmt.Sprintf(`@SP
M=M-1
@SP
A=M
D=M
@%s
D;JNE
`, l)
}

func (t *Translator) TranslateCall(f string, n int) string {
	return t.translateCall(f, n)
}

func (t *Translator) translateCall(f string, n int) string {
	t.symbolCount++
	ret := fmt.Sprintf("%s_RETURN%d", f, t.symbolCount)
	return fmt.Sprintf(`%[1]s
%[2]s 
%[3]s
%[4]s
%[5]s
@%[6]d
D=A
@5
D=D+A
@SP
D=M-D
@ARG
M=D
@SP
D=M
@LCL
M=D
@%[7]s
0;JMP
(%[8]s)
`,
		t.pushLabelAddress(ret),
		t.pushSegmentAddress("LCL"),
		t.pushSegmentAddress("ARG"),
		t.pushSegmentAddress("THIS"),
		t.pushSegmentAddress("THAT"),
		n,
		f,
		ret,
	)
}

func (t *Translator) pushLabelAddress(label string) string {
	return fmt.Sprintf(`@%s
D=A
%s`, label, t.push())
}

func (t *Translator) pushSegmentAddress(seg string) string {
	return fmt.Sprintf(`@%s
D=M
%s`, seg, t.push())
}

func (t *Translator) TranslateReturn() string {
	frame := "13"
	ret := "14"
	return fmt.Sprintf(`@LCL
D=M
@%[1]s
M=D
@5
A=D-A
D=M
@%[2]s
M=D
@SP
AM=M-1
D=M
@ARG
A=M
M=D
@ARG
D=M
@SP
M=D+1
%[3]s
%[4]s
%[5]s
%[6]s
@%[2]s
A=M
0;JMP
`,
		frame,
		ret,
		t.returnEachSegment("THAT", frame, 1),
		t.returnEachSegment("THIS", frame, 2),
		t.returnEachSegment("ARG", frame, 3),
		t.returnEachSegment("LCL", frame, 4),
	)
}

func (t *Translator) returnEachSegment(seg, base string, rel int) string {
	return fmt.Sprintf(`@%s
D=M
@%d
A=D-A
D=M
@%s
M=D`, base, rel, seg)
}

func (t *Translator) TranslateFunction(f string, k int) string {
	if k == 0 {
		return fmt.Sprintf(`(%s)
`, f)
	}

	loop := fmt.Sprintf("%s_LOOP", f)
	return fmt.Sprintf(`(%s)
@%d
D=A
(%[3]s)
@%[1]sEND
D;JEQ
@SP
AM=M+1
A=A-1
M=0
D=D-1
@%[3]s
0;JMP
(%[1]sEND)
`, f, k, loop)
}

func (t *Translator) conditionStatement(jmp string) string {
	t.symbolCount++
	lb := fmt.Sprintf("%s_TRUE%d", jmp, t.symbolCount)
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
