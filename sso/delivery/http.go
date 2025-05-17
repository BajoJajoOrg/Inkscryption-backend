package delivery

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/response"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/util"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/util/token"
	"github.com/BajoJajoOrg/Inkscryption-backend/sso"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Id           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

type RenewReq struct {
	RefreshToken string `json:"refresh_token"`
}

type handlers struct {
	cfg        *config.Config
	router     *chi.Mux
	userUC     sso.Usecase
	logger     *slog.Logger
	tokenMaker *token.JWTMaker
}

func New(cfg *config.Config, router *chi.Mux, userUC sso.Usecase, logger *slog.Logger,
	accessSecretKey string, refreshSecretKey string) sso.Handlers {
	return &handlers{
		cfg:        cfg,
		router:     router,
		userUC:     userUC,
		logger:     logger,
		tokenMaker: token.NewJWTMaker(accessSecretKey, refreshSecretKey),
	}
}

func (h *handlers) MapHandlers(tokenAuth *jwtauth.JWTAuth) error {
	h.router.Route("/register", func(r chi.Router) {
		r.Post("/", h.Register)
	})
	h.router.Route("/login", func(r chi.Router) {
		r.Post("/", h.Login)
	})

	h.router.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/logout", h.Logout)
		r.Post("/renew", h.RenewToken)
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

	accessToken, _, err := h.tokenMaker.CreateAccessToken(*gu.ID, gu.Email, 15*time.Minute)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	refreshToken, _, err := h.tokenMaker.CreateRefreshToken(*gu.ID, gu.Email, 24*time.Hour)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	userId := strconv.Itoa(*gu.ID)

	err = h.userUC.AddSession(context.TODO(), refreshToken, userId)
	if err != nil {
		h.logger.Error("error adding a session", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error adding a session"))
		return
	}

	response := LoginUserReq{
		Id:           userId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        u.Email,
	}

	render.JSON(w, r, response)
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

	// oldPass := u.Password

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

	accessToken, _, err := h.tokenMaker.CreateAccessToken(*created, u.Email, 15*time.Minute)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	refreshToken, _, err := h.tokenMaker.CreateRefreshToken(*created, u.Email, 24*time.Hour)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	userId := strconv.Itoa(*created)

	err = h.userUC.AddSession(context.TODO(), refreshToken, userId)
	if err != nil {
		h.logger.Error("error adding a session", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error adding a session"))
		return
	}

	response := LoginUserReq{
		Id:           strconv.Itoa(*created),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        u.Email,
	}

	render.JSON(w, r, response)
}

func (h *handlers) Logout(w http.ResponseWriter, r *http.Request) {
	// var u LogoutReq

	// err := render.DecodeJSON(r.Body, &u)
	// if err != nil {
	// 	h.logger.Error("failed to decode request body", slog.Attr{
	// 		Key:   "error",
	// 		Value: slog.StringValue(err.Error()),
	// 	})

	// 	w.WriteHeader(http.StatusBadRequest)
	// 	render.JSON(w, r, response.Error("failed to decode request"))
	// 	return
	// }

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	// println(strconv.FormatFloat(userID, 'f', -1, 64))

	err := h.userUC.DeleteSession(context.TODO(), strconv.FormatFloat(userID, 'f', -1, 64))
	if err != nil {
		h.logger.Error("failed to delete session", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to delete session"))
		return
	}

	w.WriteHeader(200)
}

func (h *handlers) RenewToken(w http.ResponseWriter, r *http.Request) {
	var u RenewReq

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

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "email not found", http.StatusUnauthorized)
		return
	}

	err = h.userUC.GetSession(context.TODO(), u.RefreshToken, strconv.FormatFloat(userID, 'f', -1, 64))
	if err != nil {
		h.logger.Error("wrong refresh token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("wrong refresh token"))
		return
	}

	accessToken, _, err := h.tokenMaker.CreateAccessToken(int(userID), email, 15*time.Minute)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	refreshToken, _, err := h.tokenMaker.CreateRefreshToken(int(userID), email, 24*time.Hour)
	if err != nil {
		h.logger.Error("error creating a token", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error creating a token"))
		return
	}

	err = h.userUC.AddSession(context.TODO(), refreshToken, strconv.FormatFloat(userID, 'f', -1, 64))
	if err != nil {
		h.logger.Error("error adding a session", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("error adding a session"))
		return
	}

	response := LoginUserReq{
		Id:           strconv.FormatFloat(userID, 'f', -1, 64),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Email:        email,
	}

	render.JSON(w, r, response)

}
