package user

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
	"golang.org/x/crypto/bcrypt"

	userDom "github.com/adiatma85/golang-url-shortener/src/business/domain/user"
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
)

type Interface interface {
	CreateWithoutAuthInfo(ctx context.Context, params entity.CreateUserParam) (entity.User, error)
	Create(ctx context.Context, req entity.CreateUserParam) (entity.User, error)
	Get(ctx context.Context, params entity.UserParam) (entity.User, error)
	GetListAsAdmin(ctx context.Context, params entity.UserParam) ([]entity.User, *entity.Pagination, error)
	Update(ctx context.Context, updateParam entity.UpdateUserParam, selectParam entity.UserParam) error
	Delete(ctx context.Context, selectParam entity.UserParam) error
	SignInWithPassword(ctx context.Context, req entity.UserLoginRequest) (entity.UserLoginResponse, error)
	GetSelfProfile(ctx context.Context) (entity.User, error)
	SelfDelete(ctx context.Context) error
	ChangePassword(ctx context.Context, changePasswordReq entity.ChangePasswordRequest) error
	UpdateUserProfile(ctx context.Context, updateParam entity.UpdateUserParam) error

	// Improvement kedepannya
	// CheckPassword(ctx context.Context, params entity.UserCheckPasswordParam, userParam entity.UserParam) (entity.HTTPMessage, error)
	// Activate(ctx context.Context, selectParam entity.UserParam) error
	// RefreshToken(ctx context.Context, param entity.UserRefreshTokenParam) (entity.RefreshTokenResponse, error)
}

type InitParam struct {
	Log     log.Interface
	User    userDom.Interface
	JwtAuth jwtAuth.Interface
}

type user struct {
	log     log.Interface
	user    userDom.Interface
	jwtAuth jwtAuth.Interface
}

var Now = time.Now

func Init(param InitParam) Interface {
	u := &user{
		log:     param.Log,
		user:    param.User,
		jwtAuth: param.JwtAuth,
	}

	return u
}

func (u *user) Create(ctx context.Context, req entity.CreateUserParam) (entity.User, error) {
	var result entity.User

	result, err := u.validateUser(ctx, req)
	if err != nil {
		return result, err
	}

	req.CreatedBy = null.StringFrom(fmt.Sprintf("%v", entity.SystemUser))
	req.UpdatedBy = null.StringFrom(fmt.Sprintf("%v", entity.SystemUser))

	return u.user.Create(ctx, req)
}

func (u *user) CreateWithoutAuthInfo(ctx context.Context, req entity.CreateUserParam) (entity.User, error) {
	var result entity.User
	req.ConfirmPassword = req.Password

	result, err := u.validateUser(ctx, req)
	if err != nil {
		return result, err
	}

	// Hash the password in here
	req.Password, err = u.getHashPassowrd(req.Password)
	if err != nil {
		return result, err
	}

	return u.user.Create(ctx, req)
}

func (u *user) validateUser(ctx context.Context, req entity.CreateUserParam) (entity.User, error) {
	var result entity.User

	if req.Password != req.ConfirmPassword {
		return result, errors.NewWithCode(codes.CodePasswordDoesNotMatch, "password does not match")
	}

	user, err := u.user.Get(ctx, entity.UserParam{
		Email: null.StringFrom(req.Email),
	})
	if err != nil && errors.GetCode(err) != codes.CodeSQLRecordDoesNotExist {
		return result, err
	}

	if user != result {
		return result, errors.NewWithCode(codes.CodeConflict, "email is exists")
	}

	return result, nil
}

func (u *user) Get(ctx context.Context, params entity.UserParam) (entity.User, error) {
	return u.user.Get(ctx, params)
}

