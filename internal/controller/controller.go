package controller

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"golibrary/internal/service"
)

type Controller struct {
	service *service.Service
	logger  *zap.SugaredLogger
}

func NewController(service *service.Service, log *zap.SugaredLogger) *Controller {
	return &Controller{
		service: service,
		logger:  log,
	}
}

func (c *Controller) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	users := router.Group("/users")
	{
		users.POST("/create", c.CreateUser)
		users.GET("/list", c.ListUsers)
		users.GET("/getByID", c.GetUserByID)
		users.GET("/getBySurname", c.GetUserBySurname)
		users.PUT("/update", c.UpdateUser)
	}

	friendships := router.Group("/friendships")
	{
		friendships.POST("/create", c.CreateFriendship)
		friendships.GET("/list", c.ListUserFriends)
	}

	return router
}
