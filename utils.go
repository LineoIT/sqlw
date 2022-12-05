package sqlw

import (
	"errors"
	"fmt"
	"strings"
)

func CleanSQL(query string) string {
	return strings.ReplaceAll(strings.Trim(strings.ReplaceAll(
		strings.ReplaceAll(query, "\t", ""), "\n", " "), " "), "  ", " ")
}

func Debug(query string, args ...any) string {
	s := CleanSQL(query)
	for k, v := range args {
		s = strings.Replace(s, fmt.Sprintf("$%d", k+1), fmt.Sprint(v), -1)
	}
	return s
}

func buildInsertQuery(table string, columns []string, args []any, returns []string) *strings.Builder {
	var query strings.Builder
	query.WriteString("insert into " + table)
	if len(columns) > 0 {
		query.WriteString(fmt.Sprintf("(%s)", strings.Join(columns, ",")))
	}
	query.WriteString(" values(")
	for i := range args {
		// added
		funcValue, ok := args[i].(ValueFunc)
		if ok {
			stmt, value := recFuncValue(funcValue, i+1)
			query.WriteString(stmt)
			args[i] = value
		} else {
			query.WriteString(fmt.Sprintf("$%d", i+1))
		}
		// query.WriteString(fmt.Sprintf("$%d", i+1))

		if i < len(args)-1 {
			query.WriteString(",")
		}
	}
	query.WriteString(")")
	// added
	if len(returns) > 0 {
		query.WriteString(fmt.Sprintf(" returning %s", strings.Join(returns, ",")))
	}
	return &query
}

func buildUpdateQuery(table string, columns []string, args []any, returns []string) (string, error) {
	var query strings.Builder
	query.WriteString("update " + table + " set ")
	fieldCount := len(columns)
	if fieldCount == 0 || len(args) == 0 {
		return "", errors.New("update:fields and values need")
	}
	if fieldCount != len(args) {
		return "", errors.New("update:fields and values len don't match")
	}
	for i := 0; i < fieldCount; i++ {
		// added
		funcValue, ok := args[i].(ValueFunc)
		if ok {
			stmt, value := recFuncValue(funcValue, i+1)
			query.WriteString(columns[i] + "=" + stmt)
			args[i] = value
		} else {
			query.WriteString(fmt.Sprintf("%s=$%d", columns[i], i+1))
		}
		// query.WriteString(fmt.Sprintf("%s=coalesce($%d,%s)", columns[i], i+1, columns[i]))
		if i < fieldCount-1 {
			query.WriteString(",")
		}
	}
	// added
	if len(returns) > 0 {
		query.WriteString(fmt.Sprintf(" returning %s", strings.Join(returns, ",")))
	}
	return query.String(), nil
}
