package sql_query

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	type args struct {
		condition *OrderCondition
	}

	tests := []struct {
		name    string
		args    args
		wantSQL string
	}{
		{
			name:    "Single ASC NULLS FIRST",
			args:    args{condition: NewOrder().OrderBy("name", Asc, true)},
			wantSQL: " ORDER BY name ASC NULLS FIRST",
		},
		{
			name:    "Single ASC NULLS LAST",
			args:    args{condition: NewOrder().OrderBy("name", Asc, false)},
			wantSQL: " ORDER BY name ASC NULLS LAST",
		},
		{
			name:    "Single DESC NULLS FIRST",
			args:    args{condition: NewOrder().OrderBy("name", Desc, true)},
			wantSQL: " ORDER BY name DESC NULLS FIRST",
		},
		{
			name: "Multiple fields",
			args: args{
				condition: NewOrder().OrderBy("name", Asc, true).
					OrderBy("age", Desc, false).
					OrderBy("score", Asc, false),
			},
			wantSQL: " ORDER BY name ASC NULLS FIRST,age DESC NULLS LAST,score ASC NULLS LAST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL := tt.args.condition.SQL()
			require.Equal(t, tt.wantSQL, gotSQL)
		})
	}
}

func TestSortOrderOk(t *testing.T) {
	order := "ASC"
	so, err := ParseSortOrder(order)
	require.NoError(t, err)
	assert.True(t, so.IsValid())
}

func TestSortOrderNotCorrect(t *testing.T) {
	order := "ASC2"
	_, err := ParseSortOrder(order)
	require.Error(t, err)

	assert.EqualError(t, err, fmt.Sprintf("invalid SortOrder: %s", order))

}
