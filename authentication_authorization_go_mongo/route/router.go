package route

import (
	"fmt"
	"log"
	"strconv"

	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/controllers"
	"github.com/amha-mersha/go_tasks/authentication_authorization_go_mongo/data"
	"github.com/gin-gonic/gin"
)

func Run(port int) {
	router := gin.Default()

	router.GET("/api/v1/tasks", controllers.GetTasks)
	router.GET("/api/v1/tasks/:id", controllers.GetTaskByID)
	router.POST("/api/v1/tasks", controllers.PostTask)
	router.PUT("/api/v1/tasks/:id", controllers.UpdateTask)
	router.DELETE("/api/v1/tasks/:id", controllers.DeleteTask)

	err := data.ConnecDB()
	if err != nil {
		fmt.Printf("Error occured when connecting to database. %v \n", err)
		return
	}
	router.Run("localhost:" + strconv.Itoa(port))
	log.Printf("Server up and running on port %d", port)
}
