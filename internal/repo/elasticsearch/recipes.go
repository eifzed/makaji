package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/eifzed/joona/internal/entity/recipes"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/pkg/errors"
)

func (es *elasticSearch) GetRecipeList(ctx context.Context, params recipes.GetRecipeParams) (result recipes.GetRecipeListResponse, err error) {
	query := Search{}

	if params.Keyword != "" {
		query.Query.Bool.Should = append(query.Query.Bool.Should,
			Should{Match: map[string]string{
				"description": params.Keyword,
			}},
			Should{Wildcard: map[string]Wildcard{
				"name": {
					Value:           "*" + params.Keyword + "*",
					Boost:           1.2,
					CaseInsensitive: true,
				},
			}},
			Should{Match: map[string]string{
				"tags": params.Keyword,
			}},
		)
		query.Query.Bool.MinimumShouldMatch = 1
	}

	query.SetPagination(params.Page, params.Limit)

	qb, err := json.Marshal(query)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal")
		return
	}

	resp, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex("recipes"),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithBody(bytes.NewReader(qb)),
	)
	if err != nil {
		err = errors.Wrap(err, "Search")
		return
	}

	if resp.IsError() {
		err = errors.Wrap(err, "Search."+resp.String())
		return
	}

	defer resp.Body.Close()
	var searchResult search.Response

	if err = json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		err = errors.Wrap(err, "json.NewDecoder")
		return
	}

	recipeData := []recipes.ReceipeItem{}

	total, err := bindResult(&searchResult, &recipeData)
	if err != nil {
		err = errors.Wrap(err, "bindResult")
		return
	}
	result.Items = recipeData
	result.Total = total

	return
}

func (es *elasticSearch) InsertRecipe(ctx context.Context, data *recipes.ReceipeItem) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	resp, err := es.client.Index(
		"recipes",
		bytes.NewReader(dataByte),
	)
	if err != nil {
		err = errors.Wrap(err, "Index")
	}

	if resp.IsError() {
		err = errors.Wrap(err, "Search."+resp.String())
		return
	}
	return
}

func (es *elasticSearch) UpdateRecipe(ctx context.Context, id string, data *recipes.ReceipeItem) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	resp, err := es.client.Update(
		"recipes",
		id,
		bytes.NewBuffer(dataByte),
	)
	if err != nil {
		err = errors.Wrap(err, "Update")
	}

	if resp.IsError() {
		err = errors.Wrap(err, "Update."+resp.String())
		return
	}
	return
}
