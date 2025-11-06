package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
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
	GetAnalytics(ctx context.Context, shortURLID int) (map[string]map[string]int, error)
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
	r.Router.GET("/analytics/short_url_id", r.getAnalyticsHandler)
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

// followShortURLHandler - хэндлер для перехода по короткой ссылке и лог клика
func (r *Router) followShortURLHandler(c *gin.Context) {
	shortCode := c.Param("short_url")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
		return
	}
	userAgent := c.Request.UserAgent()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()

	originalURL, err := r.shortCodeFollower.ReturnOriginalURLByShort(ctx, shortCode, userAgent)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "short URL not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

// getAnalyticsHandler - хэндлер получения аналитики
func (r *Router) getAnalyticsHandler(c *gin.Context) {
	idStr := c.Param("short_url_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid short_url_id"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	stats, err := r.analiticGetter.GetAnalytics(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get analytics: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
