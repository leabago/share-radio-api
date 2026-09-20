package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/leabago/share-radio/adder/docs/gen"
)

type Station struct {
	Id         string
	Name       string
	Url        string
	Genre      string `db:"genre_name"`
	GenreId    string `db:"-"`
	Language   string `db:"language_name"`
	LanguageId string `db:"-"`
	Icon       *string
	IsActive   *bool
	IsNew      *bool
	AddedAt    *time.Time
	UpdatedAt  *time.Time
}

func (s Station) ConvertToHttpStation() *gen.Station {

	id, _ := uuid.Parse(s.Id)

	//s.Genre.ConvertToHttpGenre()

	return &gen.Station{
		Id:        new(id),
		Name:      new(s.Name),
		Url:       new(s.Url),
		Genre:     new(s.Genre),
		Language:  new(s.Language),
		Icon:      s.Icon,
		IsActive:  s.IsActive,
		IsNew:     s.IsNew,
		AddedAt:   s.AddedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func ConvertToHttpStationList(stations []Station) []gen.Station {
	getStations := make([]gen.Station, len(stations))

	for i, v := range stations {

		station := v.ConvertToHttpStation()

		getStations[i] = *station
	}

	return getStations
}

//
//func (s Station) ConvertToHttpStationDetail() (gen.StationDetail, error) {
//	id, err := uuid.Parse(s.Id)
//	if err != nil {
//		return gen.StationDetail{}, err
//	}
//
//	return gen.StationDetail{
//		Id:        new(id),
//		Name:      new(s.Name),
//		Url:       new(s.Url),
//		Genre:     new(s.Genre),
//		Language:  new(s.Language),
//		Icon:      s.Icon,
//		IsActive:  s.IsActive,
//		IsNew:     s.IsNew,
//		AddedAt:   s.AddedAt,
//		UpdatedAt: s.UpdatedAt,
//	}, nil
//}

func ConvertCreateStationRequestObject(request gen.CreateStationRequestObject) (Station, error) {

	return Station{
		Name:       request.Body.Name,
		Url:        request.Body.Url,
		GenreId:    request.Body.GenreId.String(),
		LanguageId: request.Body.LanguageId.String(),
	}, nil
}

type StationId struct {
	Id string `db:"id"`
}
