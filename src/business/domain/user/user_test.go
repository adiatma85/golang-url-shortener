package user

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	libsql "github.com/adiatma85/own-go-sdk/sql"
	mock_log "github.com/adiatma85/own-go-sdk/tests/mock/log"
	mock_json "github.com/adiatma85/own-go-sdk/tests/mock/parser"

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

	// Type in here
	type args struct {
		ctx             context.Context
		createUserParam entity.CreateUserParam
	}

	// Mock in here
	mockCreateUserParam := entity.CreateUserParam{
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
		mockFunc    func(args args)
		want        entity.User
		wantErr     bool
	}{
		{
			name: "cannot begin tx",
			args: args{
				ctx:             context.Background(),
				createUserParam: mockCreateUserParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, _, err := sqlmock.New()
				return sqlServer, err
			},
			mockFunc: func(args args) {},
			want:     entity.User{},
			wantErr:  true,
		},
		{
			name: "cannot exec user",
			args: args{
				ctx:             context.Background(),
				createUserParam: mockCreateUserParam,
			},
			prepSqlMock: func() (*sql.DB, error) {
				sqlServer, sqlMock, err := sqlmock.New()
				sqlMock.ExpectBegin()

				createUser := regexp.QuoteMeta(`INSERT INTO user (fk_role_id, email, username, password, display_name, created_by)
				VALUES (?, ?, ?, ?, ?, ?)`)
				sqlMock.ExpectExec(createUser).WillReturnError(errors.NewWithCode(codes.CodeSQL, "cannot create user"))
				sqlMock.ExpectRollback()

				return sqlServer, err
			},
			mockFunc: func(args args) {},
			want:     entity.User{},
			wantErr:  true,
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

			tt.mockFunc(tt.args)

			// Initialize the Domain
			u := Init(InitParam{
				Log:  logger,
				Db:   sqlClient,
				Json: mockJsonParser,
			})

			got, err := u.Create(tt.args.ctx, tt.args.createUserParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
