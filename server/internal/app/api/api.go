package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"path/filepath"
	"server/internal/app/config"
	"server/internal/app/helper"
	"server/internal/app/models"
	"server/internal/app/storage"
)

type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *config.Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStorage(); err != nil {
		return err
	}

	s.logger.Info("Starting server...")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *Server) configureRouter() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/test", s.handleTest())
	api.HandleFunc("/phone_info", s.handlePhoneInfo())
	api.HandleFunc("/devices", s.handleDevices())

	mimeTypes := map[string]string{
		".css":  "text/css",
		".js":   "application/javascript",
		".json": "application/json",
	}

	fs := http.FileServer(http.Dir("./static/dist"))

	fs = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ext := filepath.Ext(path)
		if mimeType, ok := mimeTypes[ext]; ok {
			w.Header().Set("Content-Type", mimeType)
		}
		fs.ServeHTTP(w, r)
	})

	s.router.PathPrefix("/").Handler(indexHandler())
	s.router.NotFoundHandler = notFoundHandler()

	s.router.Use(corsMiddleware)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) configureStorage() error {
	st := storage.New(s.config.Storage)
	if err := st.Open(); err != nil {
		return err
	}

	s.storage = st

	return nil
}

func indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dist/index.html")
	}
}

func notFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dist/index.html")
	}
}

func (s *Server) handleTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Just test")
	}
}

func (s *Server) handlePhoneInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Info struct {
			Phone   models.Phone     `json:"phone_info"`
			SimInfo []models.SimInfo `json:"sim_info"`
			SdInfo  []models.SdInfo  `json:"sd_info"`
			AuthID  uint             `json:"authorization_id"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}
		r.Body.Close()

		var resp *Info
		err = json.Unmarshal(body, &resp)
		if err != nil {
			s.logger.Error(err)
			return
		}

		resp.Phone.ModelTag, err = helper.ConvertModelTagToMarketingName(resp.Phone.ModelTag)
		if err != nil {
			s.logger.Error(err)
			return
		}
		resp.Phone.SimSlots = len(resp.SimInfo)
		resp.Phone.SdSlots = len(resp.SdInfo)

		phone, err := s.storage.Phone().Create(&resp.Phone)
		if err != nil {
			s.logger.Error(err)
			return
		}

		s.storage.Sim().RemovePhoneId(phone.Id)
		for _, sim := range resp.SimInfo {
			_, err := s.storage.Sim().Create(&sim, phone)
			if err != nil {
				s.logger.Error(err)
				return
			}
		}

		s.storage.SdCard().RemovePhoneId(phone.Id)
		for _, sd := range resp.SdInfo {
			sd.SdManufacturerId, err = helper.ConvertManufacturerIdToCompanyName(sd.SdManufacturerId)
			if err != nil {
				s.logger.Error(err)
				return
			}
			_, err := s.storage.SdCard().Create(&sd, phone)
			if err != nil {
				s.logger.Error(err)
				return
			}
		}

		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

func (s *Server) handleDevices() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Phones   []models.Phone   `json:"phones"`
			SimCards []models.SimInfo `json:"simCards"`
			SdCards  []models.SdInfo  `json:"sdCards"`
		}

		phones, err := s.storage.Phone().SelectAll()
		if err != nil {
			s.logger.Error(err)
			return
		}
		simCards, err := s.storage.Sim().SelectAll()
		if err != nil {
			s.logger.Error(err)
			return
		}
		sdCards, err := s.storage.SdCard().SelectAll()
		if err != nil {
			s.logger.Error(err)
			return
		}

		response := Response{
			Phones:   phones,
			SimCards: simCards,
			SdCards:  sdCards,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
