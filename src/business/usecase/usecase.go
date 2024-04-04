package usecase

import (
	"github.com/adiatma85/golang-url-shortener/src/business/domain"
	"github.com/adiatma85/golang-url-shortener/src/business/usecase/user"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
)

type Usecase struct {
	user user.Interface
}

type InitParam struct {
	Log     log.Interface
	Dom     *domain.Domain
	JwtAuth jwtAuth.Interface
}

func Init(param InitParam) *Usecase {
	usecase := &Usecase{
		user: user.Init(user.InitParam{Log: param.Log, User: param.Dom.User, JwtAuth: param.JwtAuth}),
	}

	return usecase
}
