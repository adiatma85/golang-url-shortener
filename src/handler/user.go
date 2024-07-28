package handler

import (
	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/gin-gonic/gin"
)

// @Summary Get User List as an Admin
// @Description Get list all user Get User List as an Admin
// @Security BearerAuth
// @Tags Admin
// @Param limit query integer false "limit"
// @Param page query integer false "page"
// @Param disableLimit query boolean false "disable limit" Enums(true, false)
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=[]entity.User{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/admin/user [GET]
func (r *rest) GetListUserAsAdmin(ctx *gin.Context) {
	var param entity.UserParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	user, pg, err := r.uc.User.GetListAsAdmin(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, user, pg)
}

// @Summary Get User By ID
// @Description Get user details by user ID
// @Security BearerAuth
// @Tags User
// @Param user_id path integer true "user id"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.User{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/user/{user_id} [GET]
func (r *rest) GetUserByID(ctx *gin.Context) {
	var param entity.UserParam
	if err := r.BindParams(ctx, &param); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	user, err := r.uc.User.Get(ctx.Request.Context(), param)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, user, nil)
}

// @Summary Update One User
// @Description Update one user detail
// @Security BearerAuth
// @Tags Admin
// @Param user_id path integer true "user id"
// @Param user body entity.UpdateUserParam true "user data"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.User{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/admin/user/{user_id} [PUT]
func (r *rest) UpdateUser(ctx *gin.Context) {
	var updateParam entity.UpdateUserParam

	if err := r.Bind(ctx, &updateParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	var selectParam entity.UserParam
	if err := r.BindParams(ctx, &selectParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	if err := r.uc.User.Update(ctx.Request.Context(), updateParam, selectParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)
}

// @Summary Delete User
// @Description Soft delete user data
// @Security BearerAuth
// @Tags Admin
// @Param user_id path integer true "user id"
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.User{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/admin/user/{user_id} [DELETE]
func (r *rest) DeleteUser(ctx *gin.Context) {
	var selectParam entity.UserParam

	if err := r.BindParams(ctx, &selectParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	if err := r.uc.User.Delete(ctx.Request.Context(), selectParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)
}

// @Summary Get User Self Profile
// @Description Get user details self profile
// @Security BearerAuth
// @Tags User
// @Produce json
// @Success 200 {object} entity.HTTPResp{data=entity.User{}}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/user/profile [GET]
func (r *rest) UserProfile(ctx *gin.Context) {
	userProfile, err := r.uc.User.GetSelfProfile(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, userProfile, nil)
}

// @Summary Update User Profile
// @Description Update User Profile
// @Security BearerAuth
// @Tags User
// @Produce json
// @Success 200 {object} entity.HTTPResp{}
// @Param user_change_profile body entity.UpdateUserParam true "user change profile data"
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/user/profile [PUT]
func (r *rest) UpdateUserSelfProfile(ctx *gin.Context) {
	updateParam := entity.UpdateUserParam{}

	if err := r.Bind(ctx, &updateParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	err := r.uc.User.UpdateUserSelfProfile(ctx.Request.Context(), updateParam)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)

}

// @Summary Self Delete for User
// @Description Self Delete for User
// @Security BearerAuth
// @Tags User
// @Produce json
// @Success 200 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/user/profile [DELETE]
func (r *rest) UserSelfDelete(ctx *gin.Context) {
	err := r.uc.User.SelfDelete(ctx.Request.Context())
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)
}

// @Summary Change password profile for User
// @Description Change Password for User
// @Security BearerAuth
// @Tags User
// @Produce json
// @Param user_change_password body entity.ChangePasswordRequest true "user change password data"
// @Success 200 {object} entity.HTTPResp{}
// @Failure 500 {object} entity.HTTPResp{}
// @Router /v1/user/profile/change-password [PUT]
func (r *rest) UserChangePassword(ctx *gin.Context) {
	var updateParam entity.ChangePasswordRequest

	if err := r.Bind(ctx, &updateParam); err != nil {
		r.httpRespError(ctx, err)
		return
	}

	err := r.uc.User.ChangePassword(ctx.Request.Context(), updateParam)
	if err != nil {
		r.httpRespError(ctx, err)
		return
	}

	r.httpRespSuccess(ctx, codes.CodeSuccess, nil, nil)
}
