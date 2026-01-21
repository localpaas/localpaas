package server

import (
	"github.com/gin-gonic/gin"
)

func (s *HTTPServer) registerClusterRoutes(apiGroup *gin.RouterGroup) *gin.RouterGroup {
	clusterGroup := apiGroup.Group("/cluster")

	{ // node group
		nodeGroup := clusterGroup.Group("/nodes")
		// Nodes
		nodeGroup.GET("", s.handlerRegistry.clusterHandler.ListNode)
		nodeGroup.GET("/:nodeID", s.handlerRegistry.clusterHandler.GetNode)
		nodeGroup.GET("/:nodeID/inspect", s.handlerRegistry.clusterHandler.GetNodeInspection)
		nodeGroup.PUT("/:nodeID", s.handlerRegistry.clusterHandler.UpdateNode)
		nodeGroup.DELETE("/:nodeID", s.handlerRegistry.clusterHandler.DeleteNode)
		// Node join
		nodeGroup.POST("/join", s.handlerRegistry.clusterHandler.JoinNode)
		nodeGroup.GET("/join-command", s.handlerRegistry.clusterHandler.GetNodeJoinCommand)
	}
	{ // volume group
		volumeGroup := clusterGroup.Group("/volumes")
		// Volumes
		volumeGroup.GET("", s.handlerRegistry.clusterHandler.ListVolume)
		volumeGroup.GET("/:volumeID", s.handlerRegistry.clusterHandler.GetVolume)
		volumeGroup.GET("/:volumeID/inspect", s.handlerRegistry.clusterHandler.GetVolumeInspection)
		volumeGroup.POST("", s.handlerRegistry.clusterHandler.CreateVolume)
		volumeGroup.DELETE("/:volumeID", s.handlerRegistry.clusterHandler.DeleteVolume)
	}
	{ // image group
		imageGroup := clusterGroup.Group("/images")
		// Volumes
		imageGroup.GET("", s.handlerRegistry.clusterHandler.ListImage)
		imageGroup.GET("/:imageID", s.handlerRegistry.clusterHandler.GetImage)
		imageGroup.GET("/:imageID/inspect", s.handlerRegistry.clusterHandler.GetImageInspection)
		imageGroup.POST("", s.handlerRegistry.clusterHandler.CreateImage)
		imageGroup.DELETE("/:imageID", s.handlerRegistry.clusterHandler.DeleteImage)
	}

	return clusterGroup
}
