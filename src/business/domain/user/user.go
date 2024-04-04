package user

import (
	"context"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/parser"
	"github.com/adiatma85/own-go-sdk/sql"
)

type Interface interface {
	Create(ctx context.Context, userParam entity.CreateUserParam) (entity.User, error)
	Get(ctx context.Context, params entity.UserParam) (entity.User, error)
	GetList(ctx context.Context, params entity.UserParam) ([]entity.User, *entity.Pagination, error)
	Update(ctx context.Context, updateParam entity.UpdateUserParam, selectParam entity.UserParam) error
}

type InitParam struct {
	Log  log.Interface
	Db   sql.Interface
	Json parser.JSONInterface
}

type user struct {
	log  log.Interface
	db   sql.Interface
	json parser.JSONInterface
}

func Init(param InitParam) Interface {
	u := &user{
		log:  param.Log,
		db:   param.Db,
		json: param.Json,
	}

	return u
}

func (u *user) Create(ctx context.Context, userParam entity.CreateUserParam) (entity.User, error) {
	user := entity.User{}

	tx, err := u.db.Leader().BeginTx(ctx, "txcUser", sql.TxOptions{})
	if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	tx, user, err = u.createSQLUser(tx, userParam)
	if err != nil {
		return user, err
	}

	if err = tx.Commit(); err != nil {
		return user, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	return u.Get(ctx, entity.UserParam{
		ID: null.Int64From(user.ID),
	})
}

func (u *user) Get(ctx context.Context, params entity.UserParam) (entity.User, error) {
	return u.getSQLUser(ctx, params)
}

func (u *user) GetList(ctx context.Context, params entity.UserParam) ([]entity.User, *entity.Pagination, error) {
	return u.getSQLUserList(ctx, params)
}

func (u *user) Update(ctx context.Context, updateParam entity.UpdateUserParam, selectParam entity.UserParam) error {
	return u.updateSQLUser(ctx, updateParam, selectParam)
}
