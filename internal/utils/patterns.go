package utils

import (
	"github.com/PratikforCoding/CodeSentry/internal/models"
	"regexp"
)

var SecurityPatterns = map[string]*regexp.Regexp{
	// SQL Injection patterns
	"sql_injection":     regexp.MustCompile(`(?i)(select|insert|update|delete|drop|create|alter)(\s+.*(\+|concat|format).*['"])?`),
	"sql_dynamic_query": regexp.MustCompile(`(?i)(query|execute|exec)\s*\(\s*['"].*\+.*['"]`),
	"sql_string_concat": regexp.MustCompile(`(?i)(select|insert|update|delete)\s+.*\+\s*.*['"]`),

	// XSS patterns
	"xss_vulnerability":      regexp.MustCompile(`(?i)(innerHTML|outerHTML|document\.write|insertAdjacentHTML)\s*[=+].*\+`),
	"xss_template_injection": regexp.MustCompile(`(?i)\$\{.*\+.*\}|<%.*\+.*%>`),
	"xss_eval_content":       regexp.MustCompile(`(?i)(innerHTML|outerHTML)\s*=.*eval\s*\(`),

	// Path traversal
	"path_traversal":    regexp.MustCompile(`\.\./|\.\.\\|%2e%2e%2f|%2e%2e\\`),
	"path_manipulation": regexp.MustCompile(`(?i)(file|path|dir)\s*\+\s*.*input`),
	"file_include":      regexp.MustCompile(`(?i)(include|require)\s*\(.*\$`),

	// Hardcoded secrets
	"hardcoded_password":   regexp.MustCompile(`(?i)(password|pwd|pass|secret)\s*[:=]\s*['"][^'"]{3,}['"]`),
	"hardcoded_api_key":    regexp.MustCompile(`(?i)(api_key|apikey|secret|token|key)\s*[:=]\s*['"][a-zA-Z0-9]{10,}['"]`),
	"hardcoded_connection": regexp.MustCompile(`(?i)(connection|conn|url)\s*[:=]\s*['"].*://.*['"]`),
	"hardcoded_crypto_key": regexp.MustCompile(`(?i)(private_key|secret_key|encryption_key|aes_key)\s*[:=]\s*['"][^'"]+['"]`),

	// Code execution patterns
	"eval_usage":        regexp.MustCompile(`(?i)\b(eval|exec|system|shell_exec|popen|subprocess)\s*\(`),
	"command_injection": regexp.MustCompile(`(?i)(system|exec|shell_exec|passthru|popen|os\.system)\s*\(.*(\+|\$|%s)`),
	"subprocess_shell":  regexp.MustCompile(`(?i)subprocess\.(call|run|Popen).*shell\s*=\s*True`),
	"runtime_exec":      regexp.MustCompile(`(?i)Runtime\.getRuntime\(\)\.exec\s*\(`),

	// Deserialization vulnerabilities
	"unsafe_deserialization": regexp.MustCompile(`(?i)(pickle\.loads|yaml\.load|json\.loads|unserialize|readObject|ObjectInputStream).*input`),
	"xml_external_entity":    regexp.MustCompile(`(?i)(DocumentBuilderFactory|SAXParserFactory|XMLInputFactory).*setFeature.*false`),
	"java_serialization":     regexp.MustCompile(`(?i)(ObjectInputStream|readObject)\s*\(`),

	// Cryptographic issues
	"weak_crypto":    regexp.MustCompile(`(?i)(MD5|SHA1|DES|RC4|ECB)\s*\(`),
	"weak_random":    regexp.MustCompile(`(?i)(Math\.random|Random\(\)|rand\(\))`),
	"hardcoded_salt": regexp.MustCompile(`(?i)salt\s*[:=]\s*['"][^'"]+['"]`),

	// Authentication/Authorization
	"auth_bypass":          regexp.MustCompile(`(?i)(auth|login|password)\s*==\s*true`),
	"session_fixation":     regexp.MustCompile(`(?i)session_id\s*=\s*\$_(GET|POST|REQUEST)`),
	"privilege_escalation": regexp.MustCompile(`(?i)(sudo|su|setuid|setgid)\s+.*\$`),

	// Input validation
	"unvalidated_redirect": regexp.MustCompile(`(?i)(redirect|location)\s*=\s*\$_(GET|POST|REQUEST)`),
	"file_upload":          regexp.MustCompile(`(?i)(move_uploaded_file|file_put_contents)\s*\(.*\$_(FILES|POST)`),
	"ldap_injection":       regexp.MustCompile(`(?i)ldap_search\s*\(.*\+.*\$`),

	// Information disclosure
	"error_disclosure":   regexp.MustCompile(`(?i)(error_reporting|display_errors|mysqli_error|pg_last_error)\s*\(\s*.*\)`),
	"debug_info":         regexp.MustCompile(`(?i)(var_dump|print_r|console\.log|System\.out\.println|printStackTrace)\s*\(`),
	"sensitive_data_log": regexp.MustCompile(`(?i)(log|logger)\.(info|debug|error).*\(.*password`),
}

