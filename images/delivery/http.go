package delivery

import (
	"bytes"
	"strings"

	//"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	// . "github.com/BajoJajoOrg/Inkscryption-backend/configs"
	"github.com/BajoJajoOrg/Inkscryption-backend/images"
	structures "github.com/BajoJajoOrg/Inkscryption-backend/images"
	"github.com/BajoJajoOrg/Inkscryption-backend/images/usecase"
	requests "github.com/BajoJajoOrg/Inkscryption-backend/pkg"
	"github.com/emirpasic/gods/sets/hashset"
)

type UploadPayload struct {
	DataURL string `json:"dataURL"`
}

type ImageHandler struct {
	useCase images.UseCase
	mx      *http.ServeMux
}

func (deliver *ImageHandler) ListenAndServe() error {
	err := http.ListenAndServe(":8087", deliver.mx)
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

	api.mx.Handle(apiPath+"getImage", requests.AllowedMethodMiddleware(http.HandlerFunc(api.GetImageHandler()), hashset.New("GET")))
	api.mx.Handle(apiPath+"getML", requests.AllowedMethodMiddleware(http.HandlerFunc(api.GetMLHandler()), hashset.New("POST")))
	api.mx.Handle(apiPath+"add", requests.AllowedMethodMiddleware(http.HandlerFunc(api.AddImageHandler()), hashset.New("POST")))
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

		createdAt := request.URL.Query().Get("created_at")
		fmt.Print(createdAt)
		var dates []string
		if createdAt == "" {
			dates = []string{}
		} else {
			dates = strings.Split(createdAt, ":")
		}

		images, err := deliver.useCase.GetImage(1, dates, request.Context())
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusInternalServerError, err.Error())
			return
		}

		canvases := structures.Canvases{
			Canvases: images,
		}

		//requests.SendSimpleResponse(respWriter, request, http.StatusOK, images)
		requests.SendResponse(respWriter, request, http.StatusOK, canvases)
		//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Info("sent image")
	}
}

func (deliver *ImageHandler) AddImageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		err := request.ParseMultipartForm(10 << 20)
		if err != nil {
			//logger.Logger.WithFields(logrus.Fields{RequestID: logger.RequestID}).Warn(err.Error())
			requests.SendSimpleResponse(respWriter, request, http.StatusBadRequest, err.Error())
			return
		}

		img, handler, err := request.FormFile("image")

		//fileType := handler.Header.Get("Content-Type")

		filename := "1/" + fmt.Sprint(rand.Int()) + handler.Filename
		objectURL := "https://bajojajo.hb.ru-msk.vkcloud-storage.ru/" + filename

		fmt.Print(objectURL)

		userCanvas := structures.Canvas{
			Name:   filename,
			Url:    objectURL,
			Update: time.Now(),
		}

		err = deliver.useCase.AddImage(userCanvas, img, request.Context())
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

func (deliver *ImageHandler) GetMLHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(respWriter http.ResponseWriter, request *http.Request) {
		postBody, _ := json.Marshal(map[string]string{
			"image_url": "https://mumotiki.ru/sites/default/files/logokar3_0_0.png",
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

// func MetricTimeMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
// 		//start := time.Now()
// 		next.ServeHTTP(respWriter, request)
// 	})
// }
