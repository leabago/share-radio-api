package language

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/leabago/share-radio/adder/internal/entity"
	"github.com/leabago/share-radio/adder/pkg/postgres"
	"github.com/leabago/share-radio/adder/pkg/postgres/sql_query"
)

const selectLanguages = `
SELECT "id", "name"
FROM stations.languages
`

type Repo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		Postgres: pg,
	}
}

func (r *Repo) ListLanguages(ctx context.Context) ([]entity.Language, error) {

	order := sql_query.NewOrder().OrderBy("display_order", sql_query.Asc, false).SQL()

	request := selectLanguages + order

	rows, err := r.Postgres.Pool.Query(
		ctx,
		request,
	)
	if err != nil {
		return nil, err
	}

	languages, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Language])
	if err != nil {
		return nil, err
	}

	return languages, nil
}

//
//func (r *Repo) GetLanguageByName(ctx context.Context, name string) (*entity.Language, error) {
//
//	where, args := sql_query.NewWhere().And("name = ?", name).SQL()
//
//	request := selectLanguages + where
//
//	rows, err := r.Postgres.Pool.Query(
//		ctx,
//		request,
//		args,
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	language, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Language])
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, entity.ErrNotFound
//		}
//		return nil, err
//	}
//
//	return &language, nil
//}
