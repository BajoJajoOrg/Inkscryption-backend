package delivery

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	//"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	// . "github.com/BajoJajoOrg/Inkscryption-backend/configs"
	"github.com/BajoJajoOrg/Inkscryption-backend/images"
	structures "github.com/BajoJajoOrg/Inkscryption-backend/images"
	"github.com/BajoJajoOrg/Inkscryption-backend/images/usecase"
	requests "github.com/BajoJajoOrg/Inkscryption-backend/pkg"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/mailru/easyjson"
)

type UploadPayload struct {
	DataURL string `json:"dataURL"`
}

type ImageHandler struct {
	useCase images.UseCase
	mx      *http.ServeMux
}

func (deliver *ImageHandler) ListenAndServe() error {
	err := http.ListenAndServe(":6000", deliver.mx)
	if err != nil {
		return fmt.Errorf("listen and serve error: %w", err)
	}

	return nil
}

func GetApi(c *usecase.UseCase) *ImageHandler {
	api := &ImageHandler{
		useCase: c,
		mx:      http.NewServeMux(),
	}
	var apiPath = "/api/v1/"

	println("This is api path", apiPath)

	// получить все канвасы или канвасы отфильтрованные по дате
	api.mx.Handle(apiPath+"get", requests.AllowedMethodMiddleware(http.HandlerFunc(api.GetImageHandler()), hashset.New("GET")))
	// api.mx.Handle(apiPath+"getML", requests.AllowedMethodMiddleware(http.HandlerFunc(api.GetMLHandler()), hashset.New("POST")))

	// сохранить канвас
	api.mx.Handle(apiPath+"add", requests.AllowedMethodMiddleware(http.HandlerFunc(api.AddCanvasHandler()), hashset.New("POST")))

	// получить мл на картинку
	// поменять структуру запроса на мл
	api.mx.Handle(apiPath+"getML", requests.AllowedMethodMiddleware(http.HandlerFunc(api.GetMLHandler()), hashset.New("POST")))

	// апдейт существующего канваса
	api.mx.Handle(apiPath+"update", requests.AllowedMethodMiddleware(http.HandlerFunc(api.UpdateCanvasHandler()), hashset.New("POST")))

	// удалить канвас по имени
	api.mx.Handle(apiPath+"delete", requests.AllowedMethodMiddleware(http.HandlerFunc(api.DeleteCanvasHandler()), hashset.New("POST")))

	api.mx.Handle("/test", http.HandlerFunc(api.Test()))

	return api
}

// func GetApi(c *usecase.UseCase) *ImageHandler {
// 	api := &ImageHandler{
// 		useCase: c,
// 		mx:      http.NewServeMux(),
// 	}
// 	var apiPath = "/api/v1"

// 	println("This is api path", apiPath)

// }

func (deliver *ImageHandler) Test() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "vse ok")
	}
}

func (deliver *ImageHandler) GetImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		//logger := request.Context().Value(Logg).(Log)

		//var filterFields = filter.NewOptions(false, []filter.Field{})

		canvasId := request.FormValue("id")

		createdAt := request.URL.Query().Get("created_at")
		//fmt.Print(createdAt)
		var dates []string
		if createdAt == "" {
			dates = []string{}
		} else {
			dates = strings.Split(createdAt, ":")
		}

		canvasName := request.URL.Query().Get("name")

		images, err := deliver.useCase.GetImage(dates, canvasName, canvasId, request.Context())
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		//print(images)

		var canvases structures.Canvases
		if images == nil {
			canvases = structures.Canvases{
				Canvases: []structures.Canvas{},
			}
		} else {
			canvases = structures.Canvases{
				Canvases: images,
			}
		}

		//requests.SendSimpleResponse(respWriter, request, http.StatusOK, images)
		requests.SendResponse(respWriter, request, http.StatusOK, canvases)
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent image")
	}
}

func (deliver *ImageHandler) DeleteCanvasHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		var r images.CanvasRequest

		//fmt.Print("americayaa")

		body, err := io.ReadAll(request.Body)
		if err != nil {
			//log.Fatalf("Bad body %w", err.Error())
			return
		}

		err = easyjson.Unmarshal(body, &r)
		if err != nil {
			//log.Fatalf("Cant unmarshal body %w", err.Error())
			return
		}

		//fmt.Print(r.Name)

		userCanvas := images.Canvas{
			Id: r.Id,
		}

		fmt.Print(userCanvas.Id)

		err = deliver.useCase.DeleteCanvas(userCanvas, request.Context())
		if err != nil {
			//log.Fatalf("do not working %w", err.Error())
			return
		}

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "")
	}
}

func (deliver *ImageHandler) AddCanvasHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {

		name := request.FormValue("name")
		userCanvas := structures.Canvas{
			Name:   name,
			Update: time.Now(),
		}

		id, err := deliver.useCase.AddImage(userCanvas, request.Context())
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		respBody, err := json.Marshal(map[string]string{
			"id": strconv.Itoa(int(id)),
		})
		if err != nil {
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		responseBody := bytes.NewBuffer(respBody)

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, responseBody.String())
	}
}

func (deliver *ImageHandler) UpdateCanvasHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {

		id := request.FormValue("id")
		name := request.FormValue("name")

		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		img, _, err := request.FormFile("image")
		if err != nil {
			fmt.Print("err", err)
		}

		//fileType := handler.Header.Get("Content-Type")

		filename := "1/" + name
		objectURL := "https://bajojajo.hb.ru-msk.vkcloud-storage.ru/" + filename

		//fmt.Print(objectURL)

		canvas_id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		userCanvas := structures.Canvas{
			Id:     canvas_id,
			Name:   name,
			Url:    objectURL,
			Update: time.Now(),
		}

		err = deliver.useCase.UpdateCanvas(userCanvas, img, request.Context())
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, "")
	}
}

func (deliver *ImageHandler) GetMLHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		img, _, err := request.FormFile("image")
		if err != nil {
			fmt.Print("err", err)
		}

		//fileType := handler.Header.Get("Content-Type")

		filename := "1/" + "asdfjlkasdfqwerpiou123048WORKINGSTUFFFORML"
		objectURL := "https://bajojajo.hb.ru-msk.vkcloud-storage.ru/" + filename

		//fmt.Print(objectURL)

		userCanvas := structures.Canvas{
			Name:   filename,
			Url:    objectURL,
			Update: time.Now(),
		}

		err = deliver.useCase.AddML(userCanvas, img, request.Context())
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		postBody, _ := json.Marshal(map[string]string{
			"image_url": objectURL,
		})

		responseBody := bytes.NewBuffer(postBody)

		resp, err := http.Post("http://194.87.252.210:8000/predict/", "application/json", responseBody)
		if err != nil {
			fmt.Print("AHTUNG AHTUNG ZLUKEN SOBAKEN ZA YAYCEN KLAC KLAC")
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Print(err)
		}
		sb := string(body)

		requests.SendSimpleResponse(respWriter, request, http.StatusOK, sb)
	}
}

func NewImageDelivery(uc images.UseCase) *ImageHandler {
	return &ImageHandler{
		useCase: uc,
	}
}
