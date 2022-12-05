package sqlw

import (
	"fmt"
	"strings"
)

type SelectQuery struct {
	table          string
	columns        []string
	whereQueries   []WhereQuery
	havingQueries  []WhereQuery
	groupByColumns []string
	orderColumns   []string
	ascending      string
	take           uint
	skip           uint
	jointures      []JoinQuery
}

func Select(table string, columns ...string) *SelectQuery {
	return &SelectQuery{
		table:   table,
		columns: columns,
	}
}

func (q *SelectQuery) GroupBy(columns ...string) *SelectQuery {
	q.groupByColumns = columns
	return q
}

func (q *SelectQuery) OrderBy(columns ...string) *SelectQuery {
	q.orderColumns = columns
	return q
}

func (q *SelectQuery) Desc() *SelectQuery {
	q.ascending = "desc"
	return q
}

func (q *SelectQuery) Asc() *SelectQuery {
	q.ascending = "asc"
	return q
}

func (q *SelectQuery) Limit(limit uint) *SelectQuery {
	q.take = limit
	return q
}

func (q *SelectQuery) Offset(offset uint) *SelectQuery {
	q.skip = offset
	return q
}

func (q *SelectQuery) Where(column string, oper string, value any) *SelectQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "and")
	return q
}

func (q *SelectQuery) OrWhere(column string, oper string, value any) *SelectQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "or")
	return q
}

func (q *SelectQuery) WhereIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "and")
	return q
}

func (q *SelectQuery) WhereNotIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "and")
	return q
}

func (q *SelectQuery) OrWhereIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "or")
	return q
}

func (q *SelectQuery) OrWhereNotIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "or")
	return q
}

func (q *SelectQuery) Having(column string, oper string, value any) *SelectQuery {
	addClauseQuery(&q.havingQueries, column, oper, value, "and")
	return q
}

func (q *SelectQuery) OrHaving(column string, oper string, value any) *SelectQuery {
	addClauseQuery(&q.havingQueries, column, oper, value, "or")
	return q
}

func (q *SelectQuery) HavingIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.havingQueries, column, "in", value, "and")
	return q
}

func (q *SelectQuery) HavingNotIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.havingQueries, column, "not in", value, "and")
	return q
}

func (q *SelectQuery) OrHavingIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.havingQueries, column, "in", value, "or")
	return q
}

func (q *SelectQuery) OrHavingNotIn(column string, value any) *SelectQuery {
	addClauseQuery(&q.havingQueries, column, "not in", value, "or")
	return q
}

func (q *SelectQuery) Build() (string, []any) {
	var query strings.Builder
	var args []any
	// select from
	query.WriteString("select ")
	if len(q.columns) > 0 {
		query.WriteString(strings.Join(q.columns, ","))
	} else {
		query.WriteString("*")
	}
	query.WriteString(" from " + q.table)
	// jointure
	if len(q.jointures) > 0 {
		sb := getJointureQuery(&q.jointures)
		query.WriteString(sb.String())
	}
	// where
	if len(q.whereQueries) > 0 {
		query.WriteString(getClauseQuery(wHERE, q.whereQueries, &args))
	}
	// group by
	if len(q.groupByColumns) > 0 {
		query.WriteString(" group by " + strings.Join(q.groupByColumns, ","))
	}

	// having
	if len(q.havingQueries) > 0 {
		sb := getClauseQuery(hAVING, q.havingQueries, &args)
		query.WriteString(sb)
	}

	// order by
	if len(q.orderColumns) > 0 {
		query.WriteString(" order by " + strings.Join(q.orderColumns, ","))
	}

	// ascending
	if q.ascending != "" {
		query.WriteString(" " + q.ascending)
	}

	// limit
	if q.take > 0 {
		query.WriteString(fmt.Sprintf(" limit $%d", len(args)+1))
		args = append(args, q.take)
	}
	// offset
	if q.skip != 0 {
		query.WriteString(fmt.Sprintf(" offset $%d", len(args)+1))
		args = append(args, q.skip)
	}
	return query.String(), args
}
