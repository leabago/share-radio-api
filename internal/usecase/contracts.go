// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/leabago/share-radio/adder/docs/gen"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	Station interface {
		CreateStation(ctx context.Context, request gen.CreateStationRequestObject) (gen.StationId, error)
		ListStations(ctx context.Context, request gen.ListStationsRequestObject) (gen.StationListResponse, error)
		GetStation(ctx context.Context, request gen.GetStationRequestObject) (gen.Station, error)
	}

	Genre interface {
		ListGenres(ctx context.Context, request gen.ListGenresRequestObject) (gen.GenreList, error)
	}

	Language interface {
		ListLanguages(ctx context.Context, request gen.ListLanguagesRequestObject) (gen.LanguageList, error)
	}
)
