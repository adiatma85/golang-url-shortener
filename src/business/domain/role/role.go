package role

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
	Create(ctx context.Context, insertParam entity.CreateRoleParam) (entity.Role, error)
	Get(ctx context.Context, params entity.RoleParam) (entity.Role, error)
	GetList(ctx context.Context, params entity.RoleParam) ([]entity.Role, *entity.Pagination, error)
	Update(ctx context.Context, updateParam entity.UpdateRoleParam, selectParam entity.RoleParam) error
}

type InitParam struct {
	Log  log.Interface
	Db   sql.Interface
	Json parser.JSONInterface
}

type role struct {
	log  log.Interface
	db   sql.Interface
	json parser.JSONInterface
}

func Init(param InitParam) Interface {
	r := &role{
		log:  param.Log,
		db:   param.Db,
		json: param.Json,
		// redis: param.Redis,
	}

	return r
}

func (r *role) Create(ctx context.Context, insertParam entity.CreateRoleParam) (entity.Role, error) {
	result := entity.Role{}

	tx, err := r.db.Leader().BeginTx(ctx, "txcRole", sql.TxOptions{})
	if err != nil {
		return result, errors.NewWithCode(codes.CodeSQLTxBegin, err.Error())
	}
	defer tx.Rollback()

	tx, result, err = r.createSQLRole(tx, insertParam)
	if err != nil {
		return result, err
	}

	if err = tx.Commit(); err != nil {
		return result, errors.NewWithCode(codes.CodeSQLTxCommit, err.Error())
	}

	return r.Get(ctx, entity.RoleParam{
		ID: null.Int64From(result.ID),
	})
}

func (r *role) Get(ctx context.Context, params entity.RoleParam) (entity.Role, error) {
	return r.getSQLRole(ctx, params)
}

func (r *role) GetList(ctx context.Context, params entity.RoleParam) ([]entity.Role, *entity.Pagination, error) {
	return r.getSQLRoleList(ctx, params)
}

func (r *role) Update(ctx context.Context, updateParam entity.UpdateRoleParam, selectParam entity.RoleParam) error {
	return r.updateSQLRole(ctx, updateParam, selectParam)
}
