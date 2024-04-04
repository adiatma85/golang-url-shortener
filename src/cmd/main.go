package main

import (
	"context"

	"github.com/adiatma85/golang-url-shortener/utils/config"
	"github.com/adiatma85/own-go-sdk/configreader"
	"github.com/adiatma85/own-go-sdk/instrument"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/parser"
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

	log.Info(context.Background(), db)
	log.Info(context.Background(), parsers)
	log.Info(context.Background(), jwt)
}
