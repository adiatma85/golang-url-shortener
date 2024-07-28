package url

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

func Test_url_Create(t *testing.T) {
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
		createParam entity.CreateUrlParam
	}

	// Mock in here
	mockCreateParam := entity.CreateUrlParam{
		ShortenUrl:  "Shortened Url",
		OriginalUrl: "Original Url",
	}

	query := regexp.QuoteMeta(`INSERT INTO url (original_url, shorten_url, fk_user_id, created_by, updated_by)
	VALUES (?, ?, ?, ?, ?)`)
	queryGet := regexp.QuoteMeta(readUrl)

	// Test cases in here
	tests := []struct {
		name        string
		args        args
		prepSqlMock func() (*sql.DB, error)
		want        entity.Url
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
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "cannot exec url",
			args: args{
				ctx:         context.Background(),
				createParam: mockCreateParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()

				sqlMock.ExpectExec(query).WillReturnError(errors.NewWithCode(codes.CodeSQL, "cannot create url"))
				sqlMock.ExpectRollback()

				return sqlServer, err
			},
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "url no new row",
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
			want:    entity.Url{},
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
			want: entity.Url{
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
					"original_url",
					"shorten_url",
				})
				row.AddRow("1", "Original Url", "Shortened Url")
				sqlMock.ExpectQuery(queryGet).WithArgs(1).WillReturnRows(row)

				sqlMock.ExpectRollback()
				return sqlServer, err
			},
			want: entity.Url{
				ID:          1,
				ShortenUrl:  "Shortened Url",
				OriginalUrl: "Original Url",
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

func Test_url_Get(t *testing.T) {
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
		params entity.UrlParam
	}

	// Mock in here
	now := time.Now()
	query := regexp.QuoteMeta(readUrl)

	mockParam := entity.UrlParam{
		ID: null.Int64From(1),
	}

	sampleResult := entity.Url{
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
		want        entity.Url
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
			want:    entity.Url{},
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
			want:    entity.Url{},
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
			want:    sampleResult,
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

func Test_url_GetList(t *testing.T) {
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
		params entity.UrlParam
	}

	// Mock in here
	now := time.Now()
	queryExt := " WHERE 1=1 AND id=? LIMIT 0, 10;"
	query := regexp.QuoteMeta(readUrl + queryExt)
	queryCountExt := " WHERE 1=1 AND id=?;"
	queryCount := regexp.QuoteMeta(readUrlCount + queryCountExt)

	mockParams := entity.UrlParam{
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
		want        []entity.Url
		want1       *entity.Pagination
		wantErr     bool
	}{
		{
			name: "error when query-ing",
			args: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectQuery(query).WillReturnError(fmt.Errorf("failed to get list of url"))
				return sqlServer, err
			},
			want:    []entity.Url{},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "error when struct scan",
			args: args{
				ctx:    context.Background(),
				params: mockParams,
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
			want:    []entity.Url{},
			want1:   nil,
			wantErr: true,
		},
		{
			name: "all good",
			args: args{
				ctx:    context.Background(),
				params: mockParams,
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
			want: []entity.Url{
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

func Test_url_Update(t *testing.T) {
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
		updateParam entity.UpdateUrlParam
		selectParam entity.UrlParam
	}

	// Mock in here
	queryUpdate := regexp.QuoteMeta("UPDATE url SET visit=?, updated_by=? WHERE 1=1 AND status=1 AND id=?")

	selectParamSample := entity.UrlParam{
		ID: null.Int64From(1),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	updateParamSample := entity.UpdateUrlParam{
		Visit:     1,
		UpdatedBy: null.StringFrom("1"),
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
			name: "no url updated",
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
			name: "update url 0 rows affected",
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
			name: "update url success",
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
