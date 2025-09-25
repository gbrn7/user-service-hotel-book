package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
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
	"go.elastic.co/apm/module/apmgin/v2"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the serve",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()

		//setup logfile
		SetupLogfile()

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
			logrus.Errorf("failed to connect to database: %v", err)
			panic(err)
		}

		err = db.AutoMigrate(
			&models.Role{},
			&models.User{},
		)

		if err != nil {
			logrus.Errorf("failed to migrate database: %v", err)
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

		router.Use(apmgin.Middleware(router))

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

// CustomFormatter bikin format mirip log bawaan Go
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Format waktu: 2025/09/24 09:50:03
	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	// Format log: 2025/09/24 09:50:03 message
	logLine := []byte(timestamp + " " + entry.Message + "\n")
	return logLine, nil
}

func SetupLogfile() {
	logfile, err := os.OpenFile("./logs/user-service.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, logfile)
	logrus.SetOutput(mw)

	// Pakai formatter custom
	logrus.SetFormatter(&CustomFormatter{})
}
