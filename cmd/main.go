package cmd

import (
	"fmt"
	"net/http"
	"time"
	"user-service/constants"
	"user-service/database"
	"user-service/database/seeders"
	"user-service/domain/models"
	"user-service/helpers"
	"user-service/helpers/configs"
	userController "user-service/internal/controllers/users"
	"user-service/internal/middlewares"
	userRepo "user-service/internal/repositories/users"
	userRoute "user-service/internal/routes/users"
	userService "user-service/internal/services/users"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the serve",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()

		err := configs.Init(
			configs.WithConfigFolder(
				[]string{"./helpers/configs/"},
			),
			configs.WithConfigFile("config"),
			configs.WithConfigType("yaml"),
		)
		if err != nil {
			logrus.Errorf("Gagal inisiasi config: %s", err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}

		time.Local = loc

		cfg := configs.Get()

		db, err := database.InitDatabase(*cfg)
		if err != nil {
			panic(err)
		}

		err = db.AutoMigrate(
			&models.Role{},
			&models.User{},
		)

		if err != nil {
			panic(err)
		}

		seeders.NewSeederRegistry(db).Run()

		r := gin.Default()

		userRepo := userRepo.NewUserRepository(db)
		userService := userService.NewUserService(cfg, userRepo)
		userController := userController.NewUserController(r, userService)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(ctx *gin.Context) {
			ctx.JSON(http.StatusNotFound, helpers.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})

		router.GET("/", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, helpers.Response{
				Status:  constants.Success,
				Message: "Welcome to Hotel User Service",
			})
		})

		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-service-name, x-request-at, x-api-key")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		lmt := tollbooth.NewLimiter(cfg.Service.RateLimiterMaxRequest, &limiter.ExpirableOptions{
			DefaultExpirationTTL: time.Duration(cfg.Service.RateLimiterMaxSecond) * time.Second,
		})

		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")
		route := userRoute.NewUserRoute(userController, group)
		route.Run()

		port := fmt.Sprintf(":%d", cfg.Service.Port)
		router.Run(port)
	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
