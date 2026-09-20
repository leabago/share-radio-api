package station

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/leabago/share-radio/adder/docs/gen"
	"github.com/leabago/share-radio/adder/internal/entity"
	"github.com/leabago/share-radio/adder/pkg/helper/functions"
	"github.com/leabago/share-radio/adder/pkg/postgres"
	"github.com/leabago/share-radio/adder/pkg/postgres/sql_query"
)

const createStation = `
INSERT INTO stations.radio_stations
("name", url, genre, "language", added_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6)
returning id;
`

const selectStation = `
SELECT radio_stations."id", radio_stations."name", radio_stations."url", 
genres."name" as genre_name, languages."name" language_name, radio_stations."icon", 
radio_stations."is_active", radio_stations."is_new", 
radio_stations."added_at", radio_stations."updated_at"
FROM stations.radio_stations
LEFT JOIN stations.genres 
ON radio_stations.genre = stations.genres.id
LEFT JOIN stations.languages
ON radio_stations.language = stations.languages.id
`

var errStationNotFound = errors.New("station not found")

type Repo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *Repo {
	return &Repo{
		Postgres: pg,
	}
}

func (r *Repo) CreateStation(ctx context.Context, task *entity.Station) (string, error) {

	row, err := r.Postgres.Pool.Query(
		ctx,
		createStation,
		task.Name,
		task.Url,
		task.GenreId,
		task.LanguageId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return "", err
	}

	res, err := pgx.CollectOneRow(row, pgx.RowToStructByName[entity.StationId])
	if err != nil {
		return "", err
	}

	return res.Id, nil
}

func (r *Repo) ListStations(ctx context.Context, request gen.ListStationsRequestObject) ([]entity.Station, error) {

	// Where
	where := sql_query.NewWhere()

	where.And(`radio_stations."is_active" = true`)

	if request.Params.GenreId != nil {

		where.And(`radio_stations."genre" = ?`, request.Params.GenreId.String())
	}

	if request.Params.LanguageId != nil {

		where.And(`radio_stations."language" = ?`, request.Params.LanguageId.String())
	}

	if request.Params.Search != nil {
		where.And(`radio_stations."name" ILIKE ?`, "%"+*request.Params.Search+"%")
	}

	whereSql, args := where.SQL()

	// Pagination

	paginationSql := sql_query.NewPagination().
		Take(functions.GetDefaultValue(request.Params.Limit)).
		Skip(functions.GetDefaultValue(request.Params.Offset)).SQL()

	// Order
	orderSql := ""

	if request.Params.SortBy != nil {
		order := sql_query.NewOrder()

		sortOrder := sql_query.Asc

		paramsSortOrder := request.Params.SortOrder

		if paramsSortOrder != nil {
			parseOrder, err := sql_query.ParseSortOrder(string(*paramsSortOrder))
			if err != nil {
				return nil, err
			}
			sortOrder = parseOrder
		}

		switch *request.Params.SortBy {
		case gen.AddedAt:
			order.OrderBy(string(gen.AddedAt), sortOrder, false)
		case gen.Name:
			order.OrderBy(string(gen.Name), sortOrder, false)
		}

		orderSql = order.SQL()
	}

	sqlRequest := selectStation + whereSql + orderSql + paginationSql

	fmt.Printf("SQL: %s\n", whereSql)
	fmt.Printf("Args: %+v (count: %d)\n", args, len(args))
	fmt.Printf("sqlRequest: %s\n", sqlRequest)

	rows, err := r.Postgres.Pool.Query(ctx, sqlRequest, args...)
	if err != nil {
		return nil, err
	}

	resp, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.Station])
	if err != nil {
		return resp, err
	}

	return resp, nil

}

func (r *Repo) GetStation(ctx context.Context, id string) (entity.Station, error) {

	whereSql, args := sql_query.NewWhere().
		And(`radio_stations."id" = ?`, id).
		And(`radio_stations."is_active" = true`).SQL()

	sqlRequest := selectStation + whereSql

	rows, err := r.Postgres.Pool.Query(
		ctx,
		sqlRequest,
		args...,
	)
	if err != nil {
		return entity.Station{}, err
	}

	resp, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Station])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Station{}, errStationNotFound
		}
		return resp, err
	}

	return resp, nil

}
