package jobs

import "github.com/robfig/cron/v3"

var jobs Jobs

type Jobs struct {
	*cron.Cron
}

func init() {
	jobs = Jobs{
		Cron: cron.New(),
	}
	jobs.Start()
}
