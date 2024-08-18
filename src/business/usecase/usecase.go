package usecase

import (
	"github.com/adiatma85/golang-url-shortener/src/business/domain"
	"github.com/adiatma85/golang-url-shortener/src/business/usecase/url"
	"github.com/adiatma85/golang-url-shortener/src/business/usecase/user"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/redis"
)

type Usecase struct {
	User user.Interface
	Url  url.Interface
}

type InitParam struct {
	Log     log.Interface
	Dom     *domain.Domain
	JwtAuth jwtAuth.Interface
	Redis   redis.Interface
}

func Init(param InitParam) *Usecase {
	usecase := &Usecase{
		User: user.Init(user.InitParam{Log: param.Log, User: param.Dom.User, JwtAuth: param.JwtAuth}),
		Url:  url.Init(url.InitParam{Log: param.Log, Url: param.Dom.Url, JwtAuth: param.JwtAuth, Redis: param.Redis}),
	}

	return usecase
}
