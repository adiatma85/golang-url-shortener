package url

import (
	"context"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/parser"
	"github.com/adiatma85/own-go-sdk/sql"
)

type Interface interface {
	Create(ctx context.Context, insertParam entity.CreateUrlParam) (entity.Url, error)
	Get(ctx context.Context, params entity.UrlParam) (entity.Url, error)
	GetList(ctx context.Context, params entity.UrlParam) ([]entity.Url, *entity.Pagination, error)
	Update(ctx context.Context, updateParam entity.UpdateUrlParam, selectParam entity.UrlParam) error
}

type InitParam struct {
	Log  log.Interface
	Db   sql.Interface
	Json parser.JSONInterface
}

type url struct {
	log  log.Interface
	db   sql.Interface
	json parser.JSONInterface
}

func Init(params InitParam) Interface {
	u := &url{
		log:  params.Log,
		db:   params.Db,
		json: params.Json,
	}

	return u
}

func (u *url) Create(ctx context.Context, insertParam entity.CreateUrlParam) (entity.Url, error) {
	result := entity.Url{}

	tx, err := u.db.Leader().BeginTx(ctx, "txcUser", sql.TxOptions{})
	if err != nil {
		return result, errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	tx, result, err = u.createSQLUrl(tx, insertParam)
	if err != nil {
		return result, err
	}

	if err = tx.Commit(); err != nil {
		return result, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	return result, nil
}

func (u *url) Get(ctx context.Context, params entity.UrlParam) (entity.Url, error) {
	return u.getSQLUrl(ctx, params)
}

func (u *url) GetList(ctx context.Context, params entity.UrlParam) ([]entity.Url, *entity.Pagination, error) {
	return u.getSQLUrlList(ctx, params)
}

func (u *url) Update(ctx context.Context, updateParam entity.UpdateUrlParam, selectParam entity.UrlParam) error {
	return u.updateSQLUrl(ctx, updateParam, selectParam)
}
