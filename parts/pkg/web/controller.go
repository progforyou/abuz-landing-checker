package web

import (
	"AbuzLandingChecker/parts/pkg/data"
	"bytes"
	"embed"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"text/template"
)

var (
	//go:embed "static/*"
	htmlStatic embed.FS

	//go:embed "template/index.html"
	indexTemplate []byte

	//go:embed "template/create.html"
	createTemplate []byte

	//go:embed "template/table.html"
	tableTemplate []byte

	//go:embed "template/sign-in.html"
	signInTemplate []byte

	//go:embed "template/edit.html"
	editTemplate []byte

	//go:embed "template/quiz.html"
	quizTemplate []byte
)

type CreatePage struct {
	Hash string
}

type AdminPage struct {
	Data []data.Users
}

type SignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

func NewController(db *gorm.DB, r *chi.Mux, c *data.UsersController) error {
	log.Info().Msg("create web controller")
	wrap := func(f func(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			f(db, c, w, r)
		}
	}
	r.Get("/", wrap(index))
	r.Post("/abuzadmin/api/signin", func(w http.ResponseWriter, r *http.Request) {
		var obj SignIn
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("read body")
			w.WriteHeader(403)
			return
		}
		if err := json.Unmarshal(body, &obj); err != nil {
			log.Error().Err(err).Msg("decode json")
			return
		}

		if obj.Login == "admin" && obj.Password == "abuzadmin" {
			rowObj := data.Admin{
				IP:     obj.IP,
				SignIn: true,
			}
			tx := db.Model(&data.Admin{}).Create(&rowObj)
			if tx.Error != nil {
				log.Error().Err(tx.Error).Msg("db error")
				w.WriteHeader(403)
			}
			w.WriteHeader(202)
			w.Write([]byte("OK"))
			return
		} else {
			w.WriteHeader(202)
			w.Write([]byte("NO"))
			return
		}
	})
	r.Get("/abuzadmin/signin", wrap(signin))
	r.Get("/abuzadmin/table", wrap(table))
	r.Get("/abuzadmin/hash", wrap(createHash))
	r.Get("/abuzadmin/edit/{id}", wrap(edit))
	r.Post("/abuzadmin/api/edit", func(w http.ResponseWriter, r *http.Request) {
		var obj data.Users
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error().Err(err).Msg("read body")
			w.WriteHeader(403)
			return
		}
		if err := json.Unmarshal(body, &obj); err != nil {
			log.Error().Err(err).Msg("decode json")
			return
		}
		if err := c.Save(obj); err != nil {
			log.Error().Err(err).Msg("error to save")
			return
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
		return
	})
	r.Get("/quiz/{hashString}", wrap(quiz))
	r.Handle("/static/*", http.FileServer(http.FS(htmlStatic)))
	return nil
}

func index(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	realIP := r.Header.Get("X-Real-IP")
	readAddress := r.Header.Get("X-Real-IP-Address")
	log.Info().Str("Real IP", realIP).Str("readAddress", readAddress).Msg("render page")

	// Generate template
	result, err := Render(indexTemplate, nil)

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
	ip := r.Header.Get("X-Real-IP")
	UA := r.Header.Get("User-Agent")
	resp, err := http.Get("http: //ip-api.com/json/" + ip)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(500)
		w.Write(([]byte)(err.Error()))
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(500)
		w.Write(([]byte)(err.Error()))
		return
	}
	var bodyJson map[string]interface{}
	if err := json.Unmarshal(body, &bodyJson); err != nil {
		log.Error().Err(err).Msg("decode json")
		w.WriteHeader(500)
		return
	}
	var city string
	if bodyJson["city"] != nil {
		city = bodyJson["city"].(string)
	}
	err = c.UpdateIP(ip, hashString, UA, city)
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

func createHash(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	checkAdmin(w, r, db)

	log.Info().Msg("render page")
	generatedHash, err := c.CreateHash()
	dataPage := CreatePage{Hash: generatedHash}
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to create hash")
		w.Write(([]byte)(err.Error()))
		return
	}

	// Generate template
	result, err := Render(createTemplate, dataPage)

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}
	w.Write(result)
}

func table(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	checkAdmin(w, r, db)

	log.Info().Msg("render page")
	allUsers, err := c.GetAll()

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to create hash")
		w.Write(([]byte)(err.Error()))
		return
	}

	dataPage := AdminPage{Data: allUsers}

	// Generate template
	result, err := Render(tableTemplate, dataPage)

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}
	w.Write(result)
}

func signin(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	checkAdminR(w, r, db)

	log.Info().Msg("render page")

	// Generate template
	result, err := Render(signInTemplate, nil)

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}

	w.Write(result)
}

func edit(db *gorm.DB, c *data.UsersController, w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	checkAdmin(w, r, db)

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to parse id")
		w.Write(([]byte)(err.Error()))
		return
	}
	dataEdit, err := c.GetById(id)

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to get id")
		w.Write(([]byte)(err.Error()))
		return
	}

	// Generate template
	result, err := Render(editTemplate, dataEdit)

	if err != nil {
		w.WriteHeader(500)
		log.Error().Err(err).Msg("fail to render")
		w.Write(([]byte)(err.Error()))
		return
	}
	w.Write(result)
}

func checkAdmin(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var obj data.Admin
	ip := r.Header.Get("X-Real-IP")
	//#TODO FIX
	ip = "46.61.42.11"
	tx := db.Model(&data.Admin{}).Where("ip = ?", ip).Find(&obj)
	if tx.RowsAffected == 0 {
		http.Redirect(w, r, "/abuzadmin/signin", 302)
	}
}

func checkAdminR(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var obj data.Admin
	ip := r.Header.Get("X-Real-IP")
	//#TODO FIX
	ip = "46.61.42.11"
	tx := db.Model(&data.Admin{}).Where("ip = ?", ip).Find(&obj)
	if tx.RowsAffected > 0 {
		http.Redirect(w, r, "/abuzadmin/table", 302)
	}
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
