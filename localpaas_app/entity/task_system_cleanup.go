package entity

type TaskSystemCleanupOutput struct {
	DBCleanup      *DBCleanupOutput      `json:"dbCleanup"`
	ClusterCleanup *ClusterCleanupOutput `json:"clusterCleanup"`
}

type DBCleanupOutput struct {
	Error string `json:"error,omitempty"`
}

type ClusterCleanupOutput struct {
	ImagesDeleted        int    `json:"imagesDeleted"`
	ImagesPruneError     string `json:"imagesPruneError,omitempty"`
	VolumesDeleted       int    `json:"volumesDeleted"`
	VolumesPruneError    string `json:"volumesPruneError,omitempty"`
	ContainersDeleted    int    `json:"containersDeleted"`
	ContainersPruneError string `json:"containersPruneError,omitempty"`
	NetworksDeleted      int    `json:"networksDeleted"`
	NetworksPruneError   string `json:"networksPruneError,omitempty"`
	SpaceReclaimed       uint64 `json:"spaceReclaimed"`
}

func (t *Task) OutputAsSystemCleanup() (*TaskSystemCleanupOutput, error) {
	return parseTaskOutputAs(t, func() *TaskSystemCleanupOutput { return &TaskSystemCleanupOutput{} })
}
