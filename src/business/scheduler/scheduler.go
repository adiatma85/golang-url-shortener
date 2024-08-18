package scheduler

import (
	"context"
	"flag"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/adiatma85/golang-url-shortener/src/business/usecase"
	"github.com/adiatma85/golang-url-shortener/utils/config"
	"github.com/adiatma85/own-go-sdk/instrument"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/log"
	"github.com/go-co-op/gocron"
)

var (
	once = &sync.Once{}
)

type Interface interface {
	Run()
	TriggerScheduler(name string) error
}

type scheduler struct {
	cron       *gocron.Scheduler
	conf       config.SchedulerConfig
	metaconf   config.ApplicationMeta
	log        log.Interface
	jwtAuth    jwtAuth.Interface
	uc         *usecase.Usecase
	instrument instrument.Interface
}

type InitParam struct {
	Conf     config.SchedulerConfig
	MetaConf config.ApplicationMeta
	Log      log.Interface
	JwtAuth  jwtAuth.Interface
	Uc       *usecase.Usecase
	Instr    instrument.Interface
}

func Init(param InitParam) Interface {
	s := &scheduler{}

	once.Do(func() {
		cron := gocron.NewScheduler(time.UTC)
		cron.TagsUnique()

		s = &scheduler{
			cron:       cron,
			log:        param.Log,
			conf:       param.Conf,
			metaconf:   param.MetaConf,
			jwtAuth:    param.JwtAuth,
			uc:         param.Uc,
			instrument: param.Instr,
		}

		// read scheduler disabler by command or environment variable
		flagDisableScheduler := flag.Bool("disable-scheduler", false, "--disable-scheduler: Disable Scheduler on Application Startup")
		flag.Parse()
		envDisableScheduler, _ := strconv.ParseBool(os.Getenv("SERVICE_DISABLE_SCHEDULER"))

		// disable scheduler if the config is supplied in either one of command run flag, env variable, or json config
		schedulerDisabled := *flagDisableScheduler || envDisableScheduler || s.conf.DisableScheduler

		s.conf.DisableScheduler = schedulerDisabled

		s.assignScheduledTasks()
	})

	return s
}

func (s *scheduler) assignScheduledTasks() {
	// Assign the Scheduler in here

	// Assign the Scheduler for Insetring to the Database for Data Count
	s.AssignTask(s.conf.AssignCounterEndpoint, s.uc.Url.AssignCounterScheduler)
}

func (s *scheduler) Run() {
	if s.conf.DisableScheduler {
		s.log.Info(context.Background(), "Scheduler is disabled")
		return
	}

	s.cron.StartAsync()
	s.log.Info(context.Background(), "Scheduler is running")
}

func (s *scheduler) TriggerScheduler(name string) error {
	return s.cron.RunByTag(name)
}
