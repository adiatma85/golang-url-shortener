package url

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	mock_url_dom "github.com/adiatma85/golang-url-shortener/src/business/domain/mock/url"
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
	mock_jwt_auth "github.com/adiatma85/own-go-sdk/tests/mock/jwtAuth"
	mock_log "github.com/adiatma85/own-go-sdk/tests/mock/log"
	mock_redis "github.com/adiatma85/own-go-sdk/tests/mock/redis"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type mockInterface struct {
	logger  *mock_log.MockInterface
	urlDom  *mock_url_dom.MockInterface
	jwtAuth *mock_jwt_auth.MockInterface
	redis   *mock_redis.MockInterface
}

func initMockTest(t *testing.T) (Interface, url, mockInterface) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mock_log.NewMockInterface(ctrl)
	logger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	mockUrlDom := mock_url_dom.NewMockInterface(ctrl)
	mockJwtAuth := mock_jwt_auth.NewMockInterface(ctrl)
	mockRedis := mock_redis.NewMockInterface(ctrl)

	ucInterface := Init(InitParam{
		Log:     logger,
		Url:     mockUrlDom,
		JwtAuth: mockJwtAuth,
		Redis:   mockRedis,
	})

	ucStruct := url{
		log:     logger,
		url:     mockUrlDom,
		jwtAuth: mockJwtAuth,
		redis:   mockRedis,
	}

	mockInterface := mockInterface{
		logger:  logger,
		urlDom:  mockUrlDom,
		jwtAuth: mockJwtAuth,
		redis:   mockRedis,
	}

	return ucInterface, ucStruct, mockInterface
}

