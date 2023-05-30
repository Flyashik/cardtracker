package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"server/internal/app/config"
	"server/internal/app/helper"
	"server/internal/app/middlewares"
	"server/internal/app/models"
	"server/internal/app/storage"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	api.HandleFunc("/devices", middlewares.IsAuthorized(s.handleDevices()))
	api.HandleFunc("/user", middlewares.IsAuthorized(s.handleUser()))
	api.HandleFunc("/users", middlewares.IsAuthorized(s.handleUsers()))
	api.HandleFunc("/notifications", middlewares.IsAuthorized(s.handleNotifications()))
	api.HandleFunc("/login", s.handleLogin())
	api.HandleFunc("/logout", s.handleLogout())
	api.HandleFunc("/register", s.handleRegister())
	api.HandleFunc("/new_notification", s.handleNewNotification())

	fs := http.FileServer(http.Dir("./static/dist"))

	s.router.PathPrefix("/api").Handler(api)
	s.router.PathPrefix("/").HandlerFunc(staticHandler(fs))

	s.router.NotFoundHandler = notFoundHandler(fs)

	s.router.Use(corsMiddleware)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, X-HTTP-Method-Override, Content-Type, Accept, Set-Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")

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

func staticHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ext := filepath.Ext(path)

		if ext == "" {
			http.ServeFile(w, r, "./static/dist/index.html")
			return
		}

		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			w.Header().Set("Content-Type", mimeType)
		}

		fs.ServeHTTP(w, r)
	}
}

