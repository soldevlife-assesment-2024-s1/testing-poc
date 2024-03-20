package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/reugn/go-quartz/logger"
	"github.com/reugn/go-quartz/quartz"
)

// ExampleJobQueue demonstrates how to implement a custom job queue using Redis. https://github.com/reugn/go-quartz/blob/master/examples/queue/file_system.go

// printJob
type printJob struct {
	seconds int
}

var _ quartz.Job = (*printJob)(nil)

func (job *printJob) Execute(_ context.Context) error {
	logger.Infof("PrintJob: %d\n", job.seconds)
	return nil
}
func (job *printJob) Description() string {
	return fmt.Sprintf("PrintJob%s%d", quartz.Sep, job.seconds)
}

// scheduledPrintJob
type scheduledPrintJob struct {
	jobDetail   *quartz.JobDetail
	trigger     quartz.Trigger
	nextRunTime int64
}

// serializedJob
type serializedJob struct {
	Job         string                   `json:"job"`
	JobKey      string                   `json:"job_key"`
	Options     *quartz.JobDetailOptions `json:"job_options"`
	Trigger     string                   `json:"trigger"`
	NextRunTime int64                    `json:"next_run_time"`
}

var _ quartz.ScheduledJob = (*scheduledPrintJob)(nil)

func (job *scheduledPrintJob) JobDetail() *quartz.JobDetail {
	return job.jobDetail
}
func (job *scheduledPrintJob) Trigger() quartz.Trigger {
	return job.trigger
}
func (job *scheduledPrintJob) NextRunTime() int64 {
	return job.nextRunTime
}

// marshal returns the JSON encoding of the job.
func marshal(job quartz.ScheduledJob) ([]byte, error) {
	var serialized serializedJob
	serialized.Job = job.JobDetail().Job().Description()
	serialized.JobKey = job.JobDetail().JobKey().String()
	serialized.Options = job.JobDetail().Options()
	serialized.Trigger = job.Trigger().Description()
	serialized.NextRunTime = job.NextRunTime()
	return json.Marshal(serialized)
}

// unmarshal parses the JSON-encoded job.
func unmarshal(encoded []byte) (quartz.ScheduledJob, error) {
	var serialized serializedJob
	if err := json.Unmarshal(encoded, &serialized); err != nil {
		return nil, err
	}
	jobVals := strings.Split(serialized.Job, quartz.Sep)
	i, err := strconv.Atoi(jobVals[1])
	if err != nil {
		return nil, err
	}
	job := &printJob{i} // assuming we know the job type
	jobKeyVals := strings.Split(serialized.JobKey, quartz.Sep)
	jobKey := quartz.NewJobKeyWithGroup(jobKeyVals[1], jobKeyVals[0])
	jobDetail := quartz.NewJobDetailWithOptions(job, jobKey, serialized.Options)
	triggerOpts := strings.Split(serialized.Trigger, quartz.Sep)
	interval, _ := time.ParseDuration(triggerOpts[1])
	trigger := quartz.NewSimpleTrigger(interval) // assuming we know the trigger type
	return &scheduledPrintJob{
		jobDetail:   jobDetail,
		trigger:     trigger,
		nextRunTime: serialized.NextRunTime,
	}, nil
}

type jobQueue struct {
	Redis *redis.Client
}

func NewJobQueue(redis *redis.Client) quartz.JobQueue {
	return &jobQueue{
		Redis: redis,
	}
}

// Clear implements quartz.JobQueue.
func (j *jobQueue) Clear() error {
	ctx := context.Background()
	// _, err := j.Redis.FlushDB(ctx).Result()

	// Clear the job queue
	_, err := j.Redis.Del(ctx, "jobQueue").Result()

	return err
}

// Get implements quartz.JobQueue.
func (j *jobQueue) Get(jobKey *quartz.JobKey) (quartz.ScheduledJob, error) {
	ctx := context.Background()
	jobJSON, err := j.Redis.Get(ctx, "jobQueue").Result()
	if err != nil {
		return nil, err
	}
	// Deserialize jobJSON to quartz.ScheduledJob
	// ...

	job, err := unmarshal([]byte(jobJSON))
	if err != nil {
		return nil, err
	}

	return job, nil
}

// Head implements quartz.JobQueue.
func (j *jobQueue) Head() (quartz.ScheduledJob, error) {
	ctx := context.Background()
	jobJSON, err := j.Redis.LIndex(ctx, "jobQueue", 0).Result()
	if err != nil {
		return nil, err
	}
	// Deserialize jobJSON to quartz.ScheduledJob
	// ...

	job, err := unmarshal([]byte(jobJSON))
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Pop implements quartz.JobQueue.
func (j *jobQueue) Pop() (quartz.ScheduledJob, error) {
	ctx := context.Background()
	jobJSON, err := j.Redis.LPop(ctx, "jobQueue").Result()
	if err != nil {
		return nil, err
	}
	// Deserialize jobJSON to quartz.ScheduledJob
	// ...

	job, err := unmarshal([]byte(jobJSON))
	if err != nil {
		return nil, err
	}
	return job, nil
}

// Push implements quartz.JobQueue.
func (j *jobQueue) Push(job quartz.ScheduledJob) error {
	ctx := context.Background()
	// Serialize job to JSON
	jobJSON, err := marshal(job)
	if err != nil {
		return err
	}
	err = j.Redis.RPush(ctx, "jobQueue", jobJSON).Err()
	return err
}

// Remove implements quartz.JobQueue.
func (j *jobQueue) Remove(jobKey *quartz.JobKey) (quartz.ScheduledJob, error) {
	ctx := context.Background()
	jobJSON, err := j.Redis.Get(ctx, "jobQueue").Result()
	if err != nil {
		return nil, err
	}
	// Deserialize jobJSON to quartz.ScheduledJob
	// ...

	job, err := unmarshal([]byte(jobJSON))
	if err != nil {
		return nil, err
	}

	if job.JobDetail().JobKey().String() != jobKey.String() {
		return nil, fmt.Errorf("job not found")
	}

	_, err = j.Redis.Del(ctx, "jobQueue").Result()
	if err != nil {
		return nil, err
	}
	return job, nil
}

// ScheduledJobs implements quartz.JobQueue.
func (j *jobQueue) ScheduledJobs(matchers []quartz.Matcher[quartz.ScheduledJob]) []quartz.ScheduledJob {
	// Implement logic to retrieve scheduled jobs based on matchers from Redis

	ctx := context.Background()
	jobJSON, err := j.Redis.Get(ctx, "jobQueue").Result()
	if err != nil {
		return []quartz.ScheduledJob{}
	}
	// Deserialize jobJSON to quartz.ScheduledJob
	// ...

	job, err := unmarshal([]byte(jobJSON))
	if err != nil {
		return []quartz.ScheduledJob{}
	}

	return []quartz.ScheduledJob{job}
}

// Size implements quartz.JobQueue.
func (j *jobQueue) Size() int {
	ctx := context.Background()
	size, _ := j.Redis.LLen(ctx, "jobQueue").Result()
	return int(size)
}
