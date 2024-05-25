package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/eifzed/joona/internal/entity/users"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/pkg/errors"
)

func (es *elasticSearch) InsertUser(ctx context.Context, data *users.UserItem) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	resp, err := es.client.Create(
		"users",
		data.UserID,
		bytes.NewReader(dataByte),
	)
	if err != nil {
		err = errors.Wrap(err, "InsertUser.Create")
		return
	}

	if resp.IsError() {
		err = errors.Wrap(err, "InsertUser.Create."+resp.String())
		return
	}
	return
}

func (es *elasticSearch) UpdateUser(ctx context.Context, id string, data *users.UserItem) (err error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	fmt.Println(string(dataByte))

	resp, err := es.client.Update(
		"users",
		id,
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"doc":%s}`, dataByte))),
	)
	if err != nil {
		err = errors.Wrap(err, "Update")
		return
	}

	if resp.IsError() {
		fmt.Println(resp.String())
		err = errors.Wrap(err, "Update."+resp.String())
		return
	}
	return
}

func (es *elasticSearch) GetUserList(ctx context.Context, params users.GenericFilterParams) (result users.GetUserListResponse, err error) {
	query := Search{}

	if params.Keyword != "" {
		bool := Bool{}
		bool.Should = append(bool.Should,
			Should{Wildcard: map[string]Wildcard{
				"full_name": {
					Value:           "*" + params.Keyword + "*",
					Boost:           1.2,
					CaseInsensitive: true,
				},
			}},
			Should{Match: map[string]string{
				"username": params.Keyword,
			}},
		)
		bool.MinimumShouldMatch = 1
		query.Query.Bool = &bool
	}

	query.SetPagination(params.Page, params.Limit)

	qb, err := json.Marshal(query)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal")
		return
	}

	resp, err := es.client.Search(
		es.client.Search.WithContext(ctx),
		es.client.Search.WithIndex("users"),
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

	repBody, _ := io.ReadAll(resp.Body)

	if err = json.Unmarshal(repBody, &searchResult); err != nil {
		err = errors.Wrap(err, "json.NewDecoder")
		return
	}

	recipeData := []users.UserItem{}

	total, err := bindResult(&searchResult, &recipeData)
	if err != nil {
		err = errors.Wrap(err, "bindResult")
		return
	}
	result.Data = recipeData
	result.Total = total

	return
}
