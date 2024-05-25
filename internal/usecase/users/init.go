package users

import (
	"github.com/eifzed/makaji/internal/config"
	"github.com/eifzed/makaji/internal/entity/repo/db"
	"github.com/eifzed/makaji/internal/entity/repo/transactions"
	"github.com/eifzed/makaji/internal/entity/repo/users"
)

type usersUC struct {
	usersDB users.UsersDBInterface
	config  *config.Config
	tx      transactions.TransactionInterface
	elastic db.ElasticsearchInterface
}

type Options struct {
	UsersDB users.UsersDBInterface
	Config  *config.Config
	TX      transactions.TransactionInterface
	Elastic db.ElasticsearchInterface
}

func GetNewUsersUC(option *Options) *usersUC {
	if option == nil || option.UsersDB == nil {
		return nil
	}
	return &usersUC{
		usersDB: option.UsersDB,
		config:  option.Config,
		tx:      option.TX,
		elastic: option.Elastic,
	}
}
