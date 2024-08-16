package route

import (
	"log"
	"strconv"
	"time"

	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/delivery/controllers"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/infrastructure"
	repositorie "github.com/amha-mersha/go_tasks/test-go-backend-task-manager/repositories"
	"github.com/amha-mersha/go_tasks/test-go-backend-task-manager/usecases"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Run(port int, database mongo.Database, timeout time.Duration, router *gin.Engine, usercollection string, taskcollection string) {
	public := router.Group("/api/v1")
	private := router.Group("/api/v1")
	private.Use(infrastructure.AuthMiddleWare("admin"))
	public.Use(infrastructure.AuthMiddleWare("user", "admin"))
	open := router.Group("/api/v1")
	CollectionUser := database.Collection(usercollection)
	CollectionTask := database.Collection(taskcollection)
	err := infrastructure.EstablisUniqueUsernameIndex(CollectionUser, "username")
	if err != nil {
		log.Println("Error", err)
	}

	taskRepository := repositorie.NewTaskRepository(CollectionTask)
	taskUsecase := usecases.NewTaskUsecase(&taskRepository, time.Second*5)
	userRepository := repositorie.NewUserRepository(CollectionUser)
	userUsecase := usecases.NewUserUsecase(&userRepository, time.Second*5)
	controller := controllers.NewController(&taskUsecase, &userUsecase)

	private.POST("/task", controller.PostTask)
	private.PUT("/task", controller.UpdateTask)
	private.DELETE("/task/:id/:userid", controller.DeleteTask)
	private.POST("/user/assign", controller.PostUserAssign)

	open.POST("/user/register", controller.PostUserRegister)
	open.POST("/user/login", controller.PostUserLogin)
	public.GET("/task", controller.GetTasks)
	public.GET("/task/:id", controller.GetTaskByID)

	router.Run("localhost:" + strconv.Itoa(port))
	log.Println("Server is running on port:", port)
}
