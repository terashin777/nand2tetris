package models

type ITranslator interface {
	TranslateArithmetic(c string) string
}
