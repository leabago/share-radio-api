package entity

import (
	"github.com/google/uuid"
	"github.com/leabago/share-radio/adder/docs/gen"
)

type Language struct {
	Id   string
	Name string
}

func (g *Language) ConvertToHttpLanguage() *gen.Language {
	id, _ := uuid.Parse(g.Id)

	return &gen.Language{
		Id:   new(id),
		Name: new(g.Name),
	}
}

func ConvertToHttpLanguageList(arr []Language) []gen.Language {
	resp := make([]gen.Language, len(arr))
	for i, v := range arr {
		resp[i] = *v.ConvertToHttpLanguage()
	}

	return resp
}
