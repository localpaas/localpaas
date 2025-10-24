package userservice

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
	photoDirFileMode = 0644
)

func (s *userService) SaveUserPhoto(_ context.Context, user *entity.User, data []byte, fileExt string) error {
	dirPath := config.Current.App.DataPathUserPhoto()
	err := os.MkdirAll(dirPath, photoDirFileMode)
	if err != nil {
		return fmt.Errorf("error creating user photo directory: %w", err)
	}

	if !strings.HasPrefix(fileExt, ".") {
		fileExt += "."
	}
	fileName := user.ID + fileExt
	fullPath := filepath.Join(dirPath, fileName)

	err = os.WriteFile(fullPath, data, photoDirFileMode)
	if err != nil {
		return fmt.Errorf("error writing user photo: %w", err)
	}

	// Save the photo path
	user.Photo = config.Current.App.HttpPathUserPhoto() + fileName
	return nil
}
