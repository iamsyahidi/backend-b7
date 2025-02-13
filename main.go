package main

import (
	"backend-b7/controllers"
	"backend-b7/pkg/database"
	"backend-b7/repositories"
	"backend-b7/routes"
	"backend-b7/services"
	"context"
	"fmt"
	"log"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	db, err := database.InitGorm(ctx)
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.DebugMode)

	// Repositories
	meetRepository := repositories.NewMeetRepository(db)

	// Services
	meetService := services.NewMeetService(
		meetRepository,
		os.Getenv("ZOOM_CLIENT_ID"),
		os.Getenv("ZOOM_CLIENT_SECRET"),
		os.Getenv("ZOOM_AUTH_CODE"),
		os.Getenv("REDIRECT_URI"),
	)
	meetService.RequestAccessToken(os.Getenv("ZOOM_AUTH_CODE"))

	// Controllers
	meetController := controllers.NewMeetController(meetService)

	router := routes.NewRouter(meetController)

	port := fmt.Sprintf(":%s", os.Getenv("HTTP_PORT"))
	router.Run(port)
}
