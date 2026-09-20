// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/leabago/share-radio/adder/docs/gen"
	"github.com/leabago/share-radio/adder/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	StationRepo interface {
		CreateStation(ctx context.Context, station *entity.Station) (string, error)
		ListStations(ctx context.Context, request gen.ListStationsRequestObject) ([]entity.Station, error)
		GetStation(ctx context.Context, id string) (entity.Station, error)
	}

	GenreRepo interface {
		ListGenres(ctx context.Context) ([]entity.Genre, error)
	}

	LanguageRepo interface {
		ListLanguages(ctx context.Context) ([]entity.Language, error)
	}
)
