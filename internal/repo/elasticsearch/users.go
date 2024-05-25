package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/eifzed/makaji/internal/entity/users"
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
