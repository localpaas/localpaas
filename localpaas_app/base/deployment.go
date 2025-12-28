package base

type DeploymentStatus string

const (
	DeploymentStatusNotStarted DeploymentStatus = "not-started"
	DeploymentStatusInProgress DeploymentStatus = "in-progress"
	DeploymentStatusCanceled   DeploymentStatus = "canceled"
	DeploymentStatusDone       DeploymentStatus = "done"
	DeploymentStatusFailed     DeploymentStatus = "failed"
)

var (
	AllDeploymentStatuses = []DeploymentStatus{DeploymentStatusNotStarted, DeploymentStatusInProgress,
		DeploymentStatusCanceled, DeploymentStatusDone, DeploymentStatusFailed}
)
