package users

import (
	"user-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	GetUserLogin(ctx *gin.Context)
	GetUserByUUID(ctx *gin.Context)
	GetAllCustomer(ctx *gin.Context)
	GetAllAdmin(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type UserRoute struct {
	controller UserController
	group      *gin.RouterGroup
}

type IUserRoute interface {
	Run()
}

func NewUserRoute(controller UserController, group *gin.RouterGroup) IUserRoute {
	return &UserRoute{
		controller: controller,
		group:      group,
	}
}

func (u *UserRoute) Run() {
	group := u.group.Group("/auth")
	group.GET("/user", middlewares.Authenticate(), u.controller.GetUserLogin)
	group.GET("/cust", middlewares.Authenticate(), u.controller.GetAllCustomer)
	group.GET("/:uuid", middlewares.Authenticate(), u.controller.GetUserByUUID)
	group.POST("/signin", u.controller.SignIn)
	group.POST("/signup", u.controller.SignUp)
	group.PUT("/:uuid", middlewares.Authenticate(), u.controller.Update)
}
