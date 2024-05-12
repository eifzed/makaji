package recipes

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/repo/db"
	"github.com/eifzed/joona/internal/entity/repo/transactions"
	"github.com/eifzed/joona/internal/entity/repo/users"
)

type recipesUC struct {
	usersDB   users.UsersDBInterface
	recipesDB db.RecipesDBInterface
	config    *config.Config
	tx        transactions.TransactionInterface
	elastic   db.ElasticsearchInterface
	redis     db.RedisInterface
}

type Options struct {
	UsersDB   users.UsersDBInterface
	Config    *config.Config
	TX        transactions.TransactionInterface
	RecipesDB db.RecipesDBInterface
	Elastic   db.ElasticsearchInterface
	Redis     db.RedisInterface
}

func GetNewRecipesUC(option *Options) *recipesUC {
	if option == nil || option.UsersDB == nil {
		return nil
	}
	return &recipesUC{
		usersDB:   option.UsersDB,
		config:    option.Config,
		tx:        option.TX,
		recipesDB: option.RecipesDB,
		elastic:   option.Elastic,
		redis:     option.Redis,
	}
}
