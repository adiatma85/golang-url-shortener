package main

import (
	"github.com/adiatma85/golang-url-shortener/src/business/domain"
	"github.com/adiatma85/golang-url-shortener/src/business/scheduler"
	"github.com/adiatma85/golang-url-shortener/src/business/usecase"
	"github.com/adiatma85/golang-url-shortener/src/handler"
	"github.com/adiatma85/golang-url-shortener/utils/config"
	"github.com/adiatma85/own-go-sdk/configreader"
	"github.com/adiatma85/own-go-sdk/instrument"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/parser"
	"github.com/adiatma85/own-go-sdk/redis"
	"github.com/adiatma85/own-go-sdk/sql"
)

// @contact.name   Rahmadhani Lucky Adiatma

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	configfile   string = "./etc/cfg/conf.json"
	templatefile string = "./etc/tpl/conf.template.json"
)

func main() {
	// init config
	cfg := config.Init()
	configreader := configreader.Init(configreader.Options{
		ConfigFile: configfile,
	})
	configreader.ReadConfig(&cfg)

	// init logger
	log := log.Init(cfg.Log)

	// init instrument tools
	instr := instrument.Init(cfg.Instrument)

	// init db conn
	db := sql.Init(cfg.SQL, log, instr)

	// init the parser
	parsers := parser.InitParser(log, cfg.Parser)

	// Init the jwt
	jwt := jwtAuth.Init(cfg.JwtAuth)

	// Init the Redis
	cache := redis.Init(cfg.Redis, log)

	// Init the domain
	d := domain.Init(domain.InitParam{Log: log, Db: db, Json: parsers.JSONParser()})

	// Init the usecase
	uc := usecase.Init(usecase.InitParam{Log: log, Dom: d, JwtAuth: jwt, Redis: cache})

	// Init the GIN
	rest := handler.Init(handler.InitParam{Conf: cfg.Gin, Json: parsers.JSONParser(), Uc: uc, Log: log, Instrument: instr, JwtAuth: jwt})

	// Init the Scheduler
	sch := scheduler.Init(scheduler.InitParam{Conf: cfg.Scheduler, MetaConf: cfg.Meta, Log: log, JwtAuth: jwt, Uc: uc, Instr: instr})

	// Run the Scheduelr
	sch.Run()

	// Run the HTTP Server
	rest.Run()
}
