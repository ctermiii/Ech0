package router

import (
	"github.com/lin-snow/ech0/internal/di"
)

// setupResourceRoutes 设置资源路由
func setupResourceRoutes(appRouterGroup *AppRouterGroup, h *di.Handlers) {
	appRouterGroup.ResourceGroup.GET("rss", h.CommonHandler.GetRss)
}
