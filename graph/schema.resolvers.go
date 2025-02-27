package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"context"
	"fmt"

	"github.com/staszkiet/DictionaryGolang/database"
	dbmodels "github.com/staszkiet/DictionaryGolang/database/models"
	customerrors "github.com/staszkiet/DictionaryGolang/errors"
	"github.com/staszkiet/DictionaryGolang/graph/model"
)

// CreateWord is the resolver for the createWord field.
func (r *mutationResolver) CreateWord(ctx context.Context, polish string, translation model.NewTranslation) (bool, error) {
	var convertedTranslations []dbmodels.Translation
	var sentences []dbmodels.Sentence

	var count int64

	err := database.DB.Model(&dbmodels.Word{}).Where("polish = ?", polish).Count(&count).Error
	if err != nil {
		return false, err
	} else if count > 0 {
		return false, &customerrors.WordExistsError{Word: polish}
	}

	fmt.Println(count)

	sentences = make([]dbmodels.Sentence, 0)
	for _, s := range translation.Sentences {
		sentences = append(sentences, dbmodels.Sentence{Sentence: s})
	}
	convertedTranslations = append(convertedTranslations, dbmodels.Translation{
		English:   translation.English,
		Sentences: sentences,
	})

	ret := &dbmodels.Word{
		Polish:       polish,
		Translations: convertedTranslations,
	}
	database.DB.Create(ret)
	return true, nil
}

// DeleteWord is the resolver for the deleteWord field.
func (r *mutationResolver) DeleteWord(ctx context.Context, polish string) (string, error) {
	var word dbmodels.Word

	if err := database.DB.Preload("Translations.Sentences").Where("polish = ?", polish).First(&word).Error; err != nil {
		panic(err)
	}

	if err := database.DB.Preload("Translations.Sentences").Delete(&word).Error; err != nil {
		panic(err)
	}

	return polish, nil
}

// SelectWord is the resolver for the selectWord field.
func (r *queryResolver) SelectWord(ctx context.Context, polish string) (*model.Word, error) {
	var word dbmodels.Word

	if err := database.DB.Preload("Translations.Sentences").Where("polish = ?", polish).First(&word).Error; err != nil {
		return nil, err
	}

	return dbmodels.DBWordToGQLWord(&word), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