func notFoundHandler(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
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
			AuthID  int              `json:"authorization_id"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Info(`[Phone info] Error when reading request body`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.Body.Close()

		var resp *Info
		err = json.Unmarshal(body, &resp)
		if err != nil {
			s.logger.Info(`[Phone info] Error when unmarshalling request body`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp.Phone.ModelTag, err = helper.ConvertModelTagToMarketingName(resp.Phone.ModelTag)
		if err != nil {
			s.logger.Info(`[Phone info] Error when translating model tag`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp.Phone.SimSlots = len(resp.SimInfo)
		resp.Phone.SdSlots = len(resp.SdInfo)

		phone, err := s.storage.Phone().Create(&resp.Phone)
		if err != nil {
			s.logger.Info(`[Phone info] Error when creating phone`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		s.storage.Sim().RemovePhoneId(phone.Id)
		for _, sim := range resp.SimInfo {
			_, err := s.storage.Sim().Create(&sim, phone)
			if err != nil {
				s.logger.Info(`[Phone info] Error while creating sim`)
				s.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		s.storage.SdCard().RemovePhoneId(phone.Id)
		for _, sd := range resp.SdInfo {
			sd.SdManufacturerId, err = helper.ConvertManufacturerIdToCompanyName(sd.SdManufacturerId)
			if err != nil {
				s.logger.Info(`[Phone info] Error while translating sd info`)
				s.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err := s.storage.SdCard().Create(&sd, phone)
			if err != nil {
				s.logger.Info(`[Phone info] Error while creating sd card`)
				s.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		user, err := s.storage.User().SelectByCode(resp.AuthID)
		if err != nil {
			s.logger.Info(`[Phone info] Error while finding user by code`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = s.storage.UserPhone().CreateRelation(user.Id, phone.Id)
		if err != nil {
			s.logger.Info(`[Phone info] Error while creating relation`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		query := r.URL.Query().Get("user_info_needed")
		if query == "true" {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(user); err != nil {
				s.logger.Info(`[Phone info] Error while encoding json`)
				s.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		w.WriteHeader(http.StatusOK)
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
			s.logger.Info(`[Devices info] Error while fetching phones`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusInternalServerError))
			http.Error(w, "Failed fetch phones", http.StatusInternalServerError)
			return
		}
		simCards, err := s.storage.Sim().SelectAll()
		if err != nil {
			s.logger.Info(`[Devices info] Error while fetching sim cards`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusInternalServerError))
			http.Error(w, "Failed fetch simcards", http.StatusInternalServerError)
			return
		}
		sdCards, err := s.storage.SdCard().SelectAll()
		if err != nil {
			s.logger.Info(`[Devices info] Error while fetching sd cards`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusInternalServerError))
			http.Error(w, "Failed fetch sdcards", http.StatusInternalServerError)
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

func (s *Server) handleNotifications() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		modelNumber := r.URL.Query().Get("model_number")
		if modelNumber == "" {
			s.logger.Info(`[Notifications] There was no parameter in request`)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		notificationList, err := s.storage.Notification().SelectByModelTag(modelNumber)
		if err != nil {
			s.logger.Info(`[Notifications] Error while fetching notifications by tag`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusNotFound))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(notificationList); err != nil {
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

// TODO: заменить на что-то адекватное
var jwtKey = []byte("very_secret_key")

func (s *Server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var existingUser *models.User

		existingUser, err := s.storage.User().SelectByEmail(user.Email)

		// TODO check user email
		if err != nil {
			s.logger.Info(`[Login] Error while fetching user by email`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusBadRequest))
			http.Error(w, "User does not exist", http.StatusBadRequest)
			return
		}

		errHash := helper.CompareHashPassword(user.Password, existingUser.Password)
		if !errHash {
			s.logger.Info(`[Login] Error while testing password`)
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)

		claims := &models.Claims{
			Role: existingUser.Role,
			StandardClaims: jwt.StandardClaims{
				Subject:   existingUser.Email,
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			s.logger.Info(`[Login] Error while generating jwt`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusInternalServerError))
			http.Error(w, "could not generate token", http.StatusInternalServerError)
			return
		}

		cookie := &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
			Path:    "/",
		}
		http.SetCookie(w, cookie)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

func (s *Server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:    "token",
			Value:   "",
			Expires: time.Unix(0, 0),
			Path:    "/",
		}
		http.SetCookie(w, cookie)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

func (s *Server) handleRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			s.logger.Info(`[Register] Error while decoding json`)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user.Role = "user"
		_, err := s.storage.User().SelectByEmail(user.Email)
		if err == nil {
			s.logger.Info(`[Register] Error while checking for user existance`)
			http.Error(w, "user already exists", http.StatusBadRequest)
			return
		}

		var errHash error
		user.Password, errHash = helper.GenerateHashPassword(user.Password)
		if errHash != nil {
			s.logger.Info(`[Register] Error while generating passwordl`)
			http.Error(w, "could not generate password hash", http.StatusInternalServerError)
			return
		}

		random := rand.New(rand.NewSource(time.Now().UnixNano()))

		var userCode int
		for {
			userCode = random.Intn(99999-10000) + 10000
			if !s.storage.User().CheckCodeExists(userCode) {
				break
			}
		}
		user.Code = userCode
		_, err = s.storage.User().Create(&user)
		if err != nil {
			s.logger.Info(`[Register] Error while creating user`)
			s.logger.Error(fmt.Sprintf(`%s, %d`, err, http.StatusInternalServerError))
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

func (s *Server) handleUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			s.logger.Info(`[User] Error while checking cookie`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusUnauthorized))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := helper.ParseToken(cookie.Value)

		if err != nil {
			s.logger.Info(`[User] Error while checking jwt`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusUnauthorized))
			http.Error(w, "Incorrect token", http.StatusUnauthorized)
			return
		}
		u, err := s.storage.User().SelectByEmail(claims.StandardClaims.Subject)
		if err != nil {
			s.logger.Info(claims.StandardClaims.Subject)
			s.logger.Info(`[User] Error while checking user by email`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusNotFound))
			http.Error(w, "Can't fetch user", http.StatusNotFound)
		}
		u.Password = ""
		u.Role = ""
		json.NewEncoder(w).Encode(u)
		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

func (s *Server) handleNewNotification() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Info(`[NewNotification] Error while reading body`)
			s.logger.Error(err)
			return
		}
		r.Body.Close()

		var resp *models.Notification
		err = json.Unmarshal(body, &resp)
		if err != nil {
			s.logger.Info(`[NewNotification] Error while decoding json`)
			s.logger.Error(err)
			return
		}

		_, err = s.storage.Notification().Create(resp)
		if err != nil {
			s.logger.Info(`[NewNotification] Error while creating notification`)
			s.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		s.logger.Info(fmt.Sprintf(`%s %s%s %d`, r.Method, r.Host, r.RequestURI, http.StatusOK))
	}
}

func (s *Server) handleUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.storage.UserPhone().SelectUsersWithPhones()
		if err != nil {
			s.logger.Info(`[Users] Error while selecting users with phones`)
			s.logger.Error(fmt.Sprintf(`%s %d`, err, http.StatusNotFound))
			w.WriteHeader(http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
