package handler

import (
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/gin-gonic/gin"
)

// @Summary Create new Url
// @Description Create new Url
// @Security BearerAuth
// @Tags Url
// @Param data body entity.CreateUrlParam true "Input of New URL Data"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.Url{}}
// @Failure 400 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/url [POST]
func (r *rest) CreateUrl(ctx *gin.Context) {
	var insertParam entity.CreateUrlParam
	if err := r.Bind(ctx, &insertParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	newUrl, err := r.uc.Url.Create(ctx.Request.Context(), insertParam)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, newUrl, nil)
}

// @Summary Get Url List
// @Description Get Url List
// @Security BearerAuth
// @Tags Url
// @Param limit query integer false "limit"
// @Param page query integer false "page"
// @Param disableLimit query boolean false "disable limit" Enums(true, false)
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=[]entity.Url{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/url [GET]
func (r *rest) GetListUrl(ctx *gin.Context) {
	var param entity.UrlParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	urls, pg, err := r.uc.Url.GetList(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, urls, pg)
}

// @Summary Get Url List as an Admin
// @Description Get Url List as an Admin
// @Security BearerAuth
// @Tags Admin
// @Param limit query integer false "limit"
// @Param page query integer false "page"
// @Param disableLimit query boolean false "disable limit" Enums(true, false)
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=[]entity.Url{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/admin/url [GET]
func (r *rest) GetListUrlAdmin(ctx *gin.Context) {
	var param entity.UrlParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	urls, pg, err := r.uc.Url.GetListAsAdmin(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, urls, pg)
}

// @Summary Get Url By ID
// @Description Get url details by url ID
// @Security BearerAuth
// @Tags Url
// @Param url_id path integer true "url id"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.Url{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/url/{url_id} [GET]
func (r *rest) GetUrlByID(ctx *gin.Context) {
	var param entity.UrlParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	url, err := r.uc.Url.Get(ctx, param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, url, nil)
}

// @Summary Get Url By Shortened Url
// @Description Get url details by shortened version
// @Security BearerAuth
// @Tags Url
// @Param shorten_url path string true "shorten url"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.Url{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/shortened-url/{shorten_url} [GET]
func (r *rest) GetByShortenedUrl(ctx *gin.Context) {
	var param entity.UrlParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	url, err := r.uc.Url.Get(ctx, param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, url, nil)
}

// @Summary Update URL by URL ID
// @Description Update Url by Database Entry ID
// @Security BearerAuth
// @Tags Url
// @Param url_id path integer true "url id"
// @Produce json
// @Success 200 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/url/{url_id} [PUT]
func (r *rest) UpdateUrl(ctx *gin.Context) {
	var (
		updateParam entity.UpdateUrlParam
		selectParam entity.UrlParam
	)

	if err := r.Bind(ctx, &updateParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	if err := r.BindParams(ctx, &selectParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	if err := r.uc.Url.Update(ctx.Request.Context(), updateParam, selectParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)
}

// @Summary Delete URL by URL ID
// @Description Delete Url by Database Entry ID
// @Security BearerAuth
// @Tags Url
// @Param url_id path integer true "url id"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.Url{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/url/{url_id} [DELETE]
func (r *rest) DeleteUrl(ctx *gin.Context) {
	var param entity.UrlParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	if err := r.uc.Url.Delete(ctx.Request.Context(), param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)
}