func Test_url_Create(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx context.Context
		req entity.CreateUrlParam
	}

	// Mock in here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockCreateCategoryArgs := entity.CreateUrlParam{
		OriginalUrl: "This is url",
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockFinalResult := entity.Url{
		ID:          1,
		OriginalUrl: "This is url",
		CreatedBy:   null.StringFrom("10"),
		UpdatedBy:   null.StringFrom("10"),
		CreatedAt:   null.TimeFrom(mockTime),
		UpdatedAt:   null.TimeFrom(mockTime),
		Status:      null.Int64From(1),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.Url
		wantErr  bool
	}{
		{
			name: "failed to get user info",
			arg: args{
				ctx: context.Background(),
				req: mockCreateCategoryArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "failed to insert into database",
			arg: args{
				ctx: context.Background(),
				req: mockCreateCategoryArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Get(arg.ctx, gomock.Any()).Return(entity.Url{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "entry does not exist"))
				mock.urlDom.EXPECT().Create(arg.ctx, gomock.Any()).Return(entity.Url{}, assert.AnError)
			},
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "failed to insert into database",
			arg: args{
				ctx: context.Background(),
				req: mockCreateCategoryArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Get(arg.ctx, gomock.Any()).Return(entity.Url{}, errors.NewWithCode(codes.CodeSQLRecordDoesNotExist, "entry does not exist"))
				mock.urlDom.EXPECT().Create(arg.ctx, gomock.Any()).Return(mockFinalResult, nil)
			},
			want:    mockFinalResult,
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.Create(tt.arg.ctx, tt.arg.req)
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

func Test_url_Get(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UrlParam
	}

	// Mock in here
	mockUrlParam := entity.UrlParam{
		ID: null.Int64From(1),
		QueryOption: query.Option{
			IsActive: true,
		},
		UserId: null.Int64From(10),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockFinalResult := entity.Url{
		ID:          1,
		OriginalUrl: "Original url yang panjang dan lebar",
		Status:      null.Int64From(1),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.Url
		wantErr  bool
	}{
		{
			name: "failed to get user info",
			arg: args{
				ctx:    context.Background(),
				params: mockUrlParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "failed to get from category domain",
			arg: args{
				ctx:    context.Background(),
				params: mockUrlParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(entity.Url{}, assert.AnError)
			},
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:    context.Background(),
				params: mockUrlParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockFinalResult, nil)
			},
			want:    mockFinalResult,
			wantErr: false,
		},
	}

	// Iterate the tests in here
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

func Test_url_GetByShortenUrl(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UrlParam
	}

	// Mocks in here
	mockUrlParam := entity.UrlParam{
		ShortenUrl: "shortened-url",
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	assignKey := fmt.Sprintf(entity.UrlCountingRedisKey, mockUrlParam.ShortenUrl)

	mockFinalResult := entity.Url{
		OriginalUrl: "Long Original Url",
		ShortenUrl:  "shortened-url",
		Status:      null.Int64From(1),
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		want     entity.Url
		wantErr  bool
	}{
		{
			name: "failed to fetch from domain level",
			arg: args{
				ctx:    context.Background(),
				params: mockUrlParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(entity.Url{}, assert.AnError)
			},
			want:    entity.Url{},
			wantErr: true,
		},
		{
			name: "failed to insert to increment to the redis",
			arg: args{
				ctx:    context.Background(),
				params: mockUrlParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockFinalResult, nil)
				mock.redis.EXPECT().Increment(arg.ctx, assignKey).Return(assert.AnError)
			},
			want:    mockFinalResult,
			wantErr: false,
		},
		{
			name: "success",
			arg: args{
				ctx:    context.Background(),
				params: mockUrlParam,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockFinalResult, nil)
				mock.redis.EXPECT().Increment(arg.ctx, assignKey).Return(nil)
			},
			want:    mockFinalResult,
			wantErr: false,
		},
	}

	// Iterate the test cases in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			got, err := usecase.GetByShortenUrl(tt.arg.ctx, tt.arg.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetByShortenUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
				t.Errorf("usecase.GetByShortenUrl() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func Test_url_GetList(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UrlParam
	}

	// Mock in here
	mockParams := entity.UrlParam{
		ID: null.Int64From(1),
		PaginationParam: entity.PaginationParam{
			IncludePagination: true,
		},
		QueryOption: query.Option{
			IsActive: true,
		},
		UserId: null.Int64From(10),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID: 10,
		},
	}

	mockPagination := entity.Pagination{
		CurrentPage:     1,
		CurrentElements: 1,
		TotalPages:      1,
		TotalElements:   1,
	}

	mockResult := []entity.Url{
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
		want       []entity.Url
		pagination *entity.Pagination
		wantErr    bool
	}{
		{
			name: "failed to fetch from user info",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			want:       []entity.Url{},
			pagination: nil,
			wantErr:    true,
		},
		{
			name: "error fetching from domain",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().GetList(arg.ctx, mockParams).Return([]entity.Url{}, nil, assert.AnError)
			},
			want:       []entity.Url{},
			pagination: nil,
			wantErr:    true,
		},
		{
			name: "all success",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().GetList(arg.ctx, mockParams).Return(mockResult, &mockPagination, nil)
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

func Test_url_GetListAsAdmin(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx    context.Context
		params entity.UrlParam
	}

	// Mock in here
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
	}

	mockResult := []entity.Url{
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
		want       []entity.Url
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
				mock.urlDom.EXPECT().GetList(arg.ctx, mockParams).Return([]entity.Url{}, nil, assert.AnError)
			},
			want:       []entity.Url{},
			pagination: nil,
			wantErr:    true,
		},
		{
			name: "error fetching from domain",
			arg: args{
				ctx:    context.Background(),
				params: mockParams,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.urlDom.EXPECT().GetList(arg.ctx, mockParams).Return(mockResult, &mockPagination, nil)
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

func Test_category_Update(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		updateParam entity.UpdateUrlParam
		selectParam entity.UrlParam
	}

	// Mock in here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}
	mockUpdateArgs := entity.UpdateUrlParam{
		Visit: 1,
	}

	mockSelectArgs := entity.UrlParam{
		ID: null.Int64From(1),
		QueryOption: query.Option{
			IsActive: true,
		},
		UserId: null.Int64From(10),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID:     10,
			RoleID: entity.RoleIdUser,
		},
	}

	mockUpdateParam := entity.UpdateUrlParam{
		Visit:     1,
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
			name: "failed to get user info",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUpdateArgs,
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(jwtAuth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update to category domain",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUpdateArgs,
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Update(arg.ctx, mockUpdateParam, mockSelectArgs).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "success",
			arg: args{
				ctx:         context.Background(),
				updateParam: mockUpdateArgs,
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Update(arg.ctx, mockUpdateParam, mockSelectArgs).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the test in here
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

func Test_category_Delete(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx         context.Context
		selectParam entity.UrlParam
	}

	// Mock in here
	mockTime := time.Now()
	Now = func() time.Time {
		return mockTime
	}

	mockSelectArgs := entity.UrlParam{
		ID: null.Int64From(1),
		QueryOption: query.Option{
			IsActive: true,
		},
		UserId: null.Int64From(10),
	}

	mockUserInfo := jwtAuth.UserAuthInfo{
		User: jwtAuth.User{
			ID:     10,
			RoleID: entity.RoleIdUser,
		},
	}

	mockDeleteParam := entity.UpdateUrlParam{
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
			name: "failed to get user info",
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
			name: "failed to update to category domain",
			arg: args{
				ctx:         context.Background(),
				selectParam: mockSelectArgs,
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.jwtAuth.EXPECT().GetUserAuthInfo(arg.ctx).Return(mockUserInfo, nil)
				mock.urlDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectArgs).Return(assert.AnError)
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
				mock.urlDom.EXPECT().Update(arg.ctx, mockDeleteParam, mockSelectArgs).Return(nil)
			},
			wantErr: false,
		},
	}

	// Iterate the test in here
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

func Test_url_AssignCounterScheduler(t *testing.T) {
	usecase, _, mocks := initMockTest(t)

	// Type in here
	type args struct {
		ctx context.Context
	}

	// Mock in here
	scanKey := fmt.Sprintf(entity.UrlCountingRedisKey, "*")

	mockScanResult := []string{
		fmt.Sprintf(entity.UrlCountingRedisKey, "url-shortened"),
	}

	mockUrlParam := entity.UrlParam{
		ShortenUrl: "url-shortened",
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	mockUrlDomResult := entity.Url{
		Visit: 0,
	}

	mockUrlUpdateParam := entity.UpdateUrlParam{
		Visit: 1,
	}

	// Test cases in here
	tests := []struct {
		name     string
		arg      args
		mockFunc func(mock mockInterface, arg args)
		wantErr  bool
	}{
		{
			name: "failed even when scanning",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return([]string{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "zero slice of string returned when scanning",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return([]string{}, nil)
			},
			wantErr: false,
		},
		{
			name: "failed when fetching value from redis",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return(mockScanResult, nil)

				// First key
				mock.redis.EXPECT().Get(arg.ctx, mockScanResult[0]).Return("1", assert.AnError)

			},
			wantErr: false,
		},
		{
			name: "failed to str atoi for value is not number",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return(mockScanResult, nil)

				// First key
				mock.redis.EXPECT().Get(arg.ctx, mockScanResult[0]).Return("NaN", nil)

			},
			wantErr: false,
		},
		{
			name: "failed to get from url domain",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return(mockScanResult, nil)

				// First key
				mock.redis.EXPECT().Get(arg.ctx, mockScanResult[0]).Return("1", nil)
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockUrlDomResult, assert.AnError)

			},
			wantErr: false,
		},
		{
			name: "failed to update to url domain",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return(mockScanResult, nil)

				// First key
				mock.redis.EXPECT().Get(arg.ctx, mockScanResult[0]).Return("1", nil)
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockUrlDomResult, nil)
				mock.urlDom.EXPECT().Update(arg.ctx, mockUrlUpdateParam, mockUrlParam).Return(assert.AnError)

			},
			wantErr: true,
		},
		{
			name: "failed to decrement redis",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return(mockScanResult, nil)

				// First key
				mock.redis.EXPECT().Get(arg.ctx, mockScanResult[0]).Return("1", nil)
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockUrlDomResult, nil)
				mock.urlDom.EXPECT().Update(arg.ctx, mockUrlUpdateParam, mockUrlParam).Return(nil)
				mock.redis.EXPECT().DecrementBy(arg.ctx, mockScanResult[0], int64(1)).Return(assert.AnError)

			},
			wantErr: false,
		},
		{
			name: "success",
			arg: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockInterface, arg args) {
				mock.redis.EXPECT().Scan(arg.ctx, scanKey).Return(mockScanResult, nil)

				// First key
				mock.redis.EXPECT().Get(arg.ctx, mockScanResult[0]).Return("1", nil)
				mock.urlDom.EXPECT().Get(arg.ctx, mockUrlParam).Return(mockUrlDomResult, nil)
				mock.urlDom.EXPECT().Update(arg.ctx, mockUrlUpdateParam, mockUrlParam).Return(nil)
				mock.redis.EXPECT().DecrementBy(arg.ctx, mockScanResult[0], int64(1)).Return(nil)

			},
			wantErr: false,
		},
	}

	// Iterate the tests in here
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.arg)

			err := usecase.AssignCounterScheduler(tt.arg.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.AssignCounterScheduler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
