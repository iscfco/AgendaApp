package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeController struct {
	//service *services.LoginService
}

func NewHomeController() *HomeController {
	return &HomeController{}
}

// Create maneja la petición GET /
func (ctrl *HomeController) GetHome(c *gin.Context) {
	fmt.Println("hellooo33333333")
	c.String(http.StatusOK, "Home")
}
