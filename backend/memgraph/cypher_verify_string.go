package memgraph

import (
	"fmt"
	"strings"
)

var cypherKeywords = []string{
	"CREATE",
	"DELETE",
	"DETACH",
	"DETACH DELETE",
	"FOREACH",
	"LOAD CSV",
	"MERGE",
	"MATCH",
	"ON",
	"OPTIONAL MATCH",
	"REMOVE",
	"SET",
	"START",
	"UNION",
	"UNWIND",
	"WITH",
	"RETURN",
	"ORDER BY",
	"SKIP",
	"LIMIT",
	"ASC",
	"DESC",
	"EXISTS",
	"CALL",
	"USING",
	"CONSTRAINT",
	"DROP",
	"INDEX",
	"WHERE",
}

var cypherOperators = []string{
	"+",
	"-",
	"*",
	"/",
	"%",
	"^",
	"=",
	"<",
	">",
	"<=",
	">=",
	"<>",
	"AND",
	"OR",
	"XOR",
	"NOT",
	"IN",
	"STARTS WITH",
	"ENDS WITH",
	"CONTAINS",
	"IS NULL",
	"IS NOT NULL",
	"IS UNIQUE",
	"IS NODE",
	"IS RELATIONSHIP",
	"IS PROPERTY KEY",
	"IS MAP",
	"IS LIST",
	"IS BOOLEAN",
	"IS STRING",
	"IS NUMBER",
	"IS INTEGER",
	"IS FLOAT",
	"IS NODE",
	"IS RELATIONSHIP",
	"IS PATH",
	"IS POINT",
	"IS DATE",
	"IS DURATION",
}

// cypherDelimiters contains the delimiters that need to be escaped in a string to prevent cypher injection keys are the delimiters that need to be escaped and values are the escaped delimiters
var cypherDelimiters = map[string]string{
	"'":       "\\'",
	"\"":      "\\\"",
	"\\u0027": "\\\\u0027",
	"\\u0022": "\\\\u0022",
	"`":       "``",
	"\\u0060": "\\u0060\\u0060",
}

// VerifyString verifies if a string is valid and does not contain cypher injection
func VerifyString(s string) error {
	s = strings.ToUpper(s)
	for _, keyword := range cypherKeywords {
		if strings.Contains(s, keyword) {
			return fmt.Errorf("invalid string contains cypher keyword: %s", keyword)
		}
	}

	for _, operator := range cypherOperators {
		if strings.Contains(s, operator) {
			return fmt.Errorf("invalid string contains cypher operator: %s", operator)
		}
	}

	for key := range cypherDelimiters {
		if strings.Contains(s, key) {
			return fmt.Errorf("invalid string contains cypher delimiter: %s", key)
		}
	}

	return nil
}

// EscapeString escapes delimiters in a string to prevent cypher injection
func EscapeString(s string) string {
	result := s
	for k, v := range cypherDelimiters {
		result = strings.ReplaceAll(result, k, v)
	}

	return result
}
