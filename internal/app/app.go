package app

import (
	"log"

	"github.com/wb-go/wbf/config"
)

func Run() {
	cfg := config.New()
	err := cfg.LoadEnvFiles("../.env")
	if err != nil {
		log.Fatalf("[main] ошибка загрузки cfg %v", err)
	}
	cfg.EnableEnv("")
}
