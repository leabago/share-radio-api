package language

import (
	"context"

	"github.com/leabago/share-radio/adder/docs/gen"
	"github.com/leabago/share-radio/adder/internal/entity"
	"github.com/leabago/share-radio/adder/internal/repo"
)

type LanguageCase struct {
	repo repo.LanguageRepo
}

func New(r repo.LanguageRepo) *LanguageCase {
	return &LanguageCase{
		repo: r,
	}
}

func (g *LanguageCase) ListLanguages(ctx context.Context, request gen.ListLanguagesRequestObject) (gen.LanguageList, error) {

	languages, err := g.repo.ListLanguages(ctx)
	if err != nil {
		return gen.LanguageList{}, err
	}

	resp := gen.LanguageList{
		Languages: new(entity.ConvertToHttpLanguageList(languages)),
	}

	return resp, nil
}
