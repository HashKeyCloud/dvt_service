package service

import "github.com/gin-gonic/gin"

func (s *APIService) healthy(c *gin.Context) {
	c.JSON(200, resultResponse{
		Code: 200,
		Msg:  "healthy",
		Data: "healthy",
	})
}
