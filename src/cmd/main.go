package main

import (
	"context"
	"fmt"

	"github.com/adiatma85/golang-url-shortener/src/utils/config"
	"github.com/adiatma85/own-go-sdk/configreader"
	"github.com/adiatma85/own-go-sdk/instrument"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/sql"
)

const (
	configfile string = "./etc/cfg/conf.json"
)

// Read from the config file

// Get in SQL

// SQL Init

// Test

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
	// _ = sql.Init(cfg.SQL, log, instr)

	fmt.Println("Config adalah: ", cfg.Instrument)

	log.Info(context.Background(), "Ptesting log")

	fmt.Println("Database adalah: ", db.Follower())

	log.Info(context.Background(), db.Follower())
}
