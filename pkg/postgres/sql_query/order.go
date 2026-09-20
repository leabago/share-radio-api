package sql_query

import (
	"fmt"
	"strings"
)

type SortOrder string

const (
	Asc  SortOrder = "ASC"
	Desc SortOrder = "DESC"
)

func (s SortOrder) IsValid() bool {
	switch s {
	case Asc, Desc:
		return true
	}
	return false
}

func ParseSortOrder(s string) (SortOrder, error) {
	status := SortOrder(strings.ToUpper(s))
	if status.IsValid() {
		return status, nil
	}
	return "", fmt.Errorf("invalid SortOrder: %s", s)
}

type OrderCondition struct {
	parts []string
}

func NewOrder() *OrderCondition {
	return &OrderCondition{}
}

func (o *OrderCondition) OrderBy(field string, direction SortOrder, nulls bool) *OrderCondition {
	var nullStr string
	if nulls {
		nullStr = " NULLS FIRST"
	} else {
		nullStr = " NULLS LAST"
	}

	o.parts = append(o.parts, fmt.Sprintf("%s %s%s", field, direction, nullStr))

	return o
}

func (o *OrderCondition) GetLenParts() int {
	return len(o.parts)
}

func (o *OrderCondition) SQL() string {
	if len(o.parts) == 0 {
		return ""
	}

	return " ORDER BY " + strings.Join(o.parts, ",")
}
