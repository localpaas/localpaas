package gocronqueue

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/infra/logging"
	"github.com/localpaas/localpaas/localpaas_app/infra/rediscache"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	queueConcurrency = 5

	defaultTaskCheckInterval = 60 * time.Second
)

var (
	ErrTaskProcessorNotFound = errors.New("task processor not found")
)

type TaskProcessorFunc func(taskID string, payload string) error

type Server struct {
	mu        sync.RWMutex
	scheduler gocron.Scheduler
	logger    logging.Logger
	jobMap    map[string]gocron.Job // task.ID -> job
	config    StartConfig
}

type StartConfig struct {
	TaskMap           map[base.TaskType]TaskProcessorFunc
	TaskCheckFunc     func(context.Context) ([]*entity.Task, error)
	TaskCheckInterval time.Duration
}

func NewServer(redisClient rediscache.Client, logger logging.Logger) (*Server, error) {
	return &Server{
		logger: logger,
	}, nil
}

func (s *Server) init() error {
	scheduler, err := gocron.NewScheduler(gocron.WithLimitConcurrentJobs(queueConcurrency, gocron.LimitModeWait))
	if err != nil {
		return apperrors.Wrap(err)
	}
	s.scheduler = scheduler
	s.jobMap = make(map[string]gocron.Job, 20) //nolint:mnd
	return nil
}

func (s *Server) Start(startConfig StartConfig) error {
	err := s.init()
	if err != nil {
		return apperrors.Wrap(err)
	}

	if startConfig.TaskCheckInterval <= 0 {
		startConfig.TaskCheckInterval = defaultTaskCheckInterval
	}

	s.config = startConfig
	s.scheduler.Start()

	// Start a job to scan for new tasks from DB
	_, err = s.scheduler.NewJob(
		gocron.DurationJob(startConfig.TaskCheckInterval),
		gocron.NewTask(s.scanTasks),
	)
	if err != nil {
		return apperrors.Wrap(err)
	}
	_ = s.scanTasks()

	return nil
}

func (s *Server) scanTasks() error {
	tasks, err := s.config.TaskCheckFunc(context.Background())
	if err != nil {
		s.logger.Errorf("failed to scan new tasks: %v", err)
		return apperrors.Wrap(err)
	}
	err = s.scheduleTasks(tasks)
	if err != nil {
		s.logger.Errorf("failed to schedule new tasks: %v", err)
		return apperrors.Wrap(err)
	}
	return nil
}

//nolint:unparam
func (s *Server) scheduleTasks(tasks []*entity.Task) error {
	timeNow := timeutil.NowUTC()
	for _, task := range tasks {
		var startAt gocron.OneTimeJobStartAtOption
		if task.RunAt.Before(timeNow) {
			startAt = gocron.OneTimeJobStartImmediately()
		} else {
			startAt = gocron.OneTimeJobStartDateTime(task.RunAt)
		}
		job, err := s.scheduler.NewJob(
			gocron.OneTimeJob(startAt),
			gocron.NewTask(func() {
				err := s.executeTask(task)
				if err != nil {
					s.logger.Errorf("failed to execute task '%v', id %s: %v", task.Type, task.ID, err)
				}
			}),
		)
		if err != nil {
			s.logger.Errorf("failed to schedule task %s: %v", task.ID, err)
		}
		err = s.addJob(job, task)
		if err != nil {
			s.logger.Errorf("failed to add job for task %s: %v", task.ID, err)
		}
	}
	return nil
}

func (s *Server) executeTask(task *entity.Task) error {
	defer func() {
		err := s.CancelTaskJob(task)
		if err != nil {
			s.logger.Errorf("failed to remove job for task %s: %v", task.ID, err)
		}
	}()

	execFunc := s.config.TaskMap[task.Type]
	if execFunc == nil {
		return fmt.Errorf("%w: task processor func not found for task type '%v'",
			ErrTaskProcessorNotFound, task.Type)
	}
	err := execFunc(task.ID, task.Args)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *Server) Shutdown() error {
	if s.scheduler == nil {
		return nil
	}
	err := s.scheduler.Shutdown()
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *Server) addJob(job gocron.Job, task *entity.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if currJob := s.jobMap[task.ID]; currJob != nil {
		if currJob == job {
			return nil
		}
		_ = s.scheduler.RemoveJob(job.ID())
	}
	s.jobMap[task.ID] = job
	return nil
}

func (s *Server) CancelTaskJob(task *entity.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if currJob := s.jobMap[task.ID]; currJob != nil {
		_ = s.scheduler.RemoveJob(currJob.ID())
	}
	delete(s.jobMap, task.ID)
	return nil
}
