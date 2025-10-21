package app

import (
	"log"

	"github.com/wb-go/wbf/config"
)

func Run() {
	cfg := config.New()
	err := cfg.Load("../config.yaml", "../.env", "")
	if err != nil {
		log.Fatalf("[main]load cfg dissable %v", err)
	}
}
