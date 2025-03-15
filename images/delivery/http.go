package delivery

import (
	"bytes"
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

		// cell := request.FormValue("cell")
		// println(cell)

		// userId := int64(request.Context().Value(RequestUserID).(int64))

		images, err := deliver.useCase.GetImage(1, request.Context())
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
			"image_url": "https://storage.googleapis.com/kagglesdsdata/datasets/1502872/3977616/test/test0.png?X-Goog-Algorithm=GOOG4-RSA-SHA256&X-Goog-Credential=databundle-worker-v2%40kaggle-161607.iam.gserviceaccount.com%2F20250315%2Fauto%2Fstorage%2Fgoog4_request&X-Goog-Date=20250315T114815Z&X-Goog-Expires=345600&X-Goog-SignedHeaders=host&X-Goog-Signature=9c0d7afd4dbe2663906e764d788959e89daa6b60258fd0444c676e99eb1d3af1dae1bbb1722e1d72c4d9935777b35c2cc58bf8c0c5b3ef2da7bf2c91266b62c35683f7cfdfd3821e54650641dbb9abb8183d8e696fe1a86bad79921d807e9da15439b6daa687587624c3a2b124e8c964ccd5969e57e7201d2a6b82a5f7dcd6a2acb5fbe655e6d19ede0f8ac159a29e7b9388957e667199cf3b7b58192451a22ae6498d0db76ab20e7ec80415f62d9978084bda7c530406203119317b9a867bd957f3269d77da1ff0c31a0a6e071f932221ce27eecfacec68250a6904caf6233bc08a4ad223b9b2276b7400ce17538d9327a570dec946ec75987499ac36f1f694",
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
