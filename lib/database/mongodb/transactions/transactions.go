package transactions

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionKey struct{}

var (
	txKey = TransactionKey{}
)

type transaction struct {
	Client *mongo.Client
}

var dbName string

type Options struct {
	Client *mongo.Client
	DBName string
}

func GetNewMongoDBTransaction(option *Options) *transaction {
	if option == nil || option.DBName == "" {
		return nil
	}
	dbName = option.DBName
	return &transaction{
		Client: option.Client,
	}
}

func (tx *transaction) Start(ctx context.Context) (context.Context, error) {
	session, err := tx.Client.StartSession()
	if err != nil {
		return nil, err
	}
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	ctx = setSessionToContext(session, ctx)
	return ctx, nil
}

func setSessionToContext(session mongo.Session, ctx context.Context) context.Context {
	newContext := context.WithValue(ctx, txKey, session)
	return newContext
}

func GetSessionFromContext(ctx context.Context) mongo.Session {
	session, ok := ctx.Value(txKey).(mongo.Session)
	if !ok {
		return nil
	}
	return session
}

func GetCollectionFromSession(session mongo.Session, collectionName string) *mongo.Collection {
	return session.Client().Database(dbName).Collection(collectionName)
}

func (tx *transaction) Finish(ctx context.Context, err *error) {
	session := GetSessionFromContext(ctx)
	var errOrigin error

	if p := recover(); p != nil {
		_ = session.AbortTransaction(ctx)
		panic(p)
	}
	if err != nil {
		errOrigin = *err
	}
	if errOrigin == nil {
		fmt.Println("committing tx")
		_ = session.CommitTransaction(ctx)
	} else {
		fmt.Println("aborting tx")
		_ = session.AbortTransaction(ctx)
	}
}
