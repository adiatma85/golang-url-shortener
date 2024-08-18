package url

import (
	"context"
	"fmt"
	"time"

	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
	"github.com/adiatma85/own-go-sdk/redis"

	urlDom "github.com/adiatma85/golang-url-shortener/src/business/domain/url"
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	goNanoId "github.com/matoous/go-nanoid/v2"
)

type Interface interface {
	Create(ctx context.Context, insertParam entity.CreateUrlParam) (entity.Url, error)
	Get(ctx context.Context, urlParam entity.UrlParam) (entity.Url, error)
	GetList(ctx context.Context, urlParam entity.UrlParam) ([]entity.Url, *entity.Pagination, error)
	GetListAsAdmin(ctx context.Context, params entity.UrlParam) ([]entity.Url, *entity.Pagination, error)
	Update(ctx context.Context, updateParam entity.UpdateUrlParam, selectParam entity.UrlParam) error
	Delete(ctx context.Context, selectParam entity.UrlParam) error

	// Scheduler functions below
	AssignCounterScheduler(ctx context.Context) error
}

type InitParam struct {
	Log     log.Interface
	Url     urlDom.Interface
	JwtAuth jwtAuth.Interface
	Redis   redis.Interface
}

type url struct {
	log     log.Interface
	url     urlDom.Interface
	jwtAuth jwtAuth.Interface
	redis   redis.Interface
}

var Now = time.Now

func Init(params InitParam) Interface {
	u := &url{
		log:     params.Log,
		url:     params.Url,
		jwtAuth: params.JwtAuth,
		redis:   params.Redis,
	}

	return u
}

// Create with the injection of user id
func (u *url) Create(ctx context.Context, insertParam entity.CreateUrlParam) (entity.Url, error) {
	var (
		result         entity.Url
		duplicateError error
	)
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return result, err
	}

	insertParam.UserId = user.User.ID
	insertParam.CreatedBy = null.StringFrom(fmt.Sprintf("%v", user.User.ID))
	insertParam.UpdatedBy = null.StringFrom(fmt.Sprintf("%v", user.User.ID))

	// Generate shorten url
	for {
		shortenUrl, err := u.generateShortenUrl()
		if err != nil {
			return result, err
		}
		insertParam.ShortenUrl = shortenUrl
		duplicateError = u.validateShortenedUrl(ctx, insertParam)

		if duplicateError != nil {
			break
		}
	}

	result, err = u.url.Create(ctx, insertParam)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (u *url) validateShortenedUrl(ctx context.Context, urlInsertBody entity.CreateUrlParam) error {
	urlParam := entity.UrlParam{
		ShortenUrl: urlInsertBody.ShortenUrl,
		QueryOption: query.Option{
			IsActive: true,
		},
	}

	existedUrl, err := u.url.Get(ctx, urlParam)
	if err != nil && errors.GetCode(err) != codes.CodeSQLRecordDoesNotExist {
		return err
	}

	// Duplicate shorten url, but different original url
	if existedUrl.OriginalUrl != urlInsertBody.OriginalUrl {
		return errors.NewWithCode(codes.CodeConflict, "duplicate shorten url")
	}

	return nil
}

func (u *url) generateShortenUrl() (string, error) {
	shortenUrl, err := goNanoId.Generate(entity.Base62Chars, 5)
	if err != nil {
		return "", err
	}

	return shortenUrl, nil
}

func (u *url) Get(ctx context.Context, params entity.UrlParam) (entity.Url, error) {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.Url{}, err
	}

	params.QueryOption = query.Option{
		IsActive: true,
	}

	// Assign the params with user id
	params.UserId = null.Int64From(user.User.ID)

	return u.url.Get(ctx, params)
}

func (u *url) GetList(ctx context.Context, params entity.UrlParam) ([]entity.Url, *entity.Pagination, error) {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return []entity.Url{}, nil, err
	}

	params.IncludePagination = true
	params.QueryOption = query.Option{
		IsActive: true,
	}

	// Assign the params with user id
	params.UserId = null.Int64From(user.User.ID)

	return u.url.GetList(ctx, params)
}

func (u *url) GetListAsAdmin(ctx context.Context, params entity.UrlParam) ([]entity.Url, *entity.Pagination, error) {
	params.IncludePagination = true

	return u.url.GetList(ctx, params)
}

func (u *url) Update(ctx context.Context, updateParam entity.UpdateUrlParam, selectParam entity.UrlParam) error {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	// Assign the user id if they are not admin
	if user.User.RoleID != entity.RoleIdSuperAdmin {
		selectParam.UserId = null.Int64From(user.User.ID)
	}

	updateParam.UpdatedAt = null.TimeFrom(Now())
	updateParam.UpdatedBy = null.StringFrom(fmt.Sprintf("%v", user.User.ID))

	return u.url.Update(ctx, updateParam, selectParam)
}

func (u *url) Delete(ctx context.Context, selectParam entity.UrlParam) error {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	// Assign the user id if they are not admin
	if user.User.RoleID != entity.RoleIdSuperAdmin {
		selectParam.UserId = null.Int64From(user.User.ID)
	}

	deleteParam := entity.UpdateUrlParam{
		Status:    null.Int64From(-1),
		DeletedAt: null.TimeFrom(Now()),
		DeletedBy: null.StringFrom(fmt.Sprintf("%v", user.User.ID)),
	}

	return u.url.Update(ctx, deleteParam, selectParam)
}
