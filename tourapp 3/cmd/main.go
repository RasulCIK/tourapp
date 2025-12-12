package main

import (
	"log"
	"tourapp/internal/controllers"
	"tourapp/internal/logger"
	"tourapp/internal/middleware"
	"tourapp/internal/models"
	"tourapp/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	
	dsn := "host=localhost user=postgres password=Barcalove2490 dbname=tourapp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database", err)
	}

	
	db.AutoMigrate(&models.User{})

	userRepo := repository.NewUserRepository(db)
	userController := controllers.NewUserController(userRepo)

	
	r := gin.Default()
	logger.Init()
	defer logger.Sync()
	logger.Logger.Info().Msg("Starting server on port 8080")


	
	r.POST("/register", userController.Register)
	r.POST("/login", userController.Login)

	
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuthMiddleware())
	authorized.GET("/users/:id", userController.GetUser)
	authorized.PUT("/users/:id", userController.UpdateUser)
	authorized.DELETE("/users/:id", userController.DeleteUser)

	
	r.Run(":8080")
}
