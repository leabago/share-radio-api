package entity

import (
	"github.com/google/uuid"
	"github.com/leabago/share-radio/adder/docs/gen"
)

type Genre struct {
	Id   string
	Name string
}

func (g *Genre) ConvertToHttpGenre() *gen.Genre {
	id, _ := uuid.Parse(g.Id)

	return &gen.Genre{
		Id:   new(id),
		Name: new(g.Name),
	}
}

func ConvertToHttpGenreList(arr []Genre) []gen.Genre {
	resp := make([]gen.Genre, len(arr))
	for i, v := range arr {
		resp[i] = *v.ConvertToHttpGenre()
	}

	return resp
}
