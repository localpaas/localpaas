package sysbackupserviceimpl

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/infra/database"
	"github.com/localpaas/localpaas/localpaas_app/pkg/applog"
	"github.com/localpaas/localpaas/localpaas_app/pkg/jsonl"
	"github.com/localpaas/localpaas/localpaas_app/pkg/timeutil"
)

const (
	sysBackupFilePageDataSize = 10 * 1024 * 1024 // KB
)

var (
	sysBackupFileModels = []*sysBackupFileModel{}
)

type sysBackupFileModel struct {
	Type         string
	PageDataSize int
	DirPath      func() string
}

type sysBackupFileEntry struct {
	Name string `json:"name"`
	Data []byte `json:"data"`
}

func (s *service) sysBackupFiles(
	ctx context.Context,
	_ database.IDB,
	jsonlW *jsonl.Writer,
	data *sysBackupData,
) (err error) {
	start := timeutil.NowUTC()
	_ = data.LogStore.Add(ctx, applog.NewWarnFrame("Start backing up static files...", applog.TsNow))

	defer func() {
		duration := timeutil.NowUTC().Sub(start)
		if err != nil {
			_ = data.LogStore.Add(ctx, applog.NewWarnFrame("Files backup finished in "+duration.String()+
				" with error: "+err.Error(), applog.TsNow))
		} else {
			_ = data.LogStore.Add(ctx, applog.NewOutFrame("Files backup finished in "+duration.String(),
				applog.TsNow))
		}
	}()

	for _, model := range sysBackupFileModels {
		dirPath := model.DirPath()
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return apperrors.Wrap(err)
		}

		maxPageDataSize := gofn.Coalesce(model.PageDataSize, sysBackupFilePageDataSize)
		pageDataSize := 0
		savingData := make([]*sysBackupFileEntry, 0, 20) //nolint:mnd
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			fileName := entry.Name()
			fileData, err := os.ReadFile(filepath.Join(dirPath, fileName))
			if err != nil {
				return apperrors.Wrap(err)
			}
			savingData = append(savingData, &sysBackupFileEntry{
				Name: entry.Name(),
				Data: fileData,
			})
			pageDataSize += len(fileData)
			if pageDataSize >= maxPageDataSize {
				err = jsonlW.WriteChunk(jsonl.NewChunk(model.Type, savingData))
				if err != nil {
					return apperrors.Wrap(err)
				}
				savingData = savingData[:0]
				pageDataSize = 0
			}
		}
		// Last page of data
		if len(savingData) > 0 {
			err = jsonlW.WriteChunk(jsonl.NewChunk(model.Type, savingData))
			if err != nil {
				return apperrors.Wrap(err)
			}
		}
	}

	return nil
}
