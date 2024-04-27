package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/appcontext"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/header"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// timeout middleware wraps the request context with a timeout
func (r *rest) SetTimeout(ctx *gin.Context) {
	// wrap the request context with a timeout
	c, cancel := context.WithTimeout(ctx.Request.Context(), r.conf.Timeout)

	// cancel to clear resources after finished
	defer cancel()

	c = appcontext.SetRequestStartTime(c, time.Now())

	// replace request with context wrapped request
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()

}

func (r *rest) httpRespSuccess(ctx *gin.Context, code codes.Code, data interface{}, p *entity.Pagination) {
	successApp := codes.Compile(code, appcontext.GetAcceptLanguage(ctx))
	c := ctx.Request.Context()
	meta := entity.Meta{
		Path:       r.conf.Meta.Host + ctx.Request.URL.String(),
		StatusCode: successApp.StatusCode,
		Status:     http.StatusText(successApp.StatusCode),
		Message:    fmt.Sprintf("%s %s [%d] %s", ctx.Request.Method, ctx.Request.URL.RequestURI(), successApp.StatusCode, http.StatusText(successApp.StatusCode)),
		Timestamp:  time.Now().Format(time.RFC3339),
		RequestID:  appcontext.GetRequestId(c),
	}

	resp := &entity.HTTPResp{
		Message: entity.HTTPMessage{
			Title: successApp.Title,
			Body:  successApp.Body,
		},
		Meta:       meta,
		Data:       data,
		Pagination: p,
	}

	reqstart := appcontext.GetRequestStartTime(c)
	if !time.Time.IsZero(reqstart) {
		resp.Meta.TimeElapsed = fmt.Sprintf("%dms", int64(time.Since(reqstart)/time.Millisecond))
	}

	raw, err := r.json.Marshal(&resp)
	if err != nil {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeInternalServerError, err.Error()))
		return
	}

	c = appcontext.SetAppResponseCode(c, code)
	c = appcontext.SetResponseHttpCode(c, successApp.StatusCode)
	ctx.Request = ctx.Request.WithContext(c)

	ctx.Header(header.KeyRequestID, appcontext.GetRequestId(c))
	ctx.Data(successApp.StatusCode, header.ContentTypeJSON, raw)
}

func (r *rest) httpRespError(ctx *gin.Context, err error) {
	c := ctx.Request.Context()

	if errors.Is(c.Err(), context.DeadlineExceeded) {
		err = errors.NewWithCode(codes.CodeContextDeadlineExceeded, "Context Deadline Exceeded")
	}

	httpStatus, displayError := errors.Compile(err, appcontext.GetAcceptLanguage(c))
	statusStr := http.StatusText(httpStatus)

	errResp := &entity.HTTPResp{
		Message: entity.HTTPMessage{
			Title: displayError.Title,
			Body:  displayError.Body,
		},
		Meta: entity.Meta{
			Path:       r.conf.Meta.Host + ctx.Request.URL.String(),
			StatusCode: httpStatus,
			Status:     statusStr,
			Message:    fmt.Sprintf("%s %s [%d] %s", ctx.Request.Method, ctx.Request.URL.RequestURI(), httpStatus, statusStr),
			Error: &entity.MetaError{
				Code:    int(displayError.Code),
				Message: err.Error(),
			},
			Timestamp: time.Now().Format(time.RFC3339),
			RequestID: appcontext.GetRequestId(c),
		},
	}

	r.log.Error(c, err)

	c = appcontext.SetAppResponseCode(c, displayError.Code)
	c = appcontext.SetAppErrorMessage(c, fmt.Sprintf("%s - %s", displayError.Title, displayError.Body))
	c = appcontext.SetResponseHttpCode(c, httpStatus)
	ctx.Request = ctx.Request.WithContext(c)

	ctx.Header(header.KeyRequestID, appcontext.GetRequestId(c))
	ctx.AbortWithStatusJSON(httpStatus, errResp)
}

