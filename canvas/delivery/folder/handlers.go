package folder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/interfaces/folder"
	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/response"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type Request struct {
	Name   string `json:"name"`
	Parent int    `json:"parent"`
}

type Response struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type handlers struct {
	cfg      *config.Config
	router   *chi.Mux
	folderUC folder.UseCase
	logger   *slog.Logger
}

func New(cfg *config.Config, router *chi.Mux, folderUC folder.UseCase, logger *slog.Logger) folder.Handlers {
	return &handlers{
		cfg:      cfg,
		router:   router,
		folderUC: folderUC,
		logger:   logger,
	}
}

func (h *handlers) MapHandlers() error {
	h.router.Route("/folder", func(r chi.Router) {
		// r.Post("/", h.Create)
		// r.Get("/", h.Get)
		r.Post("/", h.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.Get)
			r.Delete("/", h.Delete)
			// r.Put("/", h.Update)
		})
	})

	// ml/image-to-text

	return nil
}

func (h *handlers) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	newId, _ := strconv.Atoi(id)

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	filterOptions := filter.NewOptions()
	// TODO: убрать хардкод
	name := r.URL.Query().Get("name")
	if name != "" {
		err := filterOptions.AddField("name", filter.OperatorLike, name, filter.DataTypeStr)
		if err != nil {
			h.logger.Error("failed to parse query", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("cannot add name field into filter"))
			return
		}
	}

	created_at := r.URL.Query().Get("created_at")
	if created_at != "" {
		if !util.ValidateDates(created_at) {
			h.logger.Error("wrong dates format")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ProError(400, "Invalid request parameters",
				response.Details{
					Field: "created_at",
					Error: "Invalide date format",
				}))
			return
		}
	}

	if created_at != "" {
		var operator string
		if strings.Contains(created_at, ":") {
			operator = filter.OperatorBetween
		} else {
			operator = filter.OperatorEq
		}
		err := filterOptions.AddField("created_at", operator, created_at, filter.DataTypeDate)
		if err != nil {
			h.logger.Error("failed to parse query", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("cannot add created_at field into filter"))
			return
		}
		fmt.Println(filterOptions.GetField("created_at"))
	}

	folderContent, err := h.folderUC.Get(context.TODO(), filterOptions, newId, int(userID))
	if err != nil {
		h.logger.Error("failed to get folders", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get all folders"))
		return
	}

	render.JSON(w, r, folderContent)
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {
	var req Request

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		h.logger.Error("failed to decode request body", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode request"))
		return
	}

	h.logger.Info("request body decoded", slog.Any("request", req))

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	if req.Name == "" { // TODO проверка айдишника
		h.logger.Error("folder name is empty", slog.Attr{
			Key:   "error",
			Value: slog.StringValue("folder name cannot be empty"),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("folder name cannot be empty"))
		return
	}

	folder := canvas.FolderBase{
		Name:      req.Name,
		Parent:    req.Parent,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	id, err := h.folderUC.Create(context.TODO(), folder, int(userID))
	if err != nil {
		h.logger.Error("failed to create folder", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create folder"))
		return
	}

	response := Response{
		Id:        *id,
		Name:      folder.Name,
		UpdatedAt: folder.UpdatedAt,
	}

	render.JSON(w, r, response)
}

func (h *handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	newId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("Invalid request parameters", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.ProError(400, "Invalid request parameters",
			response.Details{
				Field: "id",
				Error: "Invalid id format",
			}))
		return
	}

	if err = h.folderUC.Delete(context.TODO(), newId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Error("no such canvas to be deleted", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.ProError(404, "Canvas with such id was not found",
				response.Details{
					Field: "id",
					Error: "This id does not exist",
				}))
			return

		} else {
			h.logger.Error("failed to delete canvas", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete canvas"))
			return
		}
	}

	w.WriteHeader(204)
}
