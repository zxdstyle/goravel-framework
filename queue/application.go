package queue

import (
	configcontract "github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/queue"
)

type Application struct {
	config *Config
}

func NewApplication(config configcontract.Config) *Application {
	return &Application{
		config: NewConfig(config),
	}
}

func (app *Application) Worker(payloads ...*queue.Args) queue.Worker {
	defaultConnection := app.config.DefaultConnection()

	if len(payloads) == 0 || payloads[0] == nil {
		return NewWorker(app.config, 1, defaultConnection, app.config.Queue(defaultConnection, ""))
	}
	if payloads[0].Connection == "" {
		payloads[0].Connection = defaultConnection
	}
	if payloads[0].Concurrent == 0 {
		payloads[0].Concurrent = 1
	}

	return NewWorker(app.config, payloads[0].Concurrent, payloads[0].Connection, app.config.Queue(payloads[0].Connection, payloads[0].Queue))
}

func (app *Application) Register(jobs []queue.Job) error {
	if err := Register(jobs); err != nil {
		return err
	}

	return nil
}

func (app *Application) GetJobs() []queue.Job {
	var jobs []queue.Job
	/*JobRegistry.Range(func(key, value any) bool {
		jobs = append(jobs, value.(queue.Job))
		return true
	})*/

	return jobs
}

func (app *Application) Job(job queue.Job, args []queue.Arg) queue.Task {
	return NewTask(app.config, job, args)
}

func (app *Application) Chain(jobs []queue.Jobs) queue.Task {
	return NewChainTask(app.config, jobs)
}
