package users

import (
	"github.com/eifzed/joona/internal/config"
	"github.com/eifzed/joona/internal/entity/repo/transactions"
	"github.com/eifzed/joona/internal/entity/repo/users"
)

type usersUC struct {
	UsersDB users.UsersDBInterface
	Config  *config.Config
	TX      transactions.TransactionInterface
}

type Options struct {
	UsersDB users.UsersDBInterface
	Config  *config.Config
	TX      transactions.TransactionInterface
}

func GetNewUsersUC(option *Options) *usersUC {
	if option == nil || option.UsersDB == nil {
		return nil
	}
	return &usersUC{
		UsersDB: option.UsersDB,
		Config:  option.Config,
		TX:      option.TX,
	}
}
