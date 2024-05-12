package url

import (
	"context"
	"fmt"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/query"
	"github.com/adiatma85/own-go-sdk/sql"
)

func (u *url) createSQLUrl(tx sql.CommandTx, v entity.CreateUrlParam) (sql.CommandTx, entity.Url, error) {
	result := entity.Url{}

	res, err := tx.NamedExec("iCreateUrl", createUrl, v)
	if err != nil {
		return tx, result, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil || rowCount < 1 {
		return tx, result, errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no rows affected")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return tx, result, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	}

	result.ID = lastID

	return tx, result, nil
}

func (u *url) getSQLUrl(ctx context.Context, params entity.UrlParam) (entity.Url, error) {
	result := entity.Url{}

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &params.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&params)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := u.db.Follower().QueryRow(ctx, "rUrlByID", getUrl+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return result, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	} else if errors.Is(err, sql.ErrNotFound) {
		return result, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	}

	if err := row.StructScan(&result); err != nil && !errors.Is(err, sql.ErrNotFound) {
		return result, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	} else if errors.Is(err, sql.ErrNotFound) {
		return result, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	}

	return result, nil
}

func (u *url) getSQLUrlList(ctx context.Context, params entity.UrlParam) ([]entity.Url, *entity.Pagination, error) {
	results := []entity.Url{}

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &params.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&params)
	if err != nil {
		return results, nil, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := u.db.Follower().Query(ctx, "rListUrl", getUrl+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return results, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		temp := entity.Url{}
		if err := rows.StructScan(&temp); err != nil {
			u.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			continue
		}
		results = append(results, temp)
	}

	pg := entity.Pagination{
		CurrentPage:     params.Page,
		CurrentElements: int64(len(results)),
	}

	if len(results) > 0 && !params.QueryOption.DisableLimit && params.IncludePagination {
		if err := u.db.Follower().Get(ctx, "cUrl", readUrlCount+countExt, &pg.TotalElements, countArgs...); err != nil {
			return results, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(params.Limit)

	return results, &pg, nil
}

func (u *url) updateSQLUrl(ctx context.Context, updateParam entity.UpdateUrlParam, selectParam entity.UrlParam) error {
	u.log.Debug(ctx, fmt.Sprintf("update url data entry by: %v", selectParam))

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &selectParam.QueryOption)

	var err error
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	_, err = u.db.Leader().Exec(ctx, "uProfile", updateUrl+queryUpdate, args...)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	u.log.Debug(ctx, fmt.Sprintf("successfully updated url: %v", updateParam))

	return nil
}
