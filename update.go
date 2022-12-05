package sqlw

type UpdateQuery struct {
	table        string
	fields       []string
	whereQueries []WhereQuery
	jointures    []JoinQuery
	args         []any
	returns      []string
}

func Update(table string) *UpdateQuery {
	return &UpdateQuery{
		table: table,
	}
}

func (q *UpdateQuery) Set(column string, value any) *UpdateQuery {
	q.fields = append(q.fields, column)
	q.args = append(q.args, value)
	return q
}

func (q *UpdateQuery) Returns(columns ...string) *UpdateQuery {
	q.returns = columns
	return q
}

func (q *UpdateQuery) Where(column string, oper string, value any) *UpdateQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "and")
	return q
}

func (q *UpdateQuery) OrWhere(column string, oper string, value any) *UpdateQuery {
	addClauseQuery(&q.whereQueries, column, oper, value, "or")
	return q
}

func (q *UpdateQuery) WhereIn(column string, value any) *UpdateQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "and")
	return q
}

func (q *UpdateQuery) WhereNotIn(column string, value any) *UpdateQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "and")
	return q
}

func (q *UpdateQuery) OrWhereIn(column string, value any) *UpdateQuery {
	addClauseQuery(&q.whereQueries, column, "in", value, "or")
	return q
}

func (q *UpdateQuery) OrWhereNotIn(column string, value any) *UpdateQuery {
	addClauseQuery(&q.whereQueries, column, "not in", value, "or")
	return q
}

func (q *UpdateQuery) Build() (string, []any) {
	query, err := buildUpdateQuery(q.table, q.fields, q.args, q.returns)
	if err != nil {
		panic(err)
	}
	// jointure
	if len(q.jointures) > 0 {
		sb := getJointureQuery(&q.jointures)
		query += sb.String()
	}
	// where
	if len(q.whereQueries) > 0 {
		query += getClauseQuery(wHERE, q.whereQueries, &q.args)
	}
	return query, q.args
}
