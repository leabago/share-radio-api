package api

import (
	"context"

	"github.com/leabago/share-radio/adder/docs/gen"
	"github.com/leabago/share-radio/adder/internal/usecase"
)

func errUnhandled(err error) gen.ErrorResponse {
	return gen.ErrorResponse{
		Error: new(err.Error()),
	}
}

var _ gen.StrictServerInterface = (*Controller)(nil)

type Controller struct {
	stationService  usecase.Station
	genreService    usecase.Genre
	languageService usecase.Language
}

func NewController(stationsService usecase.Station, genresService usecase.Genre, languageService usecase.Language) Controller {
	return Controller{
		stationService:  stationsService,
		genreService:    genresService,
		languageService: languageService,
	}
}

func (c Controller) ListStations(ctx context.Context, request gen.ListStationsRequestObject) (gen.ListStationsResponseObject, error) {
	resp, err := c.stationService.ListStations(ctx, request)
	if err != nil {
		return gen.ListStations500JSONResponse(errUnhandled(err)), nil
	}

	return gen.ListStations200JSONResponse(resp), nil
}

func (c Controller) CreateStation(ctx context.Context, request gen.CreateStationRequestObject) (gen.CreateStationResponseObject, error) {

	resp, err := c.stationService.CreateStation(ctx, request)
	if err != nil {
		return gen.CreateStation500JSONResponse(errUnhandled(err)), nil
	}

	return gen.CreateStation201JSONResponse(resp), nil

}

func (c Controller) DeleteStation(ctx context.Context, request gen.DeleteStationRequestObject) (gen.DeleteStationResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (c Controller) GetStation(ctx context.Context, request gen.GetStationRequestObject) (gen.GetStationResponseObject, error) {
	resp, err := c.stationService.GetStation(ctx, request)
	if err != nil {
		return gen.GetStation500JSONResponse(errUnhandled(err)), nil
	}

	return gen.GetStation200JSONResponse(resp), nil
}

func (c Controller) UpdateStation(ctx context.Context, request gen.UpdateStationRequestObject) (gen.UpdateStationResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (c Controller) RecordPlay(ctx context.Context, request gen.RecordPlayRequestObject) (gen.RecordPlayResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (c Controller) RateStation(ctx context.Context, request gen.RateStationRequestObject) (gen.RateStationResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (c Controller) UpdateStationStatus(ctx context.Context, request gen.UpdateStationStatusRequestObject) (gen.UpdateStationStatusResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (c Controller) ListGenres(ctx context.Context, request gen.ListGenresRequestObject) (gen.ListGenresResponseObject, error) {

	resp, err := c.genreService.ListGenres(ctx, request)
	if err != nil {
		return gen.ListGenres500JSONResponse(errUnhandled(err)), nil
	}

	return gen.ListGenres200JSONResponse(resp), nil
}

func (c Controller) ListLanguages(ctx context.Context, request gen.ListLanguagesRequestObject) (gen.ListLanguagesResponseObject, error) {
	resp, err := c.languageService.ListLanguages(ctx, request)
	if err != nil {
		return gen.ListLanguages500JSONResponse(errUnhandled(err)), nil
	}

	return gen.ListLanguages200JSONResponse(resp), nil
}

func (c Controller) UploadIcon(ctx context.Context, request gen.UploadIconRequestObject) (gen.UploadIconResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
