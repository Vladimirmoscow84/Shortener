package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Vladimirmoscow84/Shortener.git/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type shortCodeCreator interface {
	CreateShortURL(ctx context.Context, longURL string) (*model.ShortURL, error)
}

type shortCodeFollower interface {
	ReturnOriginalURLByShort(ctx context.Context, shortCode, userAgent string) (string, error)
}

type analiticGetter interface {
	GetAnalytics(ctx context.Context, shortURLID uint) (map[string]map[string]int, error)
	GetShortURL(ctx context.Context, shortCode string) (*model.ShortURL, error)
}

type Router struct {
	Router            *ginext.Engine
	shortCodeCreator  shortCodeCreator
	shortCodeFollower shortCodeFollower
	analiticGetter    analiticGetter
}

func New(router *ginext.Engine, shortCodeCreator shortCodeCreator, shortCodeFollower shortCodeFollower, analiticGetter analiticGetter) *Router {
	return &Router{
		Router:            router,
		shortCodeCreator:  shortCodeCreator,
		shortCodeFollower: shortCodeFollower,
		analiticGetter:    analiticGetter,
	}
}

func (r *Router) Routes() {
	r.Router.POST("/horten", r.createShortURLHandler)
	r.Router.GET("/s/:short_url", r.followShortURLHandler)
	r.Router.GET("/analytics/short_url", r.getAnalyticsHandler)
	r.Router.GET("/", func(c *gin.Context) { c.File("./web/index.html") })
	r.Router.Static("/static", "./web")
}

// createShortURLHandler - хэндлер создания короткой ссылки
func (r *Router) createShortURLHandler(c *gin.Context) {

	var request struct {
		OriginalCode string `json:"original_code" binding:"requared,url"`
	}
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		log.Println("[handler] invalid request to create short URL")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	short, err := r.shortCodeCreator.CreateShortURL(ctx, request.OriginalCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create short url: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":            short.ID,
		"short_code":    short.ShortCode,
		"original_code": short.OriginalCode,
		"created_at":    short.CreatedAt,
	})

}

// followShortURLHandler - хэндлер для перехода по короткой ссылке
func (r *Router) followShortURLHandler(c *gin.Context) {
	shortCode := c.Param("short_code")
	userAgent := c.Request.UserAgent()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	originalURL, err := r.shortCodeFollower.ReturnOriginalURLByShort(ctx, shortCode, userAgent)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// getAnalyticsHandler - хэндлер получения аналитики
func (r *Router) getAnalyticsHandler(c *gin.Context) {
	code := c.Param("short_code")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	short, err := r.analiticGetter.GetShortURL(ctx, code)
	if err != nil || short == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	data, err := r.analiticGetter.GetAnalytics(ctx, uint(short.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get analytics: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code": short.ShortCode,
		"analytics":  data,
	})
}
