package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/adiatma85/golang-url-shortener/src/business/entity"
	"github.com/adiatma85/golang-url-shortener/utils/config"
	"github.com/adiatma85/own-go-sdk/appcontext"
	"github.com/adiatma85/own-go-sdk/codes"
	"github.com/adiatma85/own-go-sdk/errors"
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/google/uuid"
)

const (
	schedulerUserAgent string = "Cron Scheduler : %s"

	schedulerAssignError string = "Assigning Scheduler %s error: %s"

	schedulerRunning       string = "Scheduler %s is running"
	schedulerDoneError     string = "Scheduler %s is error: %v"
	schedulerDoneSuccess   string = "Scheduler %s is success"
	schedulerTimeExecution string = "Scheduler %s done in %v"

	schedulerTimeTypeExact    string = "daily"    // Everyday at exact time
	schedulerTimeTypeInterval string = "interval" // Every given interval e.g 5 * time.Minute or so
)

type handlerFunc func(ctx context.Context) error

func (s *scheduler) createContext(conf config.SchedulerTaskConf) context.Context {
	ctx := context.Background()
	ctx = appcontext.SetUserAgent(ctx, fmt.Sprintf(schedulerUserAgent, conf.Name))
	ctx = appcontext.SetRequestId(ctx, uuid.New().String())
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())
	ctx = appcontext.SetServiceVersion(ctx, s.metaconf.Version)

	schedulerUser := entity.User{
		ID:          entity.SystemUser,
		DisplayName: entity.SchedulerUserName,
		Username:    entity.SchedulerUserName,
	}

	ctx = s.jwtAuth.SetUserAuthInfo(ctx, jwtAuth.UserAuthParam{
		User: schedulerUser.ConvertToAuthUser(),
	})

	return ctx
}

func (s *scheduler) AssignTask(conf config.SchedulerTaskConf, task handlerFunc) {
	if conf.Enabled {
		var err error
		ctx := context.Background()
		schedulerFunc := s.taskWrapper(conf, task)

		switch conf.TimeType {
		case schedulerTimeTypeInterval:
			_, err = s.cron.Every(conf.Interval).Tag(conf.Name).Do(schedulerFunc)
		case schedulerTimeTypeExact:
			_, err = s.cron.Every(1).Day().Tag(conf.Name).At(conf.ScheduledTime).Do(schedulerFunc)
		default:
			err = errors.NewWithCode(codes.CodeInternalServerError, "Unknown Scheduler Task Time Type")
		}

		if err != nil {
			s.log.Fatal(ctx, fmt.Sprintf(schedulerAssignError, conf.Name, err.Error()))
		}

	}
}

func (s *scheduler) taskWrapper(conf config.SchedulerTaskConf, task handlerFunc) func() {
	return func() {
		if s.instrument.IsEnabled() {
			s.instrument.SchedulerRunningCounter(conf.Name)
			timer := s.instrument.SchedulerRunningTimer(conf.Name)
			defer timer.ObserveDuration()
		}

		ctx := s.createContext(conf)
		s.log.Info(ctx, fmt.Sprintf(schedulerRunning, conf.Name))
		if err := task(ctx); err != nil {
			s.log.Error(ctx, fmt.Sprintf(schedulerDoneError, conf.Name, err))
		} else {
			s.log.Info(ctx, fmt.Sprintf(schedulerDoneSuccess, conf.Name))
		}

		startTime := appcontext.GetRequestStartTime(ctx)
		s.log.Info(ctx, fmt.Sprintf(schedulerTimeExecution, conf.Name, time.Since(startTime)))
	}
}