var FunctionPatterns = map[string][]*regexp.Regexp{
	"javascript": {
		// Named function declarations
		regexp.MustCompile(`(?m)^\s*function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\(`),

		// Function expressions assigned to variables (const/let/var)
		regexp.MustCompile(`(?m)^\s*(const|let|var)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*function\s*\(`),

		// Arrow functions assigned to variables
		regexp.MustCompile(`(?m)^\s*(const|let|var)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*\([^)]*\)\s*=>`),

		// Class methods (inside class, simplified)
		// regexp.MustCompile(`(?m)^\s*([a-zA-Z_$][a-zA-Z0-9_$]*)\s*\([^)]*\)\s*\{`),
	},

	"python": {
		// Synchronous function definitions
		regexp.MustCompile(`(?m)^\s*def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),

		// Asynchronous function definitions
		regexp.MustCompile(`(?m)^\s*async\s+def\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),

		// Note: Lambdas are anonymous, usually not counted as named functions
	},

	"go": {
		// Function declarations
		regexp.MustCompile(`(?m)^\s*func\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),

		// Method declarations with receiver
		regexp.MustCompile(`(?m)^\s*func\s+\([^)]*\)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	},

	"java": {
		// Public, private, protected methods (with optional static)
		regexp.MustCompile(`(?m)^\s*(public|private|protected)?\s*(static\s+)?[a-zA-Z0-9_<>\[\]]+\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),

		// Constructors (class name with parentheses and brace)
		regexp.MustCompile(`(?m)^\s*(public|private|protected)?\s*([A-Z][a-zA-Z0-9_]*)\s*\([^)]*\)\s*\{`),

		// Abstract methods
		regexp.MustCompile(`(?m)^\s*abstract\s+[a-zA-Z0-9_<>\[\]]+\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),

		// Interface methods (default or not)
		regexp.MustCompile(`(?m)^\s*(default\s+)?[a-zA-Z0-9_<>\[\]]+\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\([^)]*\)\s*;`),
	},

	//"csharp": {
	//	// Similar to Java method declarations
	//	regexp.MustCompile(`(?m)^\s*(public|private|protected|internal)?\s*(static\s+)?[a-zA-Z0-9_<>\[\]]+\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	//},
	//
	//"php": {
	//	// Function declarations
	//	regexp.MustCompile(`(?m)^\s*function\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	//
	//	// Methods inside classes
	//	regexp.MustCompile(`(?m)^\s*(public|private|protected)?\s*function\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	//},
	//
	//"ruby": {
	//	// Method definitions
	//	regexp.MustCompile(`(?m)^\s*def\s+([a-zA-Z_][a-zA-Z0-9_!?=]*)`),
	//},
	//
	//"swift": {
	//	// Function declarations
	//	regexp.MustCompile(`(?m)^\s*func\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	//},
	//
	//"kotlin": {
	//	// Function declarations
	//	regexp.MustCompile(`(?m)^\s*fun\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
	//},

	// Add more languages and patterns as needed
}

var StylePatterns = map[string]*regexp.Regexp{
	"long_line":            regexp.MustCompile(`.{120,}`),
	"trailing_whitespace":  regexp.MustCompile(`\s+$`),
	"mixed_indentation":    regexp.MustCompile(`^(\t+ +| +\t+)`),
	"snake_case_violation": regexp.MustCompile(`\b[a-z]+([A-Z][a-z]*)+\b`),
	"camel_case_violation": regexp.MustCompile(`\b[a-z]+(_[a-z]+)+\b`),
}

var ComplexityPatterns = map[string]*regexp.Regexp{
	"if_statement": regexp.MustCompile(`(?i)\bif\s*\(`),
	"for_loop":     regexp.MustCompile(`(?i)\bfor\s*\(`),
	"while_loop":   regexp.MustCompile(`(?i)\bwhile\s*\(`),
	"switch_case":  regexp.MustCompile(`(?i)\b(switch|case)\b`),
	"try_catch":    regexp.MustCompile(`(?i)\b(try|catch|except)\b`),
	"ternary":      regexp.MustCompile(`\?.*:`),
	"logical_and":  regexp.MustCompile(`&&`),
	"logical_or":   regexp.MustCompile(`\|\|`),
}

func GetLanguageKeywords() map[models.Language][]string {
	return map[models.Language][]string{
		models.Go: {
			"break", "case", "chan", "const", "continue", "default", "defer", "else",
			"fallthrough", "for", "func", "go", "goto", "if", "import", "interface",
			"map", "package", "range", "return", "select", "struct", "switch", "type", "var",
		},
		models.JavaScript: {
			"async", "await", "break", "case", "catch", "class", "const", "continue",
			"debugger", "default", "delete", "do", "else", "export", "extends", "finally",
			"for", "function", "if", "import", "in", "instanceof", "let", "new", "return",
			"super", "switch", "this", "throw", "try", "typeof", "var", "void", "while", "with", "yield",
		},
		models.Python: {
			"False", "None", "True", "and", "as", "assert", "async", "await", "break",
			"class", "continue", "def", "del", "elif", "else", "except", "finally",
			"for", "from", "global", "if", "import", "in", "is", "lambda", "nonlocal",
			"not", "or", "pass", "raise", "return", "try", "while", "with", "yield",
		},
		models.Java: {
			"abstract", "boolean", "break", "byte", "case", "catch", "char", "class",
			"const", "continue", "default", "do", "double", "else", "extends", "final",
			"finally", "float", "for", "goto", "if", "implements", "import", "instanceof",
			"int", "interface", "long", "native", "new", "package", "private", "protected",
			"public", "return", "short", "static", "strictfp", "super", "switch",
			"synchronized", "this", "throw", "throws", "transient", "try", "void",
			"volatile", "while",
		},
		models.SQL: {
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
			"VARCHAR", "CHAR", "TEXT", "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT",
			"DECIMAL", "NUMERIC", "FLOAT", "DOUBLE", "REAL", "BIT", "BOOLEAN", "BOOL",
			"DATE", "TIME", "DATETIME", "TIMESTAMP", "YEAR", "BLOB", "CLOB", "BINARY",
		},
	}
}
