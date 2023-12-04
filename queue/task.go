package queue

import (
	"errors"

	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/support/carbon"
)

type Task struct {
	config     *Config
	connection string
	chain      bool
	delay      *carbon.Carbon
	driver     queue.Driver
	jobs       []queue.Jobs
	queue      string
}

func NewTask(config *Config, job queue.Job, payloads []queue.Payloads) *Task {
	return &Task{
		config:     config,
		connection: config.DefaultConnection(),
		driver:     NewDriver(config.DefaultConnection(), config),
		jobs: []queue.Jobs{
			{
				Job:      job,
				Payloads: payloads,
			},
		},
		queue: config.Queue(config.DefaultConnection(), ""),
	}
}

func NewChainTask(config *Config, jobs []queue.Jobs) *Task {
	return &Task{
		config:     config,
		connection: config.DefaultConnection(),
		chain:      true,
		driver:     NewDriver(config.DefaultConnection(), config),
		jobs:       jobs,
		queue:      config.Queue(config.DefaultConnection(), ""),
	}
}

// Delay sets a delay time for the task.
// Delay 设置任务的延迟时间。
func (receiver *Task) Delay(delay carbon.Carbon) queue.Task {
	receiver.delay = &delay

	return receiver
}

// Dispatch dispatches the task.
// Dispatch 调度任务。
func (receiver *Task) Dispatch() error {
	driver := receiver.config.Driver(receiver.connection)
	if driver == "" {
		return errors.New("unknown queue driver")
	}

	if receiver.chain {
		return receiver.driver.Bulk(receiver.jobs, receiver.queue)
	} else {
		job := receiver.jobs[0]
		return receiver.driver.Push(job.Job, job.Payloads, receiver.queue)
	}
}

// DispatchSync dispatches the task synchronously.
// DispatchSync 同步调度任务。
func (receiver *Task) DispatchSync() error {
	if receiver.chain {
		for _, job := range receiver.jobs {
			if err := Call(job.Job.Signature(), job.Payloads); err != nil {
				return err
			}
		}

		return nil
	} else {
		job := receiver.jobs[0]

		return Call(job.Job.Signature(), job.Payloads)
	}
}

// OnConnection sets the connection name.
// OnConnection 设置连接名称。
func (receiver *Task) OnConnection(connection string) queue.Task {
	receiver.connection = connection
	receiver.driver = NewDriver(connection, receiver.config)

	return receiver
}

// OnQueue sets the queue name.
// OnQueue 设置队列名称。
func (receiver *Task) OnQueue(queue string) queue.Task {
	receiver.queue = receiver.config.Queue(receiver.connection, queue)

	return receiver
}
