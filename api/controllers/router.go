package controllers

import "github.com/gin-gonic/gin"

type Router interface {
	SetRoutes(group *gin.RouterGroup)
}