func (u *user) GetListAsAdmin(ctx context.Context, params entity.UserParam) ([]entity.User, *entity.Pagination, error) {
	params.IncludePagination = true
	users, pg, err := u.user.GetList(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	return users, pg, nil
}

func (u *user) Update(ctx context.Context, updateParam entity.UpdateUserParam, selectParam entity.UserParam) error {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	updateParam.UpdatedAt = null.TimeFrom(Now())
	updateParam.UpdatedBy = null.StringFrom(fmt.Sprintf("%v", user.User.ID))

	return u.user.Update(ctx, updateParam, selectParam)
}

func (u *user) Delete(ctx context.Context, selectParam entity.UserParam) error {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	deleteParam := entity.UpdateUserParam{
		Status:    null.Int64From(-1),
		DeletedAt: null.TimeFrom(Now()),
		DeletedBy: null.StringFrom(fmt.Sprintf("%v", user.User.ID)),
	}

	return u.user.Update(ctx, deleteParam, selectParam)
}

func (u *user) getHashPassowrd(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Return true if the password is match
func (u *user) checkHashPassword(ctx context.Context, hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	if err != nil {
		u.log.Error(ctx, err)
		return false
	}

	return true
}

func (u *user) SignInWithPassword(ctx context.Context, req entity.UserLoginRequest) (entity.UserLoginResponse, error) {
	// validate body request
	if err := u.validateUserLoginRequest(req); err != nil {
		return entity.UserLoginResponse{}, err
	}

	// validate user is exist on db and status is active
	user, err := u.user.Get(ctx, entity.UserParam{
		Email: null.StringFrom(req.Email),
		QueryOption: query.Option{
			IsActive: true,
		},
	})
	if err != nil {
		if errors.GetCode(err) == codes.CodeSQLRecordDoesNotExist {
			return entity.UserLoginResponse{}, errors.NewWithCode(codes.CodeNotFound, "email not found")
		}

		u.log.Error(ctx, err)
		return entity.UserLoginResponse{}, err
	}

	// Validate the password in here
	if !u.checkHashPassword(ctx, user.Password, req.Password) {
		return entity.UserLoginResponse{}, errors.NewWithCode(codes.CodeUnauthorized, "credential does not match")
	}

	// Create the JWT token in here
	accessToken, err := u.jwtAuth.CreateAccessToken(user.ConvertToAuthUser())
	if err != nil {
		return entity.UserLoginResponse{}, err
	}

	result := entity.UserLoginResponse{
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AccessToken: accessToken,
	}

	return result, nil
}

func (u *user) validateUserLoginRequest(req entity.UserLoginRequest) error {
	if req.Email == "" {
		return errors.NewWithCode(codes.CodeBadRequest, "email is required")
	}

	if req.Password == "" {
		return errors.NewWithCode(codes.CodeBadRequest, "password is required")
	}

	return nil
}

func (u *user) GetSelfProfile(ctx context.Context) (entity.User, error) {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.User{}, err
	}

	userParam := entity.UserParam{
		ID: null.Int64From(user.User.ID),
	}

	return u.user.Get(ctx, userParam)
}

func (u *user) SelfDelete(ctx context.Context) error {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	selectParam := entity.UserParam{
		ID: null.Int64From(user.User.ID),
	}

	deleteParam := entity.UpdateUserParam{
		Status:    null.Int64From(-1),
		DeletedAt: null.TimeFrom(Now()),
		DeletedBy: null.StringFrom(fmt.Sprintf("%v", user.User.ID)),
	}

	return u.user.Update(ctx, deleteParam, selectParam)
}

func (u *user) ChangePassword(ctx context.Context, changePasswordReq entity.ChangePasswordRequest) error {
	// Check first if the password and confirm password is match
	if changePasswordReq.Password != changePasswordReq.ConfirmPassword {
		return errors.NewWithCode(codes.CodeBadRequest, "new password and confirm password does not match")
	}

	// Check if the old password is match
	userAuth, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	userDn, err := u.user.Get(ctx, entity.UserParam{
		ID: null.Int64From(userAuth.User.ID),
		QueryOption: query.Option{
			IsActive: true,
		},
	})

	if err != nil {
		return err
	}

	if !u.checkHashPassword(ctx, userDn.Password, changePasswordReq.OldPassword) {
		return errors.NewWithCode(codes.CodeUnauthorized, "credential does not match")
	}

	// Update the entry
	hashedPass, err := u.getHashPassowrd(changePasswordReq.Password)
	if err != nil {
		return err
	}

	updateParam := entity.UpdateUserParam{
		Password: hashedPass,
	}

	selectParam := entity.UserParam{
		ID: null.Int64From(userDn.ID),
	}

	return u.user.Update(ctx, updateParam, selectParam)
}

func (u *user) UpdateUserProfile(ctx context.Context, updateParam entity.UpdateUserParam) error {
	user, err := u.jwtAuth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	userParam := entity.UserParam{
		ID: null.Int64From(user.User.ID),
	}

	return u.user.Update(ctx, updateParam, userParam)
}
