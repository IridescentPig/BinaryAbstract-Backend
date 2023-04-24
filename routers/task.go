package routers

import (
	"asset-management/app/api"
	"asset-management/middleware"
	"asset-management/utils"

	"github.com/gin-gonic/gin"
)

type taskRouter struct{}

var TaskRouter *taskRouter

func newTaskRouter() *taskRouter {
	return &taskRouter{}
}

func init() {
	TaskRouter = newTaskRouter()
}

func (task *taskRouter) Init(group *gin.RouterGroup) {
	task.routerUserTask(group.Group("/users"))
	task.routerDepartmentTask(group.Group("/department"))
}

func (task *taskRouter) routerUserTask(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	group.POST("/:user_id/assets/task", utils.Handler(api.TaskApi.CreateNewTask))
	group.GET("/:user_id/assets/tasks", utils.Handler(api.TaskApi.GetUserTaskList))
}

func (task *taskRouter) routerDepartmentTask(group *gin.RouterGroup) {
	group.Use(utils.Handler(middleware.JWTMiddleware()))
	group.GET("/:department_id/assets/tasks", utils.Handler(api.TaskApi.GetDepartmentTaskList))
}
