package xorm

import (
	"context"

	"github.com/pkg/errors"

	"github.com/go-xorm/xorm"
)

type Config struct {
	DSN                  string `json:"dsn" yaml:"dsn"`
	EnableSSL            bool   `json:"enable_ssl" yaml:"enable_ssl"`
	CertificateWritePath string `json:"certificate_write_path" yaml:"certificate_write_path"`
	ClientKey            string `json:"client_key" yaml:"client_key"`
	ClientSecret         string `json:"client_secret" yaml:"client_secret"`
	ServerCA             string `json:"server_ca" yaml:"server_ca"`
}

type DBTransaction struct {
	*Connection
}

type Connection struct {
	Master *xorm.Engine
	Slave  *xorm.Engine
}

type txContext struct{}

func ConnetDB(dataSource string) (*xorm.Engine, error) {
	return xorm.NewEngine("postgres", dataSource)
}

func GetDBTransaction(conn *Connection) *DBTransaction {
	return &DBTransaction{
		Connection: conn,
	}
}

func (dt *DBTransaction) Begin(ctx context.Context) (context.Context, error) {
	session := dt.Master.NewSession().Context(ctx)
	err := session.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "Begin Transaction")
	}
	return context.WithValue(ctx, txContext{}, session), nil
}

func (dt *DBTransaction) GetDB(ctx context.Context) xorm.Interface {
	if ctx == nil {
		return nil
	}
	session, ok := ctx.Value(txContext{}).(*xorm.Session)
	if !ok {
		return dt.Master.Context(ctx)
	}
	return session
}

func GetDBSession(ctx context.Context) *xorm.Session {
	if ctx == nil {
		return nil
	}
	session, ok := ctx.Value(txContext{}).(*xorm.Session)
	if !ok {
		return nil
	}
	return session
}

func (dt *DBTransaction) commit(session *xorm.Session) error {
	err := session.Commit()
	if err != nil {
		return errors.Wrap(err, "Commit Transaction")
	}
	return nil
}

func (dt *DBTransaction) rollback(session *xorm.Session) error {
	err := session.Rollback()
	if err != nil {
		return errors.Wrap(err, "Rollback Transaction")
	}
	return nil
}

func (dt *DBTransaction) Finish(ctx context.Context, err *error) {
	var errValue error

	if err != nil {
		errValue = *err
	}
	session, ok := dt.GetDB(ctx).(*xorm.Session)
	p := recover()
	if p != nil {
		_ = dt.rollback(session)
		panic(p)
	}
	if !ok {
		panic("Finish must be called on a started transaction")
	}
	if errValue != nil {
		_ = dt.rollback(session)
	} else {
		_ = dt.commit(session)
	}
	if session != nil {
		session.Close()
	}
}
