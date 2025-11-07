package app

import (
	"log"
	"net/http"
	"time"

	"github.com/Vladimirmoscow84/Shortener.git/internal/handlers"
	"github.com/Vladimirmoscow84/Shortener.git/internal/service"
	"github.com/Vladimirmoscow84/Shortener.git/internal/storage"
	"github.com/Vladimirmoscow84/Shortener.git/internal/storage/cache"
	"github.com/Vladimirmoscow84/Shortener.git/internal/storage/postgres"
	"github.com/wb-go/wbf/config"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/redis"
)

func Run() {
	cfg := config.New()
	err := cfg.LoadEnvFiles(".env")
	if err != nil {
		log.Fatalf("[main] ошибка загрузки cfg %v", err)
	}
	cfg.EnableEnv("")

	databaseURI := cfg.GetString("DATABASE_URI")
	serverAddr := cfg.GetString("SERVER_ADDRESS")
	redisAddr := cfg.GetString("REDIS_URI")

	postgresStore, err := postgres.New(databaseURI)
	if err != nil {
		log.Fatalf("[app] failed to connect to PG DB: %v", err)
	}
	defer postgresStore.Close()

	rd := redis.New(redisAddr, "", 0)
	redisStore := cache.NewCache(rd)
	if redisStore == nil {
		log.Fatalf("[app] failed to create Redis client")

	} else {
		log.Println("[main] redis successfuly connected")
	}

	store, err := storage.New(postgresStore, rd)
	if err != nil {
		log.Fatalf("[app] failed to init unified storage: %v", err)
	}
	log.Println("[app] Unified storage initialized successfully")

	serviceURL := service.New(store, redisStore)

	engine := ginext.New("release")
	router := handlers.New(engine, serviceURL, serviceURL, serviceURL)
	router.Routes()

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("[app] starting server at %s", serverAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("[app] server error: %v", err)
	}
}
