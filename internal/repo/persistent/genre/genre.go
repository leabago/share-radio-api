package genre

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/leabago/share-radio/adder/internal/entity"
	"github.com/leabago/share-radio/adder/pkg/postgres"
	"github.com/leabago/share-radio/adder/pkg/postgres/sql_query"
)

const selectGenres = `
SELECT "id", "name"
FROM stations.genres
`

type Repo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		Postgres: pg,
	}
}

func (r *Repo) ListGenres(ctx context.Context) ([]entity.Genre, error) {

	order := sql_query.NewOrder().OrderBy("display_order", sql_query.Asc, false).SQL()

	selectGenresReq := selectGenres + order

	rows, err := r.Postgres.Pool.Query(
		ctx,
		selectGenresReq,
	)
	if err != nil {
		return nil, err
	}

	genres, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Genre])
	if err != nil {
		return nil, err
	}

	return genres, nil
}

//func (r *Repo) GetGenreByName(ctx context.Context, name string) (*entity.Genre, error) {
//
//	where, args := sql_query.NewWhere().And("name = ?", name).SQL()
//
//	request := selectGenres + where
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
//	genre, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Genre])
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return nil, entity.ErrNotFound
//		}
//		return nil, err
//	}
//
//	return &genre, nil
//}
