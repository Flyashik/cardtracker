package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"server/internal/app/config"
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
	s.router.HandleFunc("/test", s.handleTest())
	//handlers ...
	s.router.HandleFunc("/phone_info", s.handlePhoneInfo())
}

func (s *Server) configureStorage() error {
	st := storage.New(s.config.Storage)
	if err := st.Open(); err != nil {
		return err
	}

	s.storage = st

	return nil
}

func (s *Server) handleTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Just test")
	}
}

func (s *Server) handlePhoneInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			resp *models.Info
		)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}

		err = json.Unmarshal(body, &resp)
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
			_, err := s.storage.SdCard().Create(&sd, phone)
			if err != nil {
				s.logger.Error(err)
				return
			}
		}

		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}
