package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type shortCodeCreator interface {
}

type shortCodeFollower interface {
}

type analiticGetter interface {
}

type Router struct {
	Router            *ginext.Engine
	shortCodeCreator  shortCodeCreator
	shortCodeFollower shortCodeFollower
	analiticGetter    analiticGetter
}

func New(router *ginext.Engine, shshortCodeCreator shortCodeCreator, shortCodeFollower shortCodeFollower, analiticGetter analiticGetter) *Router {
	return &Router{
		Router:            router,
		shortCodeCreator:  shshortCodeCreator,
		shortCodeFollower: shortCodeFollower,
		analiticGetter:    analiticGetter,
	}
}

func (r *Router) Routes() {
	r.Router.POST("/horten")
	r.Router.GET("/s/:short_url")
	r.Router.GET("/analytics/short_url")
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
}
