package sqlw

import (
	"fmt"
	"reflect"
	"strings"
)

type clause string

const (
	wHERE  clause = "where"
	hAVING clause = "having"
)

type WhereQuery struct {
	column   string
	value    interface{}
	operator string
	logic    string
}

func addClauseQuery(whereQueries *[]WhereQuery, column, oper string, value interface{}, clause string) {
	*whereQueries = append(*whereQueries, WhereQuery{
		column:   column,
		value:    value,
		operator: oper,
		logic:    clause,
	})
}

func getClauseQuery(cls clause, fv []WhereQuery, args *[]interface{}) string {
	var query strings.Builder
	query.WriteString(fmt.Sprintf(" %s ", cls))
	for k, v := range fv {
		query.WriteString(v.column)
		if reflect.TypeOf(v.value).Kind() == reflect.Slice {
			query.WriteString(" in (")
			values := reflect.ValueOf(v.value)
			for j := 0; j < values.Len(); j++ {
				query.WriteString(fmt.Sprintf("$%d", len(*args)+1))
				*args = append(*args, values.Index(j).Interface())
				if j < values.Len()-1 {
					query.WriteString(",")
				}
			}
			query.WriteString(")")
		} else {
			query.WriteString(fmt.Sprintf(" %s $%d", v.operator, len(*args)+1))
			*args = append(*args, v.value)
		}
		if k < len(fv)-1 {
			query.WriteString(" " + fv[k+1].logic + " ")
		}
	}
	return query.String()
}
