package web

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
)

var (
	//go:embed "static/*"
	htmlStatic embed.FS

	//go:embed "template/index.html"
	indexTemplate []byte
)

func NewController(db *gorm.DB, r *chi.Mux) error {
	log.Info().Msg("create blog controller")
	wrap := func(f func(db *gorm.DB, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			f(db, w, r)
		}
	}
	r.Handle(opts.Media+"*", http.StripPrefix("/media/", convert_img.AutoThumb(opts.MediaPath, tmp, http.FileServer(http.Dir(opts.MediaPath)))))
	r.NotFound(wrap(notFound))
	r.Get("/", wrap(index))
	r.Handle("/static/*", http.FileServer(http.FS(htmlStatic)))
	return nil
}

func index(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	log.Info().Msg("render page")

	// Generate template
	result, err := axtools.Render(metricateTemplate, nil)

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}
	w.Write(result)
}
