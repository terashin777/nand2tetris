package models

type KeywordType int

const (
	UNKNOWN_KEYWORD KeywordType = iota
	CLASS
	CONSTRUCTOR
	FUNCTION
	METHOD
	FIELD
	STATIC
	VAR
	INT
	CHAR
	BOOLEAN
	VOID
	TRUE
	FALSE
	NULL
	THIS
	LET
	DO
	IF
	ELSE
	WHILE
	RETURN
)

var Keywords map[string]KeywordType = map[string]KeywordType{
	"class":       CLASS,
	"constructor": CONSTRUCTOR,
	"function":    FUNCTION,
	"method":      METHOD,
	"field":       FIELD,
	"static":      STATIC,
	"var":         VAR,
	"int":         INT,
	"char":        CHAR,
	"boolean":     BOOLEAN,
	"void":        VOID,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
	"this":        THIS,
	"let":         LET,
	"do":          DO,
	"if":          IF,
	"else":        ELSE,
	"while":       WHILE,
	"return":      RETURN,
}

var KeywordTypeString = map[KeywordType]string{}

func init() {
	for k, v := range Keywords {
		KeywordTypeString[v] = k
	}
}

func (t KeywordType) IsClass() bool {
	return t == CLASS
}

func (t KeywordType) IsSubRoutine() bool {
	return t == CONSTRUCTOR || t == FUNCTION || t == METHOD
}

func (t KeywordType) IsClassVarDec() bool {
	return t == STATIC || t == FIELD
}

func (t KeywordType) IsPrimitiveType() bool {
	return t == VOID || t == INT || t == CHAR || t == BOOLEAN
}

func (t KeywordType) IsVar() bool {
	return t == VAR
}