// Bind request body to struct using tag 'json'
func (r *rest) Bind(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindWith(obj, binding.Default(ctx.Request.Method, ctx.ContentType()))
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// Bind all query params to struct using tag 'form'
func (r *rest) BindQuery(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindWith(obj, binding.Query)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// Bind uri params to struct using tag 'uri'
func (r *rest) BindUri(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindUri(obj)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// Bind all params (query and uri params) to struct using tag 'uri' and 'form'
func (r *rest) BindParams(ctx *gin.Context, obj interface{}) error {
	err := r.BindQuery(ctx, obj)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	err = r.BindUri(ctx, obj)
	if err != nil {
		return errors.NewWithCode(codes.CodeBadRequest, err.Error())
	}

	return nil
}

// @Summary Health Check
// @Description This endpoint will hit the server
// @Tags Server
// @Produce json
// @Success 200 string example="PONG!"
// @Router /ping [GET]
func (r *rest) Ping(ctx *gin.Context) {
	resp := entity.Ping{
		Status:  "OK",
		Version: r.conf.Meta.Version,
	}
	r.httpRespSuccess(ctx, codes.CodeSuccess, resp, nil)
}

func (r *rest) DummyLogin(ctx *gin.Context) {
	// conf := redactFirebaseAccountKey(r.configreader.AllSettings())

	ctx.HTML(http.StatusOK, "login.tmpl", gin.H{
		// "config": conf,
	})
}

func (r *rest) addFieldsToContext(ctx *gin.Context) {
	reqid := ctx.GetHeader(header.KeyRequestID)
	if reqid == "" {
		reqid = uuid.New().String()
	}

	c := ctx.Request.Context()
	c = appcontext.SetRequestId(c, reqid)
	c = appcontext.SetUserAgent(c, ctx.Request.Header.Get(header.KeyUserAgent))
	c = appcontext.SetAcceptLanguage(c, ctx.Request.Header.Get(header.KeyAcceptLanguage))
	c = appcontext.SetServiceVersion(c, r.conf.Meta.Version)
	c = appcontext.SetDeviceType(c, ctx.Request.Header.Get(header.KeyDeviceType))
	c = appcontext.SetCacheControl(c, ctx.Request.Header.Get(header.KeyCacheControl))
	c = appcontext.SetServiceName(c, ctx.Request.Header.Get(header.KeyServiceName))
	ctx.Request = ctx.Request.WithContext(c)
	ctx.Next()
}

func (r *rest) BodyLogger(ctx *gin.Context) {
	if r.conf.LogRequest {
		r.log.Info(ctx.Request.Context(),
			fmt.Sprintf(infoRequest, ctx.Request.RequestURI, ctx.Request.Method))
	}

	ctx.Next()
	if r.conf.LogResponse {
		if ctx.Writer.Status() < 300 {
			r.log.Info(ctx.Request.Context(),
				fmt.Sprintf(infoResponse, ctx.Request.RequestURI, ctx.Request.Method, ctx.Writer.Status()))
		} else {
			r.log.Error(ctx.Request.Context(),
				fmt.Sprintf(infoResponse, ctx.Request.RequestURI, ctx.Request.Method, ctx.Writer.Status()))
		}
	}
}

func (r *rest) VerifyUser(ctx *gin.Context) {
	user, err := r.verifyUserAuth(ctx)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	c := ctx.Request.Context()
	c = r.jwtAuth.SetUserAuthInfo(c, jwtAuth.UserAuthParam{
		User: user.ConvertToAuthUser(),
	})
	c = appcontext.SetUserId(c, int(user.ID))
	ctx.Request = ctx.Request.WithContext(c)

	ctx.Next()
}

func (r *rest) verifyUserAuth(ctx *gin.Context) (entity.User, error) {
	var (
		user entity.User
	)

	token := ctx.Request.Header.Get(header.KeyAuthorization)
	if token == "" {
		return entity.User{}, errors.NewWithCode(codes.CodeUnauthorized, "empty token")
	}

	jwtUer, err := r.jwtAuth.ValidateToken(token)
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeUnauthorized, "token invalid or token expire")
	}

	user, err = r.uc.User.Get(ctx.Request.Context(), entity.UserParam{
		ID: null.Int64From(jwtUer.ID),
		QueryOption: query.Option{
			IsActive: true,
		},
	})
	if err != nil {
		return entity.User{}, errors.NewWithCode(codes.CodeUnauthorized, "user does not exist")
	}

	return user, nil
}

func (r *rest) isAdmin(ctx *gin.Context) {
	user, err := r.verifyUserAuth(ctx)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	if user.RoleId.Int64 != entity.RoleIdSuperAdmin {
		r.httpRespError(ctx, errors.NewWithCode(codes.CodeUnauthorized, "error role to try access on admin resources"))
		return
	}

	ctx.Next()
}
