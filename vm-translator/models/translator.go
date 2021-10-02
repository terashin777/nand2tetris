package models

type ITranslator interface {
	TranslateArithmetic(c string) string
	TranslatePush(seg string, i int) string
	TranslatePop(seg string, i int) string
	TranslateLabel(l string) string
	TranslateGoto(l string) string
	TranslateIf(l string) string
}
