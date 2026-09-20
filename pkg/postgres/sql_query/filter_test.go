package sql_query

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhere(t *testing.T) {
	type args struct {
		condition *WhereCondition
	}

	tests := []struct {
		name     string
		args     args
		wantSQL  string
		wantArgs []any
	}{
		{
			name:     "Empty where",
			args:     args{condition: NewWhere()},
			wantSQL:  "",
			wantArgs: []any(nil),
		},
		{
			name:     "Single AND",
			args:     args{condition: NewWhere().And("a = ?", 1)},
			wantSQL:  " WHERE a = $1",
			wantArgs: []any{1},
		},
		{
			name:     "Multi AND",
			args:     args{condition: NewWhere().And("a = ?", 1).And("b > ?", 2)},
			wantSQL:  " WHERE a = $1 AND b > $2",
			wantArgs: []any{1, 2},
		},
		{
			name:     "In with slice",
			args:     args{condition: NewWhere().And("id IN (?)", []int{10, 20, 30})},
			wantSQL:  " WHERE id IN ($1, $2, $3)",
			wantArgs: []any{10, 20, 30},
		},
		{
			name: "group AND only",
			args: args{condition: NewWhere().AndGroup(func(g *WhereCondition) {
				g.And("a = ?", 1).And("b > ?", 2)
			})},
			wantSQL:  " WHERE (a = $1 AND b > $2)",
			wantArgs: []any{1, 2},
		},
		{
			name: "AND with inner group OR",
			args: args{condition: NewWhere().And("a = ?", 1).AndGroup(func(g *WhereCondition) {
				g.And("b = ?", 2).Or("c = ?", 3)
			})},
			wantSQL:  " WHERE a = $1 AND (b = $2 OR c = $3)",
			wantArgs: []any{1, 2, 3},
		},
		{
			name:     "Multiple placeholder in one condition",
			args:     args{condition: NewWhere().And("(a = ? OR b = ?)", 1, 2).And("c BETWEEN ? AND ?", 3, 4)},
			wantSQL:  " WHERE (a = $1 OR b = $2) AND c BETWEEN $3 AND $4",
			wantArgs: []any{1, 2, 3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotArgs := tt.args.condition.SQL()
			require.Equal(t, tt.wantSQL, gotSQL)
			require.Equal(t, tt.wantArgs, gotArgs)
		})
	}
}

func TestOrderAndPagination(t *testing.T) {
	const selectQuery = `
select
	ex.id as "id",
	ex."pharmId" as "pharmId"
from
	exchange ex
`

	userId1 := 123
	userId2 := 124
	take := 10
	skip := 20

	where, args := NewWhere().
		And(`(ew."userId" = ? OR ew."userId" = ?)`, userId1, userId2).
		And(`ew."confirmedAt" is null`).
		And(`ew."deletedAt" is null`).
		SQL()

	order := NewOrder().
		OrderBy(`ex."startTime"`, Asc, false).
		SQL()

	pagination := NewPagination().Take(take).Skip(skip).SQL()
	query := selectQuery + where + order + pagination

	result := `
select
	ex.id as "id",
	ex."pharmId" as "pharmId"
from
	exchange ex
 WHERE (ew."userId" = $1 OR ew."userId" = $2) AND ew."confirmedAt" is null AND ew."deletedAt" is null ORDER BY ex."startTime" ASC NULLS LAST LIMIT 10 OFFSET 20`

	assert.Equal(t, query, result)

	expectArg := []any{userId1, userId2}
	assert.Equal(t, expectArg, args)
}
