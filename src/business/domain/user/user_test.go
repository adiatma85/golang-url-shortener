package user

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
	libsql "github.com/adiatma85/own-go-sdk/sql"
	mock_log "github.com/adiatma85/own-go-sdk/tests/mock/log"
	mock_json "github.com/adiatma85/own-go-sdk/tests/mock/parser"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func Test_user_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockJsonParser := mock_json.NewMockJSONInterface(ctrl)

	query := regexp.QuoteMeta(`INSERT INTO user (fk_role_id, email, username, password, display_name, created_by) VALUES (?, ?, ?, ?, ?, ?)`)
	queryGet := regexp.QuoteMeta(readUser)

	// Type in here
	type args struct {
		ctx         context.Context
		createParam entity.CreateUserParam
	}

	// Mock in here
	mockCreateParam := entity.CreateUserParam{
		RoleId:          1,
		Username:        "Ramdani Koernia",
		Email:           "random@gmail.com",
		Password:        "strongPassword",
		ConfirmPassword: "strongPassword",
		DisplayName:     "Display Name",
	}

	// Test cases in here
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "cannot begin tx",
			args: args{
				ctx:         context.Background(),
				createParam: mockCreateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "cannot exec user",
			args: args{
				ctx:         context.Background(),
				createParam: mockCreateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()

				sqlMock.ExpectExec(query).WillReturnError(errors.NewWithCode(codes.CodeSQL, "cannot create user"))
				sqlMock.ExpectRollback()

				return sqlServer, err
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "user no new row",
			args: args{
				ctx:         context.Background(),
				createParam: mockCreateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0))
				sqlMock.ExpectRollback()
				return sqlServer, err
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "cannot commit to the database",
			args: args{
				ctx:         context.Background(),
				createParam: mockCreateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit().WillReturnError(errors.NewWithCode(codes.CodeSQLTxCommit, "failed to commit"))
				sqlMock.ExpectRollback()
				return sqlServer, err
			},
			want: entity.User{
				ID: 1,
			},
			wantErr: true,
		},
		{
			name: "all good",
			args: args{
				ctx:         context.Background(),
				createParam: mockCreateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				sqlMock.ExpectCommit()

				// Add new rows
				row := sqlMock.NewRows([]string{
					"id",
					"username",
					"email",
				})
				row.AddRow("1", "Ramdani Koernia", "random@gmail.com")
				sqlMock.ExpectQuery(queryGet).WithArgs(1).WillReturnRows(row)

				sqlMock.ExpectRollback()
				return sqlServer, err
			},
			want: entity.User{
				ID:       1,
				Username: "Ramdani Koernia",
				Email:    "random@gmail.com",
			},
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger, nil)

			// Initialize the Domain
			domain := Init(InitParam{
				Log:  logger,
				Db:   sqlClient,
				Json: mockJsonParser,
			})

			got, err := domain.Create(tt.args.ctx, tt.args.createParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("domain.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockJsonParser := mock_json.NewMockJSONInterface(ctrl)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UserParam
	}

	// Mock in here
	now := time.Now()
	query := regexp.QuoteMeta(readUser)

	mockParam := entity.UserParam{
		ID: null.Int64From(1),
	}

	sampleUser := entity.User{
		ID:        1,
		CreatedAt: null.TimeFrom(now),
		CreatedBy: null.StringFrom("test"),
		UpdatedAt: null.TimeFrom(now),
		UpdatedBy: null.StringFrom("test"),
	}

	// Test cases in here
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "get empty row",
			args: args{
				ctx:    context.Background(),
				params: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WithArgs(1).WillReturnError(libsql.ErrNotFound)

				return sqlServer, err
			},
			wantErr: true,
			want:    entity.User{},
		},
		{
			name: "error struct scan",
			args: args{
				ctx:    context.Background(),
				params: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "created_at", "created_by", "updated_at", "updated_by"})
				row.AddRow("A", now, "test", now, "test")

				sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(row)

				return sqlServer, err
			},
			wantErr: true,
			want:    entity.User{},
		},
		{
			name: "all good",
			args: args{
				ctx:    context.Background(),
				params: mockParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "created_at", "created_by", "updated_at", "updated_by"})
				row.AddRow("1", now, "test", now, "test")
				sqlMock.ExpectQuery(query).WithArgs(null.Int64{Int64: 1, Valid: true}).WillReturnRows(row)

				return sqlServer, err
			},
			wantErr: false,
			want:    sampleUser,
		},
	}

	// Iterate the test in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger, nil)

			// Initialize the Domain
			domain := Init(InitParam{
				Log:  logger,
				Db:   sqlClient,
				Json: mockJsonParser,
			})

			got, err := domain.Get(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("domain.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockJsonParser := mock_json.NewMockJSONInterface(ctrl)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UserParam
	}

	// Mock in here
	now := time.Now()
	queryExt := " WHERE 1=1 AND id=? LIMIT 0, 10;"
	query := regexp.QuoteMeta(readUser + queryExt)
	queryCountExt := " WHERE 1=1 AND id=?;"
	queryCount := regexp.QuoteMeta(readUserCount + queryCountExt)

	mockUserParam := entity.UserParam{
		ID: null.Int64From(1),
		PaginationParam: entity.PaginationParam{
			IncludePagination: true,
		},
	}

	mockPagination := entity.Pagination{
		CurrentPage:     1,
		CurrentElements: 1,
		TotalPages:      1,
		TotalElements:   1,
		SortBy:          []string{},
	}

	// Test cases in here
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        []entity.User
		want1       *entity.Pagination
		wantErr     bool
	}{
		{
			name: "error when query-ing",
			args: args{
				ctx:    context.Background(),
				params: mockUserParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(fmt.Errorf("failed to get list of user"))
				return sqlServer, err
			},
			want:    []entity.User{},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "error when struct scan",
			args: args{
				ctx:    context.Background(),
				params: mockUserParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "created_at", "created_by", "updated_at", "updated_by"})
				row.AddRow("A", now, "test", now, "test")
				// error scan here
				sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(row)

				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(0)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			want:    []entity.User{},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "all good",
			args: args{
				ctx:    context.Background(),
				params: mockUserParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				row := sqlmock.NewRows([]string{"id", "created_at", "created_by", "updated_at", "updated_by"})
				row.AddRow("1", now, "test", now, "test")
				// error scan here
				sqlMock.ExpectQuery(query).WithArgs(1).WillReturnRows(row)

				rowCount := sqlMock.NewRows([]string{"COUNT(*)"}).AddRow(1)
				sqlMock.ExpectQuery(queryCount).WillReturnRows(rowCount)
				return sqlServer, err
			},
			want: []entity.User{
				{
					ID:        1,
					CreatedAt: null.TimeFrom(now),
					CreatedBy: null.StringFrom("test"),
					UpdatedAt: null.TimeFrom(now),
					UpdatedBy: null.StringFrom("test"),
				},
			},
			want1:   &mockPagination,
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger, nil)

			// Initialize the Domain
			domain := Init(InitParam{
				Log:  logger,
				Db:   sqlClient,
				Json: mockJsonParser,
			})

			got, pagination, err := domain.GetList(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.Getlist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("domain.Getlist() = %v, want %v", got, tt.want)
			}

			if !reflect.DeepEqual(pagination, tt.want1) {
				t.Errorf("domain.Getlist() = %v, want1 %v", got, tt.want)
			}
		})
	}
}

