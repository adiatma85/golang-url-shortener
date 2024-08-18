package config

import (
	"time"

	"github.com/adiatma85/own-go-sdk/instrument"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/adiatma85/own-go-sdk/parser"
	"github.com/adiatma85/own-go-sdk/redis"
	"github.com/adiatma85/own-go-sdk/sql"
)

type Application struct {
	Log         log.Config
	Meta        ApplicationMeta
	Gin         GinConfig
	SQL         sql.Config
	Parser      parser.Options
	Instrument  instrument.Config
	Concurrency ConcurrencyConfig
	Redis       redis.Config
	JwtAuth     jwtAuth.Config
	Scheduler   SchedulerConfig
}

type ApplicationMeta struct {
	Title       string
	Description string
	Host        string
	BasePath    string
	Version     string
}

type GinConfig struct {
	Port            string
	Mode            string
	LogRequest      bool
	LogResponse     bool
	Timeout         time.Duration
	ShutdownTimeout time.Duration
	CORS            CORSConfig
	Meta            ApplicationMeta
	Swagger         SwaggerConfig
	Platform        PlatformConfig
	Dummy           DummyConfig
	Instrument      InstrumentConfig
	Profiler        ProfilerConfig
}

type PlatformConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type DummyConfig struct {
	Enabled bool
	Path    string
}

type InstrumentConfig struct {
	Metrics InstrumentMetricsConfig
}

type InstrumentMetricsConfig struct {
	Enabled   bool
	Path      string
	BasicAuth BasicAuthConf
}

type CORSConfig struct {
	Mode string
}
type SwaggerConfig struct {
	Enabled    bool
	Path       string
	BasicAuth  BasicAuthConf
	IsDarkMode bool
}

type ProfilerConfig struct {
	Pprof ProfilerPprofConfig
}

type ProfilerPprofConfig struct {
	Enabled    bool
	PathPrefix string
	BasicAuth  BasicAuthConf
}

type BasicAuthConf struct {
	Username string
	Password string
}

type ConcurrencyConfig struct {
	MetricsLimit int
}

type SchedulerTaskConf struct {
	Name          string
	Enabled       bool
	TimeType      string
	Interval      time.Duration
	ScheduledTime string
}

type SchedulerConfig struct {
	DisableScheduler      bool
	AssignCounterEndpoint SchedulerTaskConf
}

func Init() Application {
	return Application{}
}
