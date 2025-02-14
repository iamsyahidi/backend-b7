package middleware

import (
	"backend-b7/models"
	"backend-b7/pkg/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// Response setting gin.JSON
func Response(c *gin.Context, req interface{}, res models.Response) {
	// LOGGER
	reqByte, _ := json.Marshal(req)
	resByte, _ := json.Marshal(res)
	logger.Infof("[backend-b7:log] [RequestURL] : %s, [RequestMethod] : %s, [RequestBody] : %s, [ResponseData] : %s", c.Request.RequestURI, c.Request.Method, string(reqByte), string(resByte))

	c.JSON(res.Code, res)
}
