package elasticsearch

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/pkg/errors"
)

func bindResult(searchResult *search.Response, data interface{}) (total int64, err error) {
	if searchResult == nil || len(searchResult.Hits.Hits) == 0 {
		return 0, nil
	}

	combined := "["

	for i, hits := range searchResult.Hits.Hits {
		combined += string(hits.Source_)
		if i < len(searchResult.Hits.Hits)-1 {
			combined += ","
		}
	}
	combined += "]"

	err = json.Unmarshal([]byte(combined), data)
	if err != nil {
		err = errors.Wrap(err, "json.Unmarshal."+combined)
		return
	}
	total = searchResult.Hits.Total.Value
	return
}
