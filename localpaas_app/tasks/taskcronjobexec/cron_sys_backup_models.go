package taskcronjobexec

import (
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	sysBackupSqlPageSize      = 1000
	sysBackupSqlPageSizeSmall = 100
)

var (
	sysBackupDBModels = []*sysBackupDBModel{
		{
			Type:   "user",
			Model:  new([]*entity.User),
			Orders: []string{"id"},
		},
		{
			Type:   "acl_permission",
			Model:  new([]*entity.ACLPermission),
			Orders: []string{"subject_id", "resource_id"},
		},
		{
			Type:         "login_trusted_device",
			Model:        new([]*entity.LoginTrustedDevice),
			Orders:       []string{"user_id", "device_id"},
			NoSoftDelete: true,
		},
		{
			Type:     "setting",
			Model:    new([]*entity.Setting),
			Orders:   []string{"id"},
			PageSize: sysBackupSqlPageSizeSmall, // some settings may have big data field
		},
		{
			Type:   "project",
			Model:  new([]*entity.Project),
			Orders: []string{"id"},
		},
		{
			Type:   "project_tag",
			Model:  new([]*entity.ProjectTag),
			Orders: []string{"project_id", "tag"},
		},
		{
			Type:   "project_shared_setting",
			Model:  new([]*entity.ProjectSharedSetting),
			Orders: []string{"project_id", "setting_id"},
		},
		{
			Type:   "app",
			Model:  new([]*entity.App),
			Orders: []string{"id"},
		},
		{
			Type:   "app_tag",
			Model:  new([]*entity.AppTag),
			Orders: []string{"app_id", "tag"},
		},
		{
			Type:   "deployment",
			Model:  new([]*entity.Deployment),
			Orders: []string{"id"},
		},
		{
			Type:   "task",
			Model:  new([]*entity.Task),
			Orders: []string{"id"},
		},
		{
			Type:         "task_log",
			Model:        new([]*entity.TaskLog),
			Orders:       []string{"id"},
			NoSoftDelete: true,
		},
		{
			Type:         "sys_error",
			Model:        new([]*entity.SysError),
			Orders:       []string{"id"},
			NoSoftDelete: true,
		},
	}
)

type sysBackupDBModel struct {
	Type         string
	Model        any
	PageSize     int
	Orders       []string
	NoSoftDelete bool
}
