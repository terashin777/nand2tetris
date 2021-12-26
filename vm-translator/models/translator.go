package models

type ITranslator interface {
	TranslateInit() string
	TranslateArithmetic(c string) string
	TranslatePush(seg string, i int) string
	TranslatePop(seg string, i int) string
	TranslateLabel(l string) string
	TranslateGoto(l string) string
	TranslateIf(l string) string
	TranslateCall(fn string, n int) string
	TranslateReturn() string
	TranslateFunction(fn string, n int) string
	SetFunctionName(n string)
}
