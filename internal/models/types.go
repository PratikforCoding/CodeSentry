package models

type TokenType int

const (
	KEYWORD TokenType = iota
	IDENTIFIER
	OPERATOR
	LITERAL
	COMMENT
	DELIMITER
	UNKNOWN
)

type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

type Language string

const (
	Go         Language = "go"
	Java       Language = "java"
	JavaScript Language = "javascript"
	Python     Language = "python"
	SQL        Language = "sql"
	Unknown    Language = "unknown"
)
