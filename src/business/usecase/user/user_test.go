package user

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	mock_user_dom "github.com/adiatma85/golang-url-shortener/src/business/domain/mock/user"
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
	mock_jwt_auth "github.com/adiatma85/own-go-sdk/tests/mock/jwtAuth"
	mock_log "github.com/adiatma85/own-go-sdk/tests/mock/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

type mockInterface struct {
	logger  *mock_log.MockInterface
	userDom *mock_user_dom.MockInterface
	jwtAuth *mock_jwt_auth.MockInterface
}

func initMockTest(t *testing.T) (Interface, user, mockInterface) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockUserDom := mock_user_dom.NewMockInterface(ctrl)
	mockJwtAuth := mock_jwt_auth.NewMockInterface(ctrl)

	ucInterface := Init(InitParam{
		Log:     logger,
		User:    mockUserDom,
		JwtAuth: mockJwtAuth,
	})

	ucStruct := user{
		log:     logger,
		user:    mockUserDom,
		jwtAuth: mockJwtAuth,
	}

	mockInterface := mockInterface{
		logger:  logger,
		userDom: mockUserDom,
		jwtAuth: mockJwtAuth,
	}

	return ucInterface, ucStruct, mockInterface
}

func Test_user_Create(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		insertParam entity.CreateUserParam
	}

	// Mock in here
	mockInsertParam := entity.CreateUserParam{
		Email:           "random@email.com",
		Password:        "samePassword",
		ConfirmPassword: "samePassword",
		CreatedBy:       null.StringFrom(fmt.Sprintf("%v", 10)),
		UpdatedBy:       null.StringFrom(fmt.Sprintf("%v", 10)),
	}

	mockUserFetchParam := entity.UserParam{
		Email: null.StringFrom("random@email.com"),
	}

	mockJwtResult := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockResult := entity.User{
		Email:     "random@email.com",
		CreatedBy: null.StringFrom(fmt.Sprintf("%v", 10)),
		UpdatedBy: null.StringFrom(fmt.Sprintf("%v", 10)),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to fetch from user domain",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "failed to fetch from user domain",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "email already exists",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{Email: "random@email.com"}, nil)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "failed to get user information from context",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, nil)
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtResult, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "failed to create user in domain layer",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, nil)
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtResult, nil)
				mock.userDom.EXPECT().Create(arg.ctx, gomock.Any()).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, nil)
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtResult, nil)
				mock.userDom.EXPECT().Create(arg.ctx, gomock.Any()).Return(mockResult, nil)
			},
			want:    mockResult,
			wantErr: false,
		},
	}

	// Iterate the test in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.Create(tt.arg.ctx, tt.arg.insertParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.Create() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_user_CreateWithoutAuthInfo(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		insertParam entity.CreateUserParam
	}

	// Mocks in here
	mockInsertParam := entity.CreateUserParam{
		Email:           "random@email.com",
		Password:        "samePassword",
		ConfirmPassword: "samePassword",
		CreatedBy:       null.StringFrom(fmt.Sprintf("%v", 10)),
		UpdatedBy:       null.StringFrom(fmt.Sprintf("%v", 10)),
	}

	mockUserFetchParam := entity.UserParam{
		Email: null.StringFrom("random@email.com"),
	}

	mockResult := entity.User{
		Email:     "random@email.com",
		CreatedBy: null.StringFrom(fmt.Sprintf("%v", 10)),
		UpdatedBy: null.StringFrom(fmt.Sprintf("%v", 10)),
	}

	// Test in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "password does not match",
			arg: args{
				ctx: context.Background(),
				insertParam: entity.CreateUserParam{
					Password:        "password 1",
					ConfirmPassword: "password 2",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {},
			want:     entity.User{},
			wantErr:  true,
		},
		{
			name: "failed to fetch from user domain",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "email already exists",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{Email: "random@email.com"}, nil)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "failed to create user in domain layer",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, nil)
				// Using Go Mock any because there are bcrypt hash password in here
				mock.userDom.EXPECT().Create(arg.ctx, gomock.Any()).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				insertParam: mockInsertParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserFetchParam).Return(entity.User{}, nil)
				// Using Go Mock any because there are bcrypt hash password in here
				mock.userDom.EXPECT().Create(arg.ctx, gomock.Any()).Return(mockResult, nil)
			},
			want:    mockResult,
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.CreateWithoutAuthInfo(tt.arg.ctx, tt.arg.insertParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.CreateWithoutAuthInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.CreateWithoutAuthInfo() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_user_Get(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UserParam
	}

	// Mock in here
	mockParams := entity.UserParam{
		ID: null.Int64From(1),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	mockResult := entity.User{
		ID:     1,
		Status: null.Int64From(1),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "error fetching from domain",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockParams).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockParams).Return(mockResult, nil)
			},
			want:    mockResult,
			wantErr: false,
		},
	}

	// Iterate the tests cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.Get(tt.arg.ctx, tt.arg.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.Get() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_user_GetList(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UserParam
	}

	// Mock in here
	mockParams := entity.UserParam{
		ID: null.Int64From(1),
		PaginationParam: entity.PaginationParam{
			IncludePagination: true,
		},
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	mockPagination := entity.Pagination{
		CurrentPage:     1,
		CurrentElements: 1,
		TotalPages:      1,
		TotalElements:   1,
	}

	mockResult := []entity.User{
		{
			ID:     1,
			Status: null.Int64From(1),
		},
	}

	// Test cases in here
	tests := []struct {
		name       string
		arg        args
		mockFunc   func(mock mockInterface, arg args)
		want       []entity.User
		pagination *entity.Pagination
		wantErr    bool
	}{
		{
			name: "error fetching from domain",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().GetList(arg.ctx, mockParams).Return([]entity.User{}, nil, assert.AnError)
			},
			want:       nil,
			pagination: nil,
			wantErr:    true,
		},
		{
			name: "success",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().GetList(arg.ctx, mockParams).Return(mockResult, &mockPagination, nil)
			},
			want:       mockResult,
			pagination: &mockPagination,
			wantErr:    false,
		},
	}

	// Iterate the tests cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, pagination, err := usecase.GetList(tt.arg.ctx, tt.arg.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.GetList() got = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(pagination, tt.pagination) {
				assert.Equal(t, tt.pagination, pagination)
				t.Errorf("usecase.GetList() got = %v, want %v", pagination, tt.pagination)
				return
			}
		})
	}
}

