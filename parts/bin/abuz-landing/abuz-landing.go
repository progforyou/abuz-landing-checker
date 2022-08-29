package main

import (
	"AbuzLandingChecker/parts/pkg/data"
	"AbuzLandingChecker/parts/pkg/web"
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os"
)

// -build-me-for: native
// -build-me-for: linux

var (
	port int
)

func init() {
	flag.IntVar(&port, "port", 80, "set port")
	flag.Parse()
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.DebugLevel)
	log.Debug().Msgf("Start Telegram Blog server on port %d", port)

	db, err := gorm.Open(sqlite.Open("file:db.sqlite3?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("fail to open database")
	}
	err = db.AutoMigrate(&data.Users{})

	if err != nil {
		log.Fatal().Err(err).Msg("fail to migrate database")
	}
	httpLogger := log.With().Str("service", "http").Logger().Level(zerolog.InfoLevel)
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(httpLogger))

	c := data.NewUsersController(db, httpLogger)
	err = web.NewController(db, r, &c)
	if err != nil {
		log.Fatal().Err(err).Msg("fail create blog")
	}

	log.Debug().Msgf("start server on port %d", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		log.Fatal().Err(err).Msg("fail start server")
	}
}
