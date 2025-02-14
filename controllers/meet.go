package controllers

import (
	"backend-b7/middleware"
	"backend-b7/models"
	"backend-b7/pkg/logger"
	"backend-b7/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type meetController struct {
	meetService services.MeetServiceInterface
}

type MeetControllerInterface interface {
	CreateMeet(c *gin.Context)
	GetMeets(c *gin.Context)
	GetMeetById(c *gin.Context)
	UpdateMeet(c *gin.Context)
	DeleteMeet(c *gin.Context)

	HandleRedirect(c *gin.Context)
}

func NewMeetController(meetService services.MeetServiceInterface) MeetControllerInterface {
	return &meetController{
		meetService: meetService,
	}
}

func (pc *meetController) CreateMeet(c *gin.Context) {
	var meetRegister models.ZoomMeetRegister
	if err := c.ShouldBindJSON(&meetRegister); err != nil {
		middleware.Response(c, meetRegister, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	meet := models.ZoomMeet{
		UserID:    meetRegister.UserID,
		Topic:     meetRegister.Topic,
		StartTime: meetRegister.StartTime,
		Duration:  meetRegister.Duration,
		Status:    meetRegister.Status,
	}

	response, err := pc.meetService.CreateMeet(&meet)
	if err != nil {
		logger.Err(err.Error())
		middleware.Response(c, meetRegister, models.Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    nil,
		})
		return
	}

	middleware.Response(c, meetRegister, *response)
}

func (pc *meetController) GetMeets(c *gin.Context) {
	filter := c.Request.URL.Query()
	response, err := pc.meetService.GetMeets(filter)
	if err != nil {
		logger.Err(err.Error())
		middleware.Response(c, filter, models.Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    nil,
		})
		return
	}

	middleware.Response(c, filter, *response)
}

func (pc *meetController) GetMeetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.Response(c, id, models.Response{
			Code:    http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Data:    nil,
		})
		return
	}

	response, err := pc.meetService.GetMeetById(id)
	if err != nil {
		logger.Err(err.Error())
		middleware.Response(c, id, models.Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    nil,
		})
		return
	}

	middleware.Response(c, id, *response)
}

func (pc *meetController) UpdateMeet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.Response(c, id, models.Response{
			Code:    http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Data:    nil,
		})
		return
	}

	var meetUpdate models.ZoomMeetUpdate
	if err := c.ShouldBindJSON(&meetUpdate); err != nil {
		middleware.Response(c, meetUpdate, models.Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	meetUpdate.ID = id
	response, err := pc.meetService.UpdateMeet(&meetUpdate)
	if err != nil {
		logger.Err(err.Error())
		middleware.Response(c, meetUpdate, models.Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    nil,
		})
		return
	}

	middleware.Response(c, meetUpdate, *response)
}

func (pc *meetController) DeleteMeet(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		middleware.Response(c, id, models.Response{
			Code:    http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Data:    nil,
		})
		return
	}

	meetDelete := models.ZoomMeetUpdate{
		ID: id,
	}
	response, err := pc.meetService.DeleteMeet(&meetDelete)
	if err != nil {
		logger.Err(err.Error())
		middleware.Response(c, id, models.Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    nil,
		})
		return
	}

	middleware.Response(c, id, *response)
}

func (pc *meetController) HandleRedirect(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		middleware.Response(c, code, models.Response{
			Code:    http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
			Data:    nil,
		})
		return
	}

	_, err := pc.meetService.RequestAccessToken(code)
	if err != nil {
		logger.Err(err.Error())
		middleware.Response(c, code, models.Response{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
			Data:    nil,
		})
		return
	}

	middleware.Response(c, code, models.Response{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    nil,
	})
}
