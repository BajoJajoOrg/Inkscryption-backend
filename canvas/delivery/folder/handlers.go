package folder

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/interfaces/folder"
	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/response"
	"github.com/go-chi/chi/v5"
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

	folderContent, err := h.folderUC.Get(context.TODO(), newId, 1)
	if err != nil {
		h.logger.Error("failed to get canvases", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get all canvases"))
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

	id, err := h.folderUC.Create(context.TODO(), folder, 1)
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
