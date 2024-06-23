package main

import (
	"net/http"

	db "github.com/hanifmaliki/chat-app/internal/db"
	route "github.com/hanifmaliki/chat-app/internal/route"
	util "github.com/hanifmaliki/chat-app/pkg/util"

	"github.com/rs/zerolog/log"
)

func main() {
	util.LoadConfig()

	db, err := db.NewGormDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing storage")
	}

	mux := http.NewServeMux()

	route.SetupRoutes(mux, db)

	port := util.GetEnv("PORT", "8080")
	log.Info().Msgf("Server started on :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
