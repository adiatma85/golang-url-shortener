package user

import (
	"context"
	"fmt"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/query"
	"github.com/adiatma85/own-go-sdk/sql"
)

func (u *user) createSQLUser(tx sql.CommandTx, v entity.CreateUserParam) (sql.CommandTx, entity.User, error) {
	user := entity.User{}

	res, err := tx.NamedExec("iCreateUser", createUser, v)
	if err != nil {
		return tx, user, errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil || rowCount < 1 {
		return tx, user, errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no rows affected")
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return tx, user, errors.NewWithCode(codes.CodeSQLNoRowsAffected, err.Error())
	}

	user.ID = lastID

	return tx, user, nil
}

func (u *user) getSQLUser(ctx context.Context, params entity.UserParam) (entity.User, error) {
	user := entity.User{}

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &params.QueryOption)
	queryExt, queryArgs, _, _, err := qb.Build(&params)
	if err != nil {
		return user, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	row, err := u.db.Follower().QueryRow(ctx, "rUserByID", readUser+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return user, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	} else if errors.Is(err, sql.ErrNotFound) {
		return user, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	}

	if err := row.StructScan(&user); err != nil && !errors.Is(err, sql.ErrNotFound) {
		return user, errors.NewWithCode(codes.CodeSQLRowScan, err.Error())
	} else if errors.Is(err, sql.ErrNotFound) {
		return user, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, err.Error())
	}

	return user, nil
}

func (u *user) getSQLUserList(ctx context.Context, params entity.UserParam) ([]entity.User, *entity.Pagination, error) {
	users := []entity.User{}

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &params.QueryOption)
	queryExt, queryArgs, countExt, countArgs, err := qb.Build(&params)
	if err != nil {
		return users, nil, errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	rows, err := u.db.Follower().Query(ctx, "rListUser", readUser+queryExt, queryArgs...)
	if err != nil && !errors.Is(err, sql.ErrNotFound) {
		return users, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		temp := entity.User{}
		if err := rows.StructScan(&temp); err != nil {
			u.log.Error(ctx, errors.NewWithCode(codes.CodeSQLRowScan, err.Error()))
			return users, nil, err
		}
		users = append(users, temp)
	}

	pg := entity.Pagination{
		CurrentPage:     params.Page,
		CurrentElements: int64(len(users)),
	}

	if len(users) > 0 && !params.QueryOption.DisableLimit && params.IncludePagination {
		if err := u.db.Follower().Get(ctx, "cUser", readUserCount+countExt, &pg.TotalElements, countArgs...); err != nil {
			return users, nil, errors.NewWithCode(codes.CodeSQLRead, err.Error())
		}
	}

	pg.ProcessPagination(params.Limit)

	return users, &pg, nil
}

func (u *user) updateSQLUser(ctx context.Context, updateParam entity.UpdateUserParam, selectParam entity.UserParam) error {
	u.log.Debug(ctx, fmt.Sprintf("update user profile by: %v", selectParam))

	qb := query.NewSQLQueryBuilder(u.db, "param", "db", &selectParam.QueryOption)

	var err error
	queryUpdate, args, err := qb.BuildUpdate(&updateParam, &selectParam)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLBuilder, err.Error())
	}

	res, err := u.db.Leader().Exec(ctx, "uProfile", updateUser+queryUpdate, args...)
	if err != nil {
		return errors.NewWithCode(codes.CodeSQLTxExec, err.Error())
	}

	rowCount, err := res.RowsAffected()
	if err != nil || rowCount < 1 {
		return errors.NewWithCode(codes.CodeSQLNoRowsAffected, "no rows affected")
	}

	u.log.Debug(ctx, fmt.Sprintf("successfully updated user: %v", updateParam))

	return nil
}
