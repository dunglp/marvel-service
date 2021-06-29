package service

import (
	"context"
	"net/http"
	"os"
	"sync"
	"time"
	"xendit-technical-assessment/pkg/config"
	"xendit-technical-assessment/pkg/domain"
	xenditHttp "xendit-technical-assessment/pkg/http"
	"xendit-technical-assessment/pkg/marvel"

	"github.com/patrickmn/go-cache"
)

func Run() Service {
	log := newLogger(os.Stdout)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to load config")
	}

	setLogLevel(cfg.LogLevel)
	setLogLevelFieldName(cfg.LogLevelFieldName)

	log.Info().Msg("Starting marvel service")
	cachingRuntime := cache.New(24*time.Hour, 24*time.Hour)
	cred, err := domain.NewCredentials(cfg.PublicKey, cfg.PrivateKey, cfg.Ts)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to init credentials")
	}

	marvelSvc := marvel.NewService(cfg.Hostname, cred, cachingRuntime)
	httpServer := xenditHttp.NewMarvelServiceHTTPServer(cfg.ServerAddress, marvelSvc)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Msgf("Listening on: %s", cfg.ServerAddress)
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Listen and serve")
		}
	}()

	return stopFunc(func() {
		log.Info().Msg("Stopping marvel service")

		log.Debug().Msg("Shutdown HTTP cmd")
		if err = httpServer.Shutdown(context.Background()); err != nil {
			log.Err(err).Msg("Failed to shutdown HTTP cmd")
		}

		wg.Wait()
	})
}
