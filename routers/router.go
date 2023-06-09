package routers

import (
	"asset-management/app/api"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type router struct{}

var Router *router

func newRouter() *router {
	return &router{}
}

func init() {
	Router = newRouter()
}

func (router *router) Init(r *gin.Engine) *gin.Engine {
	r.NoRoute(utils.Handler(RouteNotFound))
	r.NoMethod(utils.Handler(MethodNotFound))

	UserRouter.Init(r.Group("/user"))
	UsersRouter.Init(r.Group("/users"))
	EntityRouter.Init(r.Group("/entity"))
	AssetClassRouter.Init(r.Group("/department"))
	AssetRouter.Init(r.Group("/department"))
	TaskRouter.Init(r.Group(""))
	LogRouter.Init(r.Group("/entity"))
	OssRouter.Init(r.Group(""))
	AsyncRouter.Init(r.Group(""))
	r.GET("/asset/:asset_id/info", utils.Handler(api.AssetApi.GetAssetInfoByScan))
	return r
}

func RouteNotFound(ctx *utils.Context) {
	ctx.NotFound(1, "Router not found.")
}

func MethodNotFound(ctx *utils.Context) {
	ctx.NotFound(1, "Method not found.")
}
