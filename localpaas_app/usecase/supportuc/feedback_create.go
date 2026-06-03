package supportuc

import (
	"context"

	"github.com/localpaas/localpaas/localpaas_app/basedto"
	"github.com/localpaas/localpaas/localpaas_app/usecase/supportuc/supportdto"
)

func (uc *UC) CreateFeedback(
	ctx context.Context,
	auth *basedto.Auth,
	req *supportdto.CreateFeedbackReq,
) (*supportdto.CreateFeedbackResp, error) {
	// TODO: add implementation
	return &supportdto.CreateFeedbackResp{}, nil
}
