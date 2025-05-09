package delivery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/response"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/util"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/util/token"
	"github.com/BajoJajoOrg/Inkscryption-backend/sso"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type handlers struct {
	cfg        *config.Config
	router     *chi.Mux
	userUC     sso.Usecase
	logger     *slog.Logger
	tokenMaker *token.JWTMaker
}

func New(cfg *config.Config, router *chi.Mux, userUC sso.Usecase, logger *slog.Logger, secretKey string) sso.Handlers {
	return &handlers{
		cfg:        cfg,
		router:     router,
		userUC:     userUC,
		logger:     logger,
		tokenMaker: token.NewJWTMaker(secretKey),
	}
}

func (h *handlers) MapHandlers() error {
	h.router.Route("/register", func(r chi.Router) {
		r.Post("/", h.Register)
	})
	h.router.Route("/login", func(r chi.Router) {
		r.Post("/", h.Login)
	})

	return nil
}

func (h *handlers) ListenAndServe(cfg config.HTTPServer) error {
	address := ":" + cfg.SSOPort
	err := http.ListenAndServe(address, h.router)
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}
	return nil
}

func (h *handlers) Login(w http.ResponseWriter, r *http.Request) {
	var u LoginUserReq

	err := render.DecodeJSON(r.Body, &u)
	if err != nil {
		h.logger.Error("failed to decode request body", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	gu, err := h.userUC.Get(context.TODO(), u.Email)
	if err != nil {
		h.logger.Error("failed to get user", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get user"))
		return
	}

	err = util.CheckPassword(u.Password, gu.Password)
	if err != nil {
		h.logger.Error("wrong password", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, response.Error("wrong password"))
		return
	}

	accessToken, _, err := h.tokenMaker.CreateToken(*gu.ID, gu.Email, 15*time.Hour)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	u.AccessToken = accessToken

	render.JSON(w, r, u)

	// create a token and return it as response
}

func (h *handlers) Register(w http.ResponseWriter, r *http.Request) {
	var u UserReq

	err := render.DecodeJSON(r.Body, &u)
	if err != nil {
		h.logger.Error("failed to decode request body", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	if u.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("empty email"))
		return
	}

	hashed, err := util.HashPassword(u.Password)
	if err != nil {
		h.logger.Error("failed to hash password", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to hash password"))
		return
	}

	oldPass := u.Password

	u.Password = hashed

	user := sso.User{
		Email:    u.Email,
		Password: u.Password,
	}

	created, err := h.userUC.Create(context.TODO(), user)
	if err != nil {

		if strings.Contains(err.Error(), "email already exists") {
			w.WriteHeader(http.StatusConflict)
			render.JSON(w, r, response.Error("email already exists"))
			return
		}

		h.logger.Error("failed to register user", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to register user"))
		return
	}

	accessToken, _, err := h.tokenMaker.CreateToken(*created, u.Email, 15*time.Hour)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	response := LoginUserReq{
		AccessToken: accessToken,
		Email:       u.Email,
		Password:    oldPass,
	}

	render.JSON(w, r, response)
}
