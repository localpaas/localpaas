package queue

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/localpaas_app/entity/cacheentity"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
)

type TaskExecData struct {
	Task *entity.Task

	// RefObjects can be used as a cache to store objects
	RefObjects *entity.RefObjects
	LogStore   *applog.Store

	NonCancelable bool
	NonRetryable  bool
	Canceled      bool
	Done          bool

	// Callback functions
	OnCommandFunc         func(base.TaskCommand, ...any)
	OnEndTransactionFunc  func()
	OnPostTransactionFunc func()
}

func (t *TaskExecData) IsCanceled() bool {
	return t.Canceled
}

func (t *TaskExecData) IsDone() bool {
	return t.Done
}

func (t *TaskExecData) AddRefObjects(refObjects *entity.RefObjects) {
	if t.RefObjects == nil {
		t.RefObjects = refObjects
	} else {
		t.RefObjects.AddRefObjects(refObjects)
	}
}

func (t *TaskExecData) OnCommand(fn func(base.TaskCommand, ...any)) {
	if t.OnCommandFunc == nil {
		t.OnCommandFunc = fn
		return
	}
	currFunc := t.OnCommandFunc
	t.OnCommandFunc = func(cmd base.TaskCommand, args ...any) {
		currFunc(cmd, args...)
		fn(cmd, args...)
	}
}

func (t *TaskExecData) OnEndTransaction(fn func()) {
	if t.OnEndTransactionFunc == nil {
		t.OnEndTransactionFunc = fn
		return
	}
	currFunc := t.OnEndTransactionFunc
	t.OnEndTransactionFunc = func() {
		currFunc()
		fn()
	}
}

func (t *TaskExecData) OnPostTransaction(fn func()) {
	if t.OnPostTransactionFunc == nil {
		t.OnPostTransactionFunc = fn
		return
	}
	currFunc := t.OnPostTransactionFunc
	t.OnPostTransactionFunc = func() {
		currFunc()
		fn()
	}
}

type TaskExecFunc func(context.Context, database.Tx, *TaskExecData) error

type HealthcheckExecData struct {
	HealthcheckSetting *entity.Setting
	Healthcheck        *entity.Healthcheck
	Task               *entity.Task
	Project            *entity.Project
	App                *entity.App

	// RefObjects can be used as a cache to store objects
	RefObjects    *entity.RefObjects
	NotifEventMap map[string]*cacheentity.HealthcheckNotifEvent
}

func (t *HealthcheckExecData) AddRefObjects(refObjects *entity.RefObjects) {
	if t.RefObjects == nil {
		t.RefObjects = refObjects
	} else {
		t.RefObjects.AddRefObjects(refObjects)
	}
}

type HealthcheckExecFunc func(context.Context, *HealthcheckExecData) error
