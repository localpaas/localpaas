package config

import "time"

type Tasks struct {
	Queue       TaskQueue   `toml:"queue"`
	Healthcheck Healthcheck `toml:"healthcheck"`
}

type TaskQueue struct {
	Concurrency        int           `toml:"concurrency" env:"LP_TASKS_QUEUE_CONCURRENCY" default:"5"`
	TaskCheckInterval  time.Duration `toml:"task_check_interval" env:"LP_TASKS_QUEUE_TASK_CHECK_INTERVAL" default:"1h"`
	TaskCreateInterval time.Duration `toml:"task_create_interval" env:"LP_TASKS_QUEUE_TASK_CREATE_INTERVAL" default:"1h"`
}

type Healthcheck struct {
	Enabled      bool          `toml:"enabled" env:"LP_TASKS_HEALTHCHECK_ENABLED" default:"true"`
	BaseInterval time.Duration `toml:"base_interval" env:"LP_TASKS_HEALTHCHECK_BASE_INTERVAL" default:"30s"`
}