func Test_user_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockJsonParser := mock_json.NewMockJSONInterface(ctrl)

	// Type in here
	type args struct {
		ctx         context.Context
		updateParam entity.UpdateUserParam
		selectParam entity.UserParam
	}

	// Mock in here
	queryUpdate := regexp.QuoteMeta("UPDATE user SET username=?, display_name=?, updated_by=? WHERE 1=1 AND status=1 AND id=?")

	selectParamSample := entity.UserParam{
		ID: null.Int64From(1),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	updateParamSample := entity.UpdateUserParam{
		Username:    "username",
		DisplayName: "display name",
		UpdatedBy:   null.StringFrom("1"),
	}

	// Test cases in here
	tests := []struct {
		name        string
		prepSqlMock func() (*sql.DB, error)
		args        args
		wantErr     bool
	}{
		{
			name: "failed to exec update",
			args: args{
				ctx:         context.Background(),
				updateParam: updateParamSample,
				selectParam: selectParamSample,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectExec(queryUpdate).WillReturnError(assert.AnError)
				return sqlServer, err
			},
			wantErr: true,
		},
		{
			name: "no user updated",
			args: args{
				ctx:         context.Background(),
				updateParam: updateParamSample,
				selectParam: selectParamSample,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectExec(queryUpdate).WillReturnResult(driver.ResultNoRows)
				return sqlServer, err
			},
			wantErr: true,
		},
		{
			name: "update user 0 rows affected",
			args: args{
				ctx:         context.Background(),
				updateParam: updateParamSample,
				selectParam: selectParamSample,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectExec(queryUpdate).WillReturnResult(driver.RowsAffected(0))
				return sqlServer, err
			},
			wantErr: true,
		},
		{
			name: "update user success",
			args: args{
				ctx:         context.Background(),
				updateParam: updateParamSample,
				selectParam: selectParamSample,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectExec(queryUpdate).WillReturnResult(driver.RowsAffected(1))
				return sqlServer, err
			},
			wantErr: false,
		},
	}

	// Iterate the test in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqlServer, err := tt.prepSqlMock()
			if err != nil {
				t.Error(err)
			}
			defer sqlServer.Close()

			sqlClient := libsql.Init(libsql.Config{
				Driver: "sqlmock",
				Leader: libsql.ConnConfig{
					MockDB: sqlServer,
				},
				Follower: libsql.ConnConfig{
					MockDB: sqlServer,
				},
			}, logger, nil)

			// Initialize the Domain
			domain := Init(InitParam{
				Log:  logger,
				Db:   sqlClient,
				Json: mockJsonParser,
			})

			err = domain.Update(tt.args.ctx, tt.args.updateParam, tt.args.selectParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("domain.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
