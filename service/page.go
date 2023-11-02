package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func BackgroundPage(c *gin.Context) {
	c.HTML(http.StatusOK, "background.html", nil)
}
