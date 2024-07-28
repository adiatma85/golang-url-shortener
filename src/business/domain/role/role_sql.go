package role

import (
	"context"
	"fmt"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/query"
	"github.com/adiatma85/own-go-sdk/sql"
)

func (r *role) createSQLRole(tx sql.CommandTx, v entity.CreateRoleParam) (sql.CommandTx, entity.Role, error) {
	role := entity.Role{}

	res, err := tx.NamedExec("iCreateRole", createRole, v)
	if err != nil {
		return tx, role, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil || rowCount < 1 {
		return tx, role, errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no rows affected")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return tx, role, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	}

	role.ID = lastID

	return tx, role, nil
}

func (r *role) getSQLRole(ctx context.Context, params entity.RoleParam) (entity.Role, error) {
	result := entity.Role{}

	qb := query.NewSQLQueryBuilder(r.db, "param", "db", &params.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&params)
	if err != nil {
		return result, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := r.db.Follower().QueryRow(ctx, "rRoleByID", readRole+queryExt, queryArgs...)
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

func (r *role) getSQLRoleList(ctx context.Context, params entity.RoleParam) ([]entity.Role, *entity.Pagination, error) {
	results := []entity.Role{}

	qb := query.NewSQLQueryBuilder(r.db, "param", "db", &params.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&params)
	if err != nil {
		return results, nil, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := r.db.Follower().Query(ctx, "rListRole", readRole+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return results, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		temp := entity.Role{}
		if err := rows.StructScan(&temp); err != nil {
			r.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			// Need discussiion, ini jadinya mau bagaimana
			return results, nil, err
		}
		results = append(results, temp)
	}

	pg := entity.Pagination{
		CurrentPage:     params.Page,
		CurrentElements: int64(len(results)),
	}

	if len(results) > 0 && !params.QueryOption.DisableLimit && params.IncludePagination {
		if err := r.db.Follower().Get(ctx, "cRole", readRoleCount+countExt, &pg.TotalElements, countArgs...); err != nil {
			return results, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(params.Limit)

	return results, &pg, nil
}

func (r *role) updateSQLRole(ctx context.Context, updateParam entity.UpdateRoleParam, selectParam entity.RoleParam) error {
	r.log.Debug(ctx, fmt.Sprintf("update role by: %v", selectParam))

	qb := query.NewSQLQueryBuilder(r.db, "param", "db", &selectParam.QueryOption)

	var err error
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := r.db.Leader().Exec(ctx, "uRole", updateRole+queryUpdate, args...)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil || rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no rows affected")
	}

	r.log.Debug(ctx, fmt.Sprintf("successfully updated role: %v", updateParam))

	return nil
}
