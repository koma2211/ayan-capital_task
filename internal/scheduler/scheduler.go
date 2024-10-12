package scheduler

import (
	"context"

	"github.com/jasonlvhit/gocron"
	"github.com/koma2211/ayan-capital_task/internal/service"
)

type JobSheduler struct {
	services *service.Service
}

func NewJobSheduler(services *service.Service) *JobSheduler {
	return &JobSheduler{services: services}
}

func (jobScheduler *JobSheduler) StartSheduler(ctx context.Context) {
	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(func ()  {
		jobScheduler.services.Eventer.NotifyAllEvents(ctx)
	})
	<-s.Start()
}
