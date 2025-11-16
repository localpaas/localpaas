package clusterhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/permission"
	"github.com/localpaas/localpaas/localpaas_app/usecase/cluster/imageuc/imagedto"
)

// To keep `apperrors` pkg imported and swag gen won't fail
type _ *apperrors.ErrorInfo

// ListImage Lists cluster images
// @Summary Lists cluster images
// @Description Lists cluster images
// @Tags    cluster_images
// @Produce json
// @Id      listImage
// @Param   status query string false "`status=<target>`"
// @Param   search query string false "`search=<target> (support *)`"
// @Param   pageOffset query int false "`pageOffset=offset`"
// @Param   pageLimit query int false "`pageLimit=limit`"
// @Param   sort query string false "`sort=[-]field1|field2...`"
// @Success 200 {object} imagedto.ListImageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/images [get]
func (h *ClusterHandler) ListImage(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleCluster,
		ResourceType:   base.ResourceTypeImage,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := imagedto.NewListImageReq()
	if err = h.ParseRequest(ctx, req, &req.Paging); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.imageUC.ListImage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetImage Gets image details
// @Summary Gets image details
// @Description Gets image details
// @Tags    cluster_images
// @Produce json
// @Id      getImage
// @Param   imageID path string true "image ID"
// @Success 200 {object} imagedto.GetImageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/images/{imageID} [get]
func (h *ClusterHandler) GetImage(ctx *gin.Context) {
	imageID, err := h.ParseStringParam(ctx, "imageID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleCluster,
		ResourceType:   base.ResourceTypeImage,
		ResourceID:     imageID,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := imagedto.NewGetImageReq()
	req.ImageID = imageID
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.imageUC.GetImage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetImageInspection Gets image details
// @Summary Gets image details
// @Description Gets image details
// @Tags    cluster_images
// @Produce json
// @Id      getImageInspection
// @Param   imageID path string true "image ID"
// @Success 200 {object} imagedto.GetImageInspectionResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/images/{imageID}/inspect [get]
func (h *ClusterHandler) GetImageInspection(ctx *gin.Context) {
	imageID, err := h.ParseStringParam(ctx, "imageID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleCluster,
		ResourceType:   base.ResourceTypeImage,
		ResourceID:     imageID,
		Action:         base.ActionTypeRead,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := imagedto.NewGetImageInspectionReq()
	req.ImageID = imageID
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.imageUC.GetImageInspection(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// CreateImage Creates an image
// @Summary Creates an image
// @Description Creates an image
// @Tags    cluster_images
// @Produce json
// @Id      createImage
// @Param   body body imagedto.CreateImageReq true "request data"
// @Success 200 {object} imagedto.CreateImageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/images [post]
func (h *ClusterHandler) CreateImage(ctx *gin.Context) {
	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleCluster,
		ResourceType:   base.ResourceTypeImage,
		Action:         base.ActionTypeWrite,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := imagedto.NewCreateImageReq()
	if err := h.ParseJSONBody(ctx, req); err != nil {
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.imageUC.CreateImage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// DeleteImage Deletes an image
// @Summary Deletes an image
// @Description Deletes an image
// @Tags    cluster_images
// @Produce json
// @Id      deleteImage
// @Param   imageID path string true "image ID"
// @Success 200 {object} imagedto.DeleteImageResp
// @Failure 400 {object} apperrors.ErrorInfo
// @Failure 500 {object} apperrors.ErrorInfo
// @Router  /cluster/images/{imageID} [delete]
func (h *ClusterHandler) DeleteImage(ctx *gin.Context) {
	imageID, err := h.ParseStringParam(ctx, "imageID")
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	auth, err := h.authHandler.GetCurrentAuth(ctx, &permission.AccessCheck{
		ResourceModule: base.ResourceModuleCluster,
		ResourceType:   base.ResourceTypeImage,
		ResourceID:     imageID,
		Action:         base.ActionTypeDelete,
	})
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	req := imagedto.NewDeleteImageReq()
	req.ImageID = imageID
	if err = h.ParseRequest(ctx, req, nil); err != nil { // to make sure Validate() to be called
		h.RenderError(ctx, err)
		return
	}

	resp, err := h.imageUC.DeleteImage(h.RequestCtx(ctx), auth, req)
	if err != nil {
		h.RenderError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
