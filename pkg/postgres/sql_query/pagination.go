package sql_query

import (
	"fmt"
	"strings"
)

type PaginationCondition struct {
	parts []string
}

func NewPagination() *PaginationCondition {
	return &PaginationCondition{}
}

func (p *PaginationCondition) Take(v int) *PaginationCondition {
	p.parts = append(p.parts, fmt.Sprintf(" LIMIT %d", v))

	return p
}

func (p *PaginationCondition) Skip(v int) *PaginationCondition {
	p.parts = append(p.parts, fmt.Sprintf(" OFFSET %d", v))

	return p
}

func (p *PaginationCondition) SQL() string {
	if len(p.parts) == 0 {
		return ""
	}

	return strings.Join(p.parts, "")
}
