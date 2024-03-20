package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/quartz"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create scheduler
	sched := quartz.NewStdScheduler()

	// async start scheduler
	sched.Start(ctx)

	go jobsAndScheduler(sched)

	// // create jobs
	// cronTrigger, _ := quartz.NewCronTrigger("* * * * * *")
	// shellJob := job.NewShellJob("ls -la")

	// request, _ := http.NewRequest(http.MethodGet, "https://worldtimeapi.org/api/timezone/utc", nil)
	// curlJob := job.NewCurlJob(request)

	// functionJob := job.NewFunctionJob(func(_ context.Context) (int, error) { return 42, nil })

	// // register jobs to scheduler
	// sched.ScheduleJob(quartz.NewJobDetail(shellJob, quartz.NewJobKey("shellJob")),
	// 	cronTrigger)
	// sched.ScheduleJob(quartz.NewJobDetail(curlJob, quartz.NewJobKey("curlJob")),
	// 	quartz.NewSimpleTrigger(time.Second*7))
	// sched.ScheduleJob(quartz.NewJobDetail(functionJob, quartz.NewJobKey("functionJob")),
	// 	quartz.NewSimpleTrigger(time.Second*5))

	// stop scheduler
	defer sched.Stop()
	time.Sleep(time.Minute) // Sleep for 1 minute

	// wait for all workers to exit
	sched.Wait(ctx)
}

func jobsAndScheduler(sched quartz.Scheduler) {
	// create jobs
	cronTrigger, _ := quartz.NewCronTrigger("1/5 * * * * *")
	shellJob := job.NewShellJob("ls -la")

	// request, _ := http.NewRequest(http.MethodGet, "https://worldtimeapi.org/api/timezone/utc", nil)
	// curlJob := job.NewCurlJob(request)

	functionJob := job.NewFunctionJob(func(_ context.Context) (int, error) {
		fmt.Println("test hello job trigger once")
		return 42, nil
	})

	dateEnd := time.Now().Add(time.Second * 10)
	duration := durationCalculation(dateEnd)

	// register jobs to scheduler
	sched.ScheduleJob(quartz.NewJobDetail(shellJob, quartz.NewJobKey("shellJob")),
		cronTrigger)
	sched.ScheduleJob(quartz.NewJobDetail(functionJob, quartz.NewJobKey("curlJob")),
		quartz.NewRunOnceTrigger(duration))
	// sched.ScheduleJob(quartz.NewJobDetail(functionJob, quartz.NewJobKey("functionJob")),
	// 	quartz.NewSimpleTrigger(time.Second*5))
}

// create duration calculator job that calculate duration from specific time to time now
func durationCalculation(dateEnd time.Time) time.Duration {
	dateStart := time.Now()
	return dateEnd.Sub(dateStart)
}
