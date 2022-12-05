package sqlw

import (
	"strings"
)

type DeleteQuery struct {
	table        string
	whereQueries []WhereQuery
}

func Delete(table string) *DeleteQuery {
	return &DeleteQuery{
		table: table,
	}
}

func (q *DeleteQuery) Where(column string, oper string, value any) *DeleteQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "and")
	return q
}

func (q *DeleteQuery) OrWhere(column string, oper string, value any) *DeleteQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "or")
	return q
}

func (q *DeleteQuery) WhereIn(column string, value any) *DeleteQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "and")
	return q
}

func (q *DeleteQuery) WhereNotIn(column string, value any) *DeleteQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "and")
	return q
}

func (q *DeleteQuery) OrWhereIn(column string, value any) *DeleteQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "or")
	return q
}

func (q *DeleteQuery) OrWhereNotIn(column string, value any) *DeleteQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "or")
	return q
}

func (q *DeleteQuery) Build(args *[]any) string {
	var query strings.Builder
	query.WriteString("delete from " + q.table)
	// where
	if len(q.whereQueries) > 0 {
		sb := getClauseQuery(wHERE, q.whereQueries, args)
		query.WriteString(sb)
	}
	return query.String()
}
