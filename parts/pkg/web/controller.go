package web

import (
	"AbuzLandingChecker/parts/pkg/data"
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

	//go:embed "template/admin.html"
	adminTemplate []byte

	//go:embed "template/quiz.html"
	quizTemplate []byte
)

type AdminPage struct {
	Hash string
}

func NewController(db *gorm.DB, r *chi.Mux, c *data.UsersController) error {
	log.Info().Msg("create blog controller")
	wrap := func(f func(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			f(db, c, w, r)
		}
	}
	r.Get("/", wrap(index))
	r.Get("/abuzadmincreatehash/spDdryEuGzeNulzOBgZHekfOI", wrap(admin))
	r.Get("/quiz/{hashString}", wrap(quiz))
	r.Handle("/static/*", http.FileServer(http.FS(htmlStatic)))
	return nil
}

func index(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
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

func quiz(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	hashString := chi.URLParam(r, "hashString")
	ip := "122.3.3.2"
	UA := r.Header.Get("User-Agent")
	err := c.UpdateIP(ip, hashString, UA)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to write user")
		w.Write(([]byte)(err.Error()))
		return
	}

	log.Info().Msg("render page")

	// Generate template
	result, err := Render(quizTemplate, nil)
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}
	w.Write(result)
}

func admin(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	fmt.Println("HERE!")
	w.Header().Add("Content-Type", "text/html")

	log.Info().Msg("render page")
	generatedHash, err := c.CreateHash()
	dataPage := AdminPage{Hash: generatedHash}
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to create hash")
		w.Write(([]byte)(err.Error()))
		return
	}

	// Generate template
	result, err := Render(adminTemplate, dataPage)

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
