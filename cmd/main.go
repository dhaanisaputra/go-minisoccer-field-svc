package cmd

import (
	"field-service/clients"
	"field-service/common/storage"
	"field-service/common/response"
	"field-service/config"
	"field-service/constants"
	"field-service/controllers"
	"field-service/domain/models"
	"field-service/middlewares"
	"field-service/repositories"
	"field-service/routes"
	"field-service/services"
	"fmt"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {

		_ = godotenv.Load()
		config.Init()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			panic(err)
		}

		time.Local = loc

		err = db.AutoMigrate(
			&models.Field{},
			&models.FieldSchedule{},
			&models.Time{},
		)

		if err != nil {
			panic(err)
		}

		// set seeder
		store := initStorage()
		client := clients.NewClientRegistry()
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, store)
		controller := controllers.NewControllerRegistry(service)

		router := gin.Default()
		router.Use(middlewares.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  constants.Error,
				Message: fmt.Sprintf("Path %s", http.StatusText(http.StatusNotFound)),
			})
		})
		router.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, response.Response{
				Status:  constants.Success,
				Message: "Welcome to field Service",
			})
		})

		// handle CORS
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, x-api-key, x-request-at, x-service-name")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		// handle rate limiter
		lmt := tollbooth.NewLimiter(
			config.Config.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.Config.RateLimiterTimeSecond) * time.Second,
			})

		router.Use(middlewares.RateLimiter(lmt))

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controller, group, client)
		route.Serve()

		port := fmt.Sprintf(":%d", config.Config.Port)
		router.Run(port)

	},
}

func Run() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}

func initStorage() storage.IStorageClient {
	client, err := storage.NewMinIOClient(config.Config.MinIO)
	if err != nil {
		panic(err)
	}
	return client
}
