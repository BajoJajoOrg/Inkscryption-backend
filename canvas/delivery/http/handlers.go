package delivery

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/config"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/response"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type Request struct {
	Name     string `json:"name"`
	FolderId int    `json:"folder_id"`
}

type MLRequest struct {
	Text string `json:"text"`
}

type MLResponse struct {
	Text string `json:"text"`
}

type Response struct {
	Id         int       `json:"id"`
	CanvasName string    `json:"canvas_name"`
	FolderId   int       `json:"folder_id"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type handlers struct {
	cfg      *config.Config
	router   *chi.Mux
	canvasUC canvas.UseCase
	logger   *slog.Logger
}

func New(cfg *config.Config, router *chi.Mux, canvasUC canvas.UseCase, logger *slog.Logger) canvas.Handlers {
	return &handlers{
		cfg:      cfg,
		router:   router,
		canvasUC: canvasUC,
		logger:   logger,
	}
}

func (h *handlers) ListenAndServe(cfg config.HTTPServer) error {
	address := ":" + cfg.Port
	err := http.ListenAndServe(address, h.router)
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}
	return nil
}

func (h *handlers) MapHandlers() error {
	h.router.Route("/canvas", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.GetAll)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetByID)
			r.Delete("/", h.Delete)
			r.Put("/", h.Update)
		})
	})
	h.router.Route("/ml", func(r chi.Router) {
		r.Post("/image-to-text", h.ImageToText)
		r.Post("/text-to-image", h.TextToImage)
	})

	h.router.Route("/sound-predict", func(r chi.Router) {
		r.Post("/", h.SoundPredict)
	})

	// ml/image-to-text

	return nil
}

func (h *handlers) GetAll(w http.ResponseWriter, r *http.Request) {

	folder_id := r.Form.Get("folder_id")

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

	new_folder_id, err := strconv.Atoi(folder_id)
	if err != nil {
		h.logger.Error("failed to get canvases", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get all canvases"))
		return
	}

	canvases, err := h.canvasUC.GetAll(context.TODO(), filterOptions, new_folder_id, int(userID))
	if err != nil {
		h.logger.Error("failed to get canvases", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get all canvases"))
		return
	}

	render.JSON(w, r, canvases)
}

func (h *handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	// id := r.Context().Value("id").(int)

	// id := r.URL.Query().Get("id")

	// TODO: полнейшая хуита, нужна или мидлвара или че то еще
	id := chi.URLParam(r, "id")

	// TODO: хуитта
	newId, _ := strconv.Atoi(id)

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	canvasFound, err := h.canvasUC.GetByID(context.TODO(), newId, int(userID))
	if err != nil {
		h.logger.Error("failed to get id", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to get id"))
		return
	}

	// canvasFound := canvas.CanvasBase{
	// 	CanvasID:  &newId,
	// 	Name:      canvasFound.Name,
	// 	UpdatedAt: canvasFound.UpdatedAt,
	// 	Url:       canvasFound.Url,
	// }

	if canvasFound.Url != "" {
		file, err := http.Get(canvasFound.Url)
		if err != nil {
			h.logger.Error("failed to download file from s3", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ProError(500, "failed to download file from s3",
				response.Details{
					Field: "aws",
					Error: "failed to download",
				}))
		}
		if file.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.ProError(500, "failed to download file from s3",
				response.Details{
					Field: "aws",
					Error: "failed to download",
				}))
		}
		defer file.Body.Close()

		fileContent, err := io.ReadAll(file.Body)
		if err != nil {
			http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		}

		decoded, _, err := transform.Bytes(charmap.Windows1251.NewDecoder(), fileContent)
		if err != nil {
			// обработка ошибки
		}

		encodedFile := base64.StdEncoding.EncodeToString(decoded)

		canvasFound.Data = encodedFile

		// decoder := charmap.Windows1251.NewDecoder()
		// utf8Content, err := decoder.Bytes(fileContent)
		// if err != nil {
		// 	http.Error(w, "Failed to decode file content", http.StatusInternalServerError)
		// 	return
		// }

		// encodedFile := base64.StdEncoding.EncodeToString(utf8Content)
		// canvasFound.Data = encodedFile
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(canvasFound); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

	// w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	// w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
	// w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%s", filepath.Base(canvas.Url)))

	//render.JSON(w, r, canvas)
}

func (h *handlers) Create(w http.ResponseWriter, r *http.Request) {
	var req Request

	// TODO: вынести ошибки в отдельный модуль
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

	if req.Name == "" {
		h.logger.Error("canvas_name is empty", slog.Attr{
			Key:   "error",
			Value: slog.StringValue("canvas_name cannot be empty"),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("canvas name cannot be empty"))
		return
	}

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	// url := h.cfg.AWSConfig.SecretEndpoint + "/1/"

	canvas := canvas.CanvasBase{
		Name:      req.Name,
		UpdatedAt: time.Now(),
		FolderId:  req.FolderId,
		UserId:    int(userID),
		CreatedAt: time.Now(),
	}

	id, err := h.canvasUC.Create(context.TODO(), canvas)
	if err != nil {
		h.logger.Error("failed to create canvas", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create canvas"))
		return
	}

	response := Response{
		Id:         *id,
		CanvasName: canvas.Name,
		FolderId:   canvas.FolderId,
		UpdatedAt:  canvas.UpdatedAt,
		CreatedAt:  canvas.CreatedAt,
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

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	if err := h.canvasUC.Delete(context.TODO(), newId, int(userID)); err != nil {

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

func (h *handlers) Update(w http.ResponseWriter, r *http.Request) {

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

	file, _, err := r.FormFile("file")
	// if err != nil {
	// 	h.logger.Error("failed to read file", slog.Attr{
	// 		Key:   "error",
	// 		Value: slog.StringValue(err.Error()),
	// 	})

	// 	w.WriteHeader(http.StatusBadRequest)
	// 	render.JSON(w, r, response.Error("failed to read file"))
	// 	return
	// }

	name := r.FormValue("name")

	_, claims, _ := jwtauth.FromContext(r.Context())
	userID, ok := claims["id"].(float64)
	if !ok {
		http.Error(w, "user_id not found", http.StatusUnauthorized)
		return
	}

	var url string

	if file != nil {
		url = h.cfg.AWSConfig.SecretEndpoint + "/1/" + id
	}

	canvasFound, err := h.canvasUC.Update(context.TODO(), newId, int(userID), url, name, &file)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.logger.Error("no such canvas", slog.Attr{
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
			h.logger.Error("failed to update canvas", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})

			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to update canvas"))
			return
		}
	}
	// w.WriteHeader(http.StatusOK)

	if canvasFound == nil {
		w.WriteHeader(200)
		render.JSON(w, r, "ok")
	} else {
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, *canvasFound)
	}

}

// TODO: возможно перенести все взаимодействие с МЛ в отдельную сущность
func (h *handlers) ImageToText(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("failed to read file", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to read file"))
		return
	}

	string_id := r.FormValue("id")
	id, err := strconv.Atoi(string_id)
	if err != nil {
		h.logger.Error("failed to get id", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to get id"))
		return
	}

	url := h.cfg.AWSConfig.SecretEndpoint + "/1/" + "9999999"

	err = h.canvasUC.MLUpdate(context.TODO(), 9999999, &file)
	if err != nil {
		h.logger.Error("internal server error", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("internal server error"))
		return
	}

	postBody, _ := json.Marshal(map[string]string{
		"image_url": url,
	})

	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("https://ml.hooli-pishem.ru/predict/", "application/json", responseBody)
	if err != nil {
		h.logger.Error("ML service unavaliable", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("ML service unavaliable"))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("cannot read from ml service", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("cannot read from ml service"))
		return
	}

	var mlResp MLResponse
	err = json.Unmarshal(body, &mlResp)
	if err != nil {
		h.logger.Error("failed to parse ML response", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to parse ML response"))
		return
	}

	sb := mlResp.Text

	h.logger.Info("sb ", sb, "id", id)

	err = h.canvasUC.UpdateText(context.TODO(), id, sb)

	// var textResp MLRequest
	// textResp.Text = sb

	// response, _ := json.Marshal(map[string]string{
	// 	"text": sb,
	// })
	render.JSON(w, r, mlResp)
}

func (h *handlers) TextToImage(w http.ResponseWriter, r *http.Request) {

	var req MLRequest

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

	postBody, _ := json.Marshal(map[string]string{
		"text": req.Text,
	})

	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://194.87.252.210:8001/draw", "application/json", responseBody)
	if err != nil {
		h.logger.Error("ML service unavaliable", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("ML service unavaliable"))
		return
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	r
	// }

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("cant read body from ML", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("cant read body from ML"))
		return
	}

	// w.Header().Set("Content-Type", "image/svg+xml")
	// w.WriteHeader(http.StatusOK)
	// w.Write(data)
	// render.JSON(w, r, data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		h.logger.Error("Cannot write response to client", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
}

func (h *handlers) SoundPredict(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("audio_file")
	if err != nil {
		h.logger.Error("failed to read file", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to read file"))
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("audio_file", header.Filename)
	if err != nil {
		h.logger.Error("failed to create form file", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create form file"))
		return
	}

	_, err = io.Copy(part, file)
	if err != nil {
		h.logger.Error("failed to copy file", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to copy file"))
		return
	}

	err = writer.Close()
	if err != nil {
		h.logger.Error("failed to close multipart writer", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to close multipart writer"))
		return
	}

	req, err := http.NewRequest("POST", "https://sound.hooli-pishem.ru/predict/", &buf)
	if err != nil {
		h.logger.Error("failed to create request", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create request"))
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Error("ML service unavailable", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("ML service unavailable"))
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		h.logger.Error("can't read body from ML", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("can't read body from ML"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		h.logger.Error("Cannot write response to client", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
}
