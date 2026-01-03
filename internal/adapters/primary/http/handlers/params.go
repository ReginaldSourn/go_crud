package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseIDParam(c *gin.Context, name string) (int64, error) {
	return strconv.ParseInt(c.Param(name), 10, 64)
}
