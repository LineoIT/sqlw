package sqlw

import (
	"fmt"
	"strings"
)

type JoinQuery struct {
	table, leftColumn, rightColumn, jointure string
}

// select
func (s *SelectQuery) Join(table string, leftColumn, rightColumn string) *SelectQuery {
	addJointure(&s.jointures, "join", table, leftColumn, rightColumn)
	return s
}

func (s *SelectQuery) LetfJoin(table string, leftColumn, rightColumn string) *SelectQuery {
	addJointure(&s.jointures, "left join", table, leftColumn, rightColumn)
	return s
}

func (s *SelectQuery) RightJoin(table string, leftColumn, rightColumn string) *SelectQuery {
	addJointure(&s.jointures, "right join", table, leftColumn, rightColumn)
	return s
}

func (s *SelectQuery) InnerJoin(table string, leftColumn, rightColumn string) *SelectQuery {
	addJointure(&s.jointures, "inner join", table, leftColumn, rightColumn)
	return s
}

func (s *UpdateQuery) Join(table string, leftColumn, rightColumn string) *UpdateQuery {
	addJointure(&s.jointures, "join", table, leftColumn, rightColumn)
	return s
}

func (s *UpdateQuery) LetfJoin(table string, leftColumn, rightColumn string) *UpdateQuery {
	addJointure(&s.jointures, "left join", table, leftColumn, rightColumn)
	return s
}

func (s *UpdateQuery) RightJoin(table string, leftColumn, rightColumn string) *UpdateQuery {
	addJointure(&s.jointures, "right join", table, leftColumn, rightColumn)
	return s
}

func (s *UpdateQuery) InnerJoin(table string, leftColumn, rightColumn string) *UpdateQuery {
	addJointure(&s.jointures, "inner join", table, leftColumn, rightColumn)
	return s
}

func getJointureQuery(jointures *[]JoinQuery) strings.Builder {
	var sb strings.Builder
	for _, v := range *jointures {
		if v.leftColumn != "" && v.rightColumn != "" {
			sb.WriteString(fmt.Sprintf(" %s %s on %s=%s ", v.jointure, v.table, v.leftColumn, v.rightColumn))
		}
	}
	return sb
}

func addJointure(jointures *[]JoinQuery, jointure, table, leftColumn, rightColumn string) {
	*jointures = append(*jointures, JoinQuery{
		table:       table,
		rightColumn: rightColumn,
		leftColumn:  leftColumn,
		jointure:    jointure,
	})
}
