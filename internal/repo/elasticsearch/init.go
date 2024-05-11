package elasticsearch

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
)

type elasticSearch struct {
	config *config.Config
	client *elasticsearch.Client
}

type Option struct {
	APIKey  string
	CloudID string
	Config  *config.Config
}

func New(opt Option) (client *elasticSearch, err error) {
	esCfg := elasticsearch.Config{
		CloudID: opt.CloudID,
		APIKey:  opt.APIKey,
	}
	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		err = errors.Wrap(err, "elasticsearch.NewClient")
		return
	}

	// boolQ
	// esClient.Search().Index("recipes").Request(&search.Request{Query: })
	client = &elasticSearch{
		config: opt.Config,
		client: esClient,
	}
	return
}
