package sqlw

type InsertQuery struct {
	table   string
	fields  []string
	args    []any
	returns []string
}

func Insert(table string) *InsertQuery {
	return &InsertQuery{
		table: table,
	}
}

func (q *InsertQuery) Value(column string, value any) *InsertQuery {
	q.fields = append(q.fields, column)
	q.args = append(q.args, value)
	return q
}

func (q *InsertQuery) Returns(columns ...string) *InsertQuery {
	q.returns = columns
	return q
}

func (q *InsertQuery) Build() (string, []any) {
	return buildInsertQuery(q.table, q.fields, q.args, q.returns).String(), q.args
}
