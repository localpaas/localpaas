package projectservice

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

const (
	// 0755 grants read/write/execute for owner, read/execute for group/others
	// 0644 grants read/write for owner, read-only for group/others
	photoDirFileMode = 0o755
)

func (s *projectService) SaveProjectPhoto(
	_ context.Context,
	project *entity.Project,
	data []byte,
	fileExt string,
) error {
	dirPath := filepath.Join(config.Current.DataPathPhoto(), "project")
	err := os.MkdirAll(dirPath, photoDirFileMode)
	if err != nil {
		return fmt.Errorf("error creating project photo directory: %w", err)
	}

	if !strings.HasPrefix(fileExt, ".") {
		fileExt += "."
	}
	fileName := project.ID + fileExt
	fullPath := filepath.Join(dirPath, fileName)

	err = os.WriteFile(fullPath, data, photoDirFileMode)
	if err != nil {
		return fmt.Errorf("error writing project photo: %w", err)
	}

	// Save the photo path
	project.Photo = filepath.Join(config.Current.HttpPathPhoto(), "project", fileName)
	return nil
}
