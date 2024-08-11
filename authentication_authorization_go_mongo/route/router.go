package route

import (
	"fmt"
	"log"
	"strconv"

	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/controllers"
	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/data"
	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/middleware"
	"github.com/gin-gonic/gin"
)

func Run(port int) {
	router := gin.Default()
	private := router.Group("/api/v1/tasks")
	private.Use(middleware.AuthMiddleware())
	public := router.Group("/api/v1/user")

	private.GET("/", controllers.GetTasks)
	private.GET("/:id", controllers.GetTaskByID)
	private.POST("/", controllers.PostTask)
	private.POST("/assign", controllers.PostUserAssign)
	private.PUT("/:id", controllers.UpdateTask)
	private.DELETE("/:id", controllers.DeleteTask)

	public.POST("/register", controllers.PostUserRegister)
	public.POST("/login", controllers.PostUserLogin)

	err := data.ConnecDB()
	if err != nil {
		fmt.Printf("Error occured when connecting to database. %v \n", err)
		return
	}
	router.Run("localhost:" + strconv.Itoa(port))
	log.Printf("Server up and running on port %d", port)
}
