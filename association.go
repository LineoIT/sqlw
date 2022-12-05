package sqlw

import (
	"fmt"
	"strings"
)

type HasQuery struct {
	table          string
	columns        []string
	associated     []associate
	whereQueries   []WhereQuery
	havingQueries  []WhereQuery
	groupByColumns []string
	orderColumns   []string
	ascending      string
	take           uint
	skip           uint
}

type associate struct {
	table, leftColumn, rightColumn, alias string
	hasMany                               bool
}

func (s *SelectQuery) HasMany(table, leftColumn, rightColumn, alias string) *HasQuery {
	return &HasQuery{
		table:   s.table,
		columns: s.columns,
		associated: []associate{{
			table:       table,
			leftColumn:  leftColumn,
			rightColumn: rightColumn,
			alias:       alias,
			hasMany:     true,
		}},
	}
}

func (s *SelectQuery) Has(table, leftColumn, rightColumn, alias string) *HasQuery {
	return &HasQuery{
		table:   s.table,
		columns: s.columns,
		associated: []associate{{
			table:       table,
			leftColumn:  leftColumn,
			rightColumn: rightColumn,
			alias:       alias,
			hasMany:     false,
		}},
	}
}

func (q *HasQuery) HasMany(table, leftColumn, rightColumn, alias string) *HasQuery {
	q.associated = append(q.associated, associate{
		table:       table,
		leftColumn:  leftColumn,
		rightColumn: rightColumn,
		alias:       alias,
		hasMany:     true,
	})
	return q
}

func (q *HasQuery) Has(table, leftColumn, rightColumn, alias string) *HasQuery {
	q.associated = append(q.associated, associate{
		table:       table,
		leftColumn:  leftColumn,
		rightColumn: rightColumn,
		alias:       alias,
		hasMany:     false,
	})
	return q
}

func (q *HasQuery) Where(column string, oper string, value any) *HasQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "and")
	return q
}

func (q *HasQuery) OrWhere(column string, oper string, value any) *HasQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "or")
	return q
}

func (q *HasQuery) WhereIn(column string, value any) *HasQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "and")
	return q
}

func (q *HasQuery) WhereNotIn(column string, value any) *HasQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "and")
	return q
}

func (q *HasQuery) OrWhereIn(column string, value any) *HasQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "or")
	return q
}

func (q *HasQuery) OrWhereNotIn(column string, value any) *HasQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "or")
	return q
}

func (q *HasQuery) Having(column string, oper string, value any) *HasQuery {
	addClauseQuery(&q.havingQueries, column, oper, value, "and")
	return q
}

func (q *HasQuery) OrHaving(column string, oper string, value any) *HasQuery {
	addClauseQuery(&q.havingQueries, column, oper, value, "or")
	return q
}

func (q *HasQuery) HavingIn(column string, value any) *HasQuery {
	addClauseQuery(&q.havingQueries, column, "in", value, "and")
	return q
}

func (q *HasQuery) HavingNotIn(column string, value any) *HasQuery {
	addClauseQuery(&q.havingQueries, column, "not in", value, "and")
	return q
}

func (q *HasQuery) OrHavingIn(column string, value any) *HasQuery {
	addClauseQuery(&q.havingQueries, column, "in", value, "or")
	return q
}

func (q *HasQuery) OrHavingNotIn(column string, value any) *HasQuery {
	addClauseQuery(&q.havingQueries, column, "not in", value, "or")
	return q
}

func (q *HasQuery) GroupBy(columns ...string) *HasQuery {
	q.groupByColumns = columns
	return q
}

func (q *HasQuery) OrderBy(columns ...string) *HasQuery {
	q.orderColumns = columns
	return q
}

func (q *HasQuery) Desc() *HasQuery {
	q.ascending = "desc"
	return q
}

func (q *HasQuery) Asc() *HasQuery {
	q.ascending = "asc"
	return q
}

func (q *HasQuery) Limit(limit uint) *HasQuery {
	q.take = limit
	return q
}

func (q *HasQuery) Offset(offset uint) *HasQuery {
	q.skip = offset
	return q
}

func (q *HasQuery) Build() (string, []any) {
	var args []any
	var query strings.Builder
	query.WriteString("select replace(concat(")
	if len(q.columns) > 0 {
		query.WriteString("json_build_object(")
		keyVal := []string{}
		for _, v := range q.columns {
			keyVal = append(keyVal, fmt.Sprintf("'%s',%s", strings.Split(v, ".")[1], v))
		}
		query.WriteString(strings.Join(keyVal, ","))
		query.WriteString(")")
	} else {
		query.WriteString(fmt.Sprintf("row_to_json(%s.*)", q.table))
	}
	query.WriteString(", ''), '}', ',')")
	if len(q.associated) > 0 {
		assoc := []string{}
		query.WriteString(" ||")
		for _, v := range q.associated {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf(" concat('\"%s\":'", v.alias))

			if v.hasMany {
				sb.WriteString(",json_agg(")
			} else {
				sb.WriteString(",row_to_json(")
			}
			sb.WriteString(v.table + ".*")
			sb.WriteString("))")
			assoc = append(assoc, sb.String())
		}
		query.WriteString(strings.Join(assoc, " || "))
	}

	query.WriteString(" || '}' from " + q.table)
	if len(q.associated) > 0 {
		for _, v := range q.associated {
			query.WriteString(fmt.Sprintf(" join %s on %s=%s", v.table, v.leftColumn, v.rightColumn))
		}
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
