package config

import "time"

type TaskQueue struct {
	Concurrency        int           `toml:"concurrency" env:"LP_TASKQUEUE_CONCURRENCY" default:"2"`
	TaskCheckInterval  time.Duration `toml:"task_check_interval" env:"LP_TASKQUEUE_TASK_CHECK_INTERVAL" default:"1m"`
	TaskCreateInterval time.Duration `toml:"task_create_interval" env:"LP_TASKQUEUE_TASK_CREATE_INTERVAL" default:"1m"`
}
