package sqlw

import (
	"fmt"
	"strings"
)

type ValueFunc struct {
	value       any
	alternative any
	fun         string
	cast        string
}

func Truncate(table string, restartIdentity bool, cascade bool) string {
	query := "truncate table " + table
	if restartIdentity {
		query += " restart identity"
	}
	if cascade {
		query += " cascade"
	}
	return query
}

func Drop(table string, cascade bool) string {
	query := "drop table if exists " + table
	if cascade {
		query += " cascade"
	}
	return query
}

func Coalesce(value any, alternative any, cast ...string) ValueFunc {
	v := ValueFunc{
		value:       value,
		fun:         "coalesce",
		alternative: alternative,
	}
	if len(cast) > 0 {
		v.cast = cast[0]
	}
	return v
}

func Nullif(value any, alternative any, cast ...string) ValueFunc {
	v := ValueFunc{
		value:       value,
		fun:         "nullif",
		alternative: alternative,
	}
	if len(cast) > 0 {
		v.cast = cast[0]
	}
	return v
}

func recFuncValue(fv ValueFunc, argIndex int) (string, any) {
	var castype string
	if fv.cast != "" {
		castype = "::" + fv.cast
	}
	val, ok := fv.value.(ValueFunc)
	s := fmt.Sprintf("%s($%v%s,%v)", fv.fun, argIndex, castype, fv.alternative)
	value := fv.value
	if ok {
		t, v := recFuncValue(val, argIndex)
		value = v
		s = strings.ReplaceAll(s, fmt.Sprintf("$%v", argIndex), t)
	}
	return s, value
}