func Test_user_GetAsAdmin(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UserParam
	}

	// Mock in here
	mockParams := entity.UserParam{
		ID: null.Int64From(1),
	}

	mockResult := entity.User{
		ID:     1,
		Status: null.Int64From(1),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "error fetching from domain",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockParams).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockParams).Return(mockResult, nil)
			},
			want:    mockResult,
			wantErr: false,
		},
	}

	// Iterate the tests cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.GetAsAdmin(tt.arg.ctx, tt.arg.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetAsAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.GetAsAdmin() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_user_GetListAsAdmin(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UserParam
	}

	// Mock in here
	mockParams := entity.UserParam{
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
	}

	mockResult := []entity.User{
		{
			ID:     1,
			Status: null.Int64From(1),
		},
	}

	// Test cases in here
	tests := []struct {
		name       string
		arg        args
		mockFunc   func(mock mockInterface, arg args)
		want       []entity.User
		pagination *entity.Pagination
		wantErr    bool
	}{
		{
			name: "error fetching from domain",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().GetList(arg.ctx, mockParams).Return([]entity.User{}, nil, assert.AnError)
			},
			want:       nil,
			pagination: nil,
			wantErr:    true,
		},
		{
			name: "success",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().GetList(arg.ctx, mockParams).Return(mockResult, &mockPagination, nil)
			},
			want:       mockResult,
			pagination: &mockPagination,
			wantErr:    false,
		},
	}

	// Iterate the tests cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, pagination, err := usecase.GetListAsAdmin(tt.arg.ctx, tt.arg.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetListAsAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.GetListAsAdmin() got = %v, want %v", got, tt.want)
				return
			}

			if !reflect.DeepEqual(pagination, tt.pagination) {
				assert.Equal(t, tt.pagination, pagination)
				t.Errorf("usecase.GetListAsAdmin() got = %v, want %v", pagination, tt.pagination)
				return
			}
		})
	}
}

