package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(port int, database mongo.Database, timeout time.Duration, gin *gin.Engine) {
	public := gin.Group("/api/v1")
	private := gin.Group("")
	private.Use()

	private.GET("/task", controllers.GetTasks)
	private.GET("/task/:id", controllers.GetTaskByID)
	private.POST("/task/", controllers.PostTask)
	private.PUT("/task/:id", controllers.UpdateTask)
	private.DELETE("/task/:id", controllers.DeleteTask)
	private.POST("/user/assign", controllers.PostUserAssign)

	public.POST("/user/register", controllers.PostUserRegister)
	public.POST("/user/login", controllers.PostUserLogin)
}
