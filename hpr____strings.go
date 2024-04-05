package filter

import "strings"

func to_escaped_string(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

func escape_single_quote(v string) string {
	return strings.ReplaceAll(v, "'", "''")
}

func wrap_with_single_quote(v string) string {
	return "'" + v + "'"
}