func Test_user_Update(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		updateParam entity.UpdateUserParam
		selectParam entity.UserParam
	}

	// Mock in here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockUpdateArg := entity.UpdateUserParam{
		DisplayName: "Ini nama yang diganti hehehe",
	}

	mockSelectArgs := entity.UserParam{
		ID: null.Int64From(1),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockUpdateParam := entity.UpdateUserParam{
		DisplayName: "Ini nama yang diganti hehehe",
		UpdatedAt:   null.TimeFrom(mockTime),
		UpdatedBy:   null.StringFrom("10"),
	}

	// Tests cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "failed to fetch user info",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUpdateArg,
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update in domain layer",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUpdateArg,
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockUpdateParam, mockSelectArgs).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUpdateArg,
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockUpdateParam, mockSelectArgs).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the tests cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.Update(tt.arg.ctx, tt.arg.updateParam, tt.arg.selectParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_Delete(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		selectParam entity.UserParam
	}

	// Mock In here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockSelectArgs := entity.UserParam{
		ID: null.Int64From(1),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockDeleteParam := entity.UpdateUserParam{
		Status:    null.Int64From(-1),
		DeletedAt: null.TimeFrom(mockTime),
		DeletedBy: null.StringFrom("10"),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "failed to fetch user info",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update in domain layer",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectArgs).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectArgs).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the test cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.Delete(tt.arg.ctx, tt.arg.selectParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_Activate(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		selectParam entity.UserParam
	}

	// Mock In here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockSelectArgs := entity.UserParam{
		ID: null.Int64From(1),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockDeleteParam := entity.UpdateUserParam{
		Status:    null.Int64From(1),
		UpdatedAt: null.TimeFrom(mockTime),
		UpdatedBy: null.StringFrom("10"),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "failed to fetch user info",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update in domain layer",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectArgs).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectArgs).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the test cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.Activate(tt.arg.ctx, tt.arg.selectParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Activate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_SignWithPassword(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx context.Context
		req entity.UserLoginRequest
	}

	// Mock in here
	mockLoginReq := entity.UserLoginRequest{
		Email:    "random@email.com",
		Password: "password",
	}

	mockUserParam := entity.UserParam{
		Email: null.StringFrom("random@email.com"),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(mockLoginReq.Password), bcrypt.DefaultCost)

	mockUserResult := entity.User{
		Email:       "random@email.com",
		DisplayName: "random name",
		Password:    string(hashedPassword),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.UserLoginResponse
		wantErr  bool
	}{
		{
			name: "failed because email is zero string",
			arg: args{
				ctx: context.Background(),
				req: entity.UserLoginRequest{},
			},
			mockFunc: func(mock mockInterface, arg args) {},
			want:     entity.UserLoginResponse{},
			wantErr:  true,
		},
		{
			name: "failed because password is zero string",
			arg: args{
				ctx: context.Background(),
				req: entity.UserLoginRequest{
					Email: "random@email.com",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {},
			want:     entity.UserLoginResponse{},
			wantErr:  true,
		},
		{
			name: "failed to fetch user from domain layer",
			arg: args{
				ctx: context.Background(),
				req: mockLoginReq,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(entity.User{}, assert.AnError)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "failed to fetch user from domain layer (error code does not exists)",
			arg: args{
				ctx: context.Background(),
				req: mockLoginReq,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(entity.User{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "row not found"))
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "failed to check hash passowrd",
			arg: args{
				ctx: context.Background(),
				req: mockLoginReq,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(entity.User{Password: "randomPassword"}, nil)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "failed to create access token",
			arg: args{
				ctx: context.Background(),
				req: mockLoginReq,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(mockUserResult, nil)
				mock.jwtAuth.EXPECT().CreateAccessToken(mockUserResult.ConvertToAuthUser()).Return("", assert.AnError)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "failed to create refresh token",
			arg: args{
				ctx: context.Background(),
				req: mockLoginReq,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(mockUserResult, nil)
				mock.jwtAuth.EXPECT().CreateAccessToken(mockUserResult.ConvertToAuthUser()).Return("ini access token yang panjang", nil)
				mock.jwtAuth.EXPECT().CreateRefreshToken(mockUserResult.ConvertToAuthUser()).Return("", assert.AnError)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx: context.Background(),
				req: mockLoginReq,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(mockUserResult, nil)
				mock.jwtAuth.EXPECT().CreateAccessToken(mockUserResult.ConvertToAuthUser()).Return("ini access token yang panjang", nil)
				mock.jwtAuth.EXPECT().CreateRefreshToken(mockUserResult.ConvertToAuthUser()).Return("ini refresh token yang panjang", nil)
			},
			want: entity.UserLoginResponse{
				Email:        "random@email.com",
				DisplayName:  "random name",
				AccessToken:  "ini access token yang panjang",
				RefreshToken: "ini refresh token yang panjang",
			},
			wantErr: false,
		},
	}

	// Iterate the tests cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.SignInWithPassword(tt.arg.ctx, tt.arg.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SignInWithPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.SignInWithPassword() got = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func Test_user_GetSelfProfile(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx context.Context
	}

	// Mock in here
	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockUserParam := entity.UserParam{
		ID: null.Int64From(10),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to fetch user info",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "failed to fetch user from user domain",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(entity.User{ID: 10}, nil)
			},
			want: entity.User{
				ID: 10,
			},
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.GetSelfProfile(tt.arg.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetSelfProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.GetSelfProfile() got = %v, want %v", got, tt.want)
				return
			}

		})
	}
}

func Test_user_SelfDelete(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx context.Context
	}

	// Mock in here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockDeleteParam := entity.UpdateUserParam{
		Status:    null.Int64From(-1),
		DeletedAt: null.TimeFrom(Now()),
		DeletedBy: null.StringFrom(fmt.Sprintf("%v", 10)),
	}

	mockSelectParam := entity.UserParam{
		ID: null.Int64From(10),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get user info",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectParam).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the test in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.SelfDelete(tt.arg.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SelfDelete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_ChangePassword(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx               context.Context
		changePasswordReq entity.ChangePasswordRequest
	}

	// Mock In here
	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockUserParam := entity.UserParam{
		ID: null.Int64From(10),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password_lama"), bcrypt.DefaultCost)

	mockUser := entity.User{
		Password: string(hashedPassword),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "change req password does not match with confirm password",
			arg: args{
				ctx: context.Background(),
				changePasswordReq: entity.ChangePasswordRequest{
					Password:        "passwordSama",
					ConfirmPassword: "passwordTidakSama",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {},
			wantErr:  true,
		},
		{
			name: "failed to get user info",
			arg: args{
				ctx: context.Background(),
				changePasswordReq: entity.ChangePasswordRequest{
					Password:        "password_sama",
					ConfirmPassword: "password_sama",
					OldPassword:     "password_lama",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to get user from user domain",
			arg: args{
				ctx: context.Background(),
				changePasswordReq: entity.ChangePasswordRequest{
					Password:        "password_sama",
					ConfirmPassword: "password_sama",
					OldPassword:     "password_lama",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(entity.User{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "old password is wrong",
			arg: args{
				ctx: context.Background(),
				changePasswordReq: entity.ChangePasswordRequest{
					Password:        "password_sama",
					ConfirmPassword: "password_sama",
					OldPassword:     "password_lama_salah",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(mockUser, nil)
			},
			wantErr: true,
		},
		{
			name: "failed to update user domain",
			arg: args{
				ctx: context.Background(),
				changePasswordReq: entity.ChangePasswordRequest{
					Password:        "password_sama",
					ConfirmPassword: "password_sama",
					OldPassword:     "password_lama",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(mockUser, nil)
				mock.userDom.EXPECT().Update(arg.ctx, gomock.Any(), mockUserParam).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx: context.Background(),
				changePasswordReq: entity.ChangePasswordRequest{
					Password:        "password_sama",
					ConfirmPassword: "password_sama",
					OldPassword:     "password_lama",
				},
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.userDom.EXPECT().Get(arg.ctx, mockUserParam).Return(mockUser, nil)
				mock.userDom.EXPECT().Update(arg.ctx, gomock.Any(), mockUserParam).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.ChangePassword(tt.arg.ctx, tt.arg.changePasswordReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.ChangePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_UpdateUserSelfProfile(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		updateParam entity.UpdateUserParam
	}

	// Mock in here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockUserUpdateArgs := entity.UpdateUserParam{
		DisplayName: "Nama editan yang panjang dan lebar",
	}

	mockJwtInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockUserUpdateParam := entity.UpdateUserParam{
		DisplayName: "Nama editan yang panjang dan lebar",
		UpdatedAt:   null.TimeFrom(mockTime),
		UpdatedBy:   null.StringFrom("10"),
	}

	mockUserParam := entity.UserParam{
		ID: null.Int64From(10),
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get user info",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUserUpdateArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update user in user domain",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUserUpdateArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockUserUpdateParam, mockUserParam).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUserUpdateArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtInfo, nil)
				mock.userDom.EXPECT().Update(arg.ctx, mockUserUpdateParam, mockUserParam).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.UpdateUserSelfProfile(tt.arg.ctx, tt.arg.updateParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.UpdateUserSelfProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_user_RefreshToken(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx context.Context
	}

	// Mock in here
	mockJwtInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.UserLoginResponse
		wantErr  bool
	}{
		{
			name: "failed to get user info",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "failed to create access token",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtInfo, nil)
				mock.jwtAuth.EXPECT().CreateAccessToken(mockJwtInfo.User).Return("", assert.AnError)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "failed to create refresh token",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtInfo, nil)
				mock.jwtAuth.EXPECT().CreateAccessToken(mockJwtInfo.User).Return("access token", nil)
				mock.jwtAuth.EXPECT().CreateRefreshToken(mockJwtInfo.User).Return("", assert.AnError)
			},
			want:    entity.UserLoginResponse{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockJwtInfo, nil)
				mock.jwtAuth.EXPECT().CreateAccessToken(mockJwtInfo.User).Return("access token", nil)
				mock.jwtAuth.EXPECT().CreateRefreshToken(mockJwtInfo.User).Return("refresh token", nil)
			},
			want: entity.UserLoginResponse{
				AccessToken:  "access token",
				RefreshToken: "refresh token",
			},
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.RefreshToken(tt.arg.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.RefreshToken() got = %v, want %v", got, tt.want)
				return
			}

		})
	}
}
