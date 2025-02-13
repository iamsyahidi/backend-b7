package routes

import (
	"backend-b7/controllers"
	"backend-b7/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(meetController controllers.MeetControllerInterface) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	baseRouter := router.Group("/v1")

	//* health
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "health",
		})
	})

	//* redirect
	redirect := baseRouter.Group("/redirect")
	redirect.GET("", meetController.HandleRedirect)

	//* meets
	meets := baseRouter.Group("/meets")
	meets.POST("", meetController.CreateMeet)
	meets.GET("", meetController.GetMeets)
	meets.GET("/:id", meetController.GetMeetById)
	meets.PUT("/:id", meetController.UpdateMeet)
	meets.DELETE("/:id", meetController.DeleteMeet)

	return router
}
