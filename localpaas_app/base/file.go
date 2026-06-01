package base

const (
	// 0755 grants read/write/execute for owner, read/execute for group/others
	DirModeDefault = 0755
)

type FileStatus string

const (
	FileStatusActive   FileStatus = "active"
	FileStatusPending  FileStatus = "pending"
	FileStatusDisabled FileStatus = "disabled"
)

var (
	AllFileStatuses = []FileStatus{FileStatusActive, FileStatusPending, FileStatusDisabled}
)

type FileType string

const (
	FileTypeSystemBackup FileType = "system-backup"
	FileTypeCache        FileType = "cache"
)

var (
	AllFileTypes = []FileType{FileTypeSystemBackup, FileTypeCache}
)

type FileStorageType string

const (
	FileStorageLocal FileStorageType = "local"
	FileStorageCloud FileStorageType = "cloud"
)

var (
	AllFileStorageTypes = []FileStorageType{FileStorageLocal, FileStorageCloud}
)
