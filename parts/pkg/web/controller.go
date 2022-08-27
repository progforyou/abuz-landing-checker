package web

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"net/http"
	"text/template"
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
	r.Get("/", wrap(index))
	r.Handle("/static/*", http.FileServer(http.FS(htmlStatic)))
	return nil
}

func index(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	log.Info().Msg("render page")

	// Generate template
	result, err := Render(indexTemplate, nil)

	//ip := "122.3.3.3"

	//check if exist
	//var res data.Users
	//tx := db.Model(&data.Users{}).Find(&res, "ip", ip)
	/*if tx.RowsAffected == 0 {
		return nil
	}*/

	fmt.Println("r.Header", r.Header.Get("User-Agent"))
	//rowData := data.Users{FP: r.Header.Get("User-Agent")}

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}
	w.Write(result)
}

func Render(templateByte []byte, data interface{}) ([]byte, error) {
	t, err := template.New("").Parse(string(templateByte))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create template")
		return nil, err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to render template")
		return nil, err
	}
	return tpl.Bytes(), nil
}
