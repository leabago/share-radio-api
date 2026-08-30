package station

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/leabago/share-radio/adder/config"
	"github.com/leabago/share-radio/adder/docs/gen"
	"github.com/leabago/share-radio/adder/internal/entity"
	"github.com/leabago/share-radio/adder/internal/repo"
	"github.com/leabago/share-radio/adder/pkg/helper/chronos"
)

// UseCase -.
type StationCase struct {
	repo repo.StationRepo
	cfg  *config.Config
}

// New returns a Task usecase instrumented with OpenTelemetry tracing spans.
func New(cfg *config.Config, r repo.StationRepo) *StationCase {
	return &StationCase{
		repo: r,
		cfg:  cfg,
	}
}

// Create -.
func (sc *StationCase) CreateStation(ctx context.Context, request gen.CreateStationRequestObject) (gen.StationId, error) {

	station, err := entity.ConvertCreateStationRequestObject(request)
	if err != nil {
		return gen.StationId{}, err
	}

	id, err := sc.repo.CreateStation(ctx, &station)
	if err != nil {
		return gen.StationId{}, err
	}

	idUuid, err := uuid.Parse(id)
	if err != nil {
		return gen.StationId{}, err
	}

	jwtToken, err := sc.createSessionToken(request.Params.XSessionId, id)
	if err != nil {
		return gen.StationId{}, err
	}

	resp := gen.StationId{
		Id:          new(idUuid),
		AccessToken: new(jwtToken),
		ExpiresIn:   new(60 * 60 * 24),
		TokenType:   new("Bearer"),
	}

	return resp, nil
}

func (sc *StationCase) createSessionToken(sessionId string, stationId string) (string, error) {

	claims := entity.SessionClaims{
		SessionId: sessionId,
		StationId: stationId,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(chronos.Now().AddDays(1).Time),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "adder-app",
			Subject:   "user-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(sc.cfg.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (sc *StationCase) ListStations(ctx context.Context, request gen.ListStationsRequestObject) (gen.StationListResponse, error) {

	stations, err := sc.repo.ListStations(ctx, request)
	if err != nil {
		return gen.StationListResponse{}, err
	}

	getStations := entity.ConvertToHttpStationList(stations)

	resp := gen.StationListResponse{
		Stations: new(getStations),
		Pagination: &gen.Pagination{
			Limit:  request.Params.Limit,
			Offset: request.Params.Offset,
			Total:  new(len(getStations)),
		},
	}

	return resp, nil
}

func (sc *StationCase) GetStation(ctx context.Context, request gen.GetStationRequestObject) (gen.Station, error) {
	stationEntity, err := sc.repo.GetStation(ctx, request.Id.String())
	if err != nil {
		return gen.Station{}, err
	}

	station := stationEntity.ConvertToHttpStation()

	return *station, nil
}
