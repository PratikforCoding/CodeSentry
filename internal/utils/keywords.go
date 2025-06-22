package utils

import "strings"

func (p *Parser) isKeyword(word string) bool {
	keywords := []string{
		// Go keywords (complete set)
		"break", "case", "chan", "const", "continue", "default", "defer", "else",
		"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
		"map", "package", "range", "return", "select", "struct", "switch", "type", "var",
		// Go built-in types and functions
		"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "complex64", "complex128", "byte", "rune", "string", "bool",
		"error", "make", "len", "cap", "new", "append", "copy", "delete", "panic", "recover",
		"close", "nil", "iota",

		// JavaScript keywords (complete set)
		"async", "await", "class", "extends", "function", "let", "const", "var",
		"try", "catch", "finally", "throw", "new", "this", "super", "static",
		"get", "set", "of", "in", "instanceof", "typeof", "void", "delete",
		"export", "import", "from", "as", "default", "with", "debugger",
		// JavaScript built-ins and common identifiers
		"console", "window", "document", "Array", "Object", "String", "Number",
		"Boolean", "Date", "RegExp", "Math", "JSON", "Promise", "Symbol",
		"undefined", "NaN", "Infinity", "prototype", "constructor",

		// Java keywords (complete set)
		"abstract", "assert", "boolean", "byte", "case", "catch", "char", "class",
		"const", "continue", "default", "do", "double", "else", "enum", "extends",
		"final", "finally", "float", "for", "goto", "if", "implements", "import",
		"instanceof", "int", "interface", "long", "native", "new", "package",
		"private", "protected", "public", "return", "short", "static", "strictfp",
		"super", "switch", "synchronized", "this", "throw", "throws", "transient",
		"try", "void", "volatile", "while",
		// Java built-ins and common classes
		"String", "System", "Object", "Integer", "Double", "Boolean", "Character",
		"Math", "Thread", "Exception", "ArrayList", "HashMap", "List", "Map", "Set",

		// Python keywords (complete set)
		"False", "None", "True", "and", "as", "assert", "async", "await", "break",
		"class", "continue", "def", "del", "elif", "else", "except", "finally",
		"for", "from", "global", "if", "import", "in", "is", "lambda", "nonlocal",
		"not", "or", "pass", "raise", "return", "try", "while", "with", "yield",
		// Python built-ins
		"int", "float", "str", "list", "dict", "tuple", "set", "bool", "bytes",
		"len", "range", "enumerate", "zip", "map", "filter", "sorted", "sum",
		"min", "max", "abs", "all", "any", "print", "input", "open", "type",
		"isinstance", "hasattr", "getattr", "setattr", "__init__", "__str__", "__name__",

		// SQL keywords (common across databases)
		"SELECT", "FROM", "WHERE", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP",
		"ALTER", "TABLE", "INDEX", "VIEW", "DATABASE", "SCHEMA", "PROCEDURE", "FUNCTION",
		"TRIGGER", "JOIN", "INNER", "LEFT", "RIGHT", "FULL", "OUTER", "ON", "USING",
		"GROUP", "BY", "HAVING", "ORDER", "ASC", "DESC", "LIMIT", "OFFSET", "UNION",
		"INTERSECT", "EXCEPT", "ALL", "DISTINCT", "TOP", "INTO", "VALUES", "SET",
		"AND", "OR", "NOT", "IN", "EXISTS", "BETWEEN", "LIKE", "IS", "NULL",
		"PRIMARY", "KEY", "FOREIGN", "REFERENCES", "UNIQUE", "CHECK", "DEFAULT",
		"AUTO_INCREMENT", "IDENTITY", "SERIAL", "CONSTRAINT", "CASCADE", "RESTRICT",
		"GRANT", "REVOKE", "COMMIT", "ROLLBACK", "TRANSACTION", "BEGIN", "END",
		"DECLARE", "CURSOR", "FETCH", "CLOSE", "DEALLOCATE", "EXEC", "EXECUTE",
		// SQL data types
		"VARCHAR", "CHAR", "TEXT", "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT",
		"DECIMAL", "NUMERIC", "FLOAT", "DOUBLE", "REAL", "BIT", "BOOLEAN", "BOOL",
		"DATE", "TIME", "DATETIME", "TIMESTAMP", "YEAR", "BLOB", "CLOB", "BINARY",

		// Common keywords across all languages
		"if", "else", "for", "while", "do", "switch", "case", "default", "break",
		"continue", "return", "true", "false", "null", "class", "function", "var",
		"let", "const", "try", "catch", "throw", "new", "this", "super", "static",
		"public", "private", "protected", "abstract", "final", "interface", "enum",
		"extends", "implements", "import", "export", "package", "namespace", "using",
		"include", "require", "module", "async", "await", "yield", "lambda", "def",
		"end", "begin", "then", "when", "unless", "until", "loop", "match", "case",
	}

	for _, keyword := range keywords {
		if strings.EqualFold(word, keyword) {
			return true
		}
	}
	return false
}
