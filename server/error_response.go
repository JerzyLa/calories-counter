package server

import (
	"calories-counter/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func handleErrorResponse(c *gin.Context, err error) {
	log.Info("error response caused by err: ", err)
	if e, ok := err.(common.ApiErr); ok {
		c.JSON(e.Code, gin.H{"error": e.Error()})
	} else {
		e = ErrInternalServerError
		c.JSON(e.Code, gin.H{"error": e.Error()})
	}
}
