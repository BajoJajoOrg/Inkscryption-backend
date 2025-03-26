package requests

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mailru/easyjson"
)

func SendSimpleResponse(w http.ResponseWriter, _ *http.Request, code int, Body string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Csrft")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Access-Control-Max-Age", "86400")

	w.WriteHeader(code)
	if _, err := w.Write([]byte(Body)); err != nil {
		//logrus.Info(err.Error())
		fmt.Print("FUCK YOU")
		return
	}
}

func SendResponse[T easyjson.Marshaler](w http.ResponseWriter, r *http.Request, code int, Body T) {
	jsonResponse, err := easyjson.Marshal(Body)
	if err != nil {
		//logrus.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Csrft")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Access-Control-Max-Age", "86400")

	w.WriteHeader(code)
	_, err = w.Write(jsonResponse)
	if err != nil {
		//logrus.Info(err.Error())
		return
	}
}

func SendFileResponse(w http.ResponseWriter, file *os.File, filename string) error {
	// jsonResponse, err := easyjson.Marshal(Body)
	// if err != nil {
	// 	//logrus.Info(err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	if file == nil {
		return fmt.Errorf("file descriptor is nil")
	}

	// fileInfo, err := file.Stat()
	// if err != nil {
	// 	return fmt.Errorf("failed to get file info: %w", err)
	// }

	w.Header().Set("Content-Type", http.DetectContentType(make([]byte, 0, 512)))
	//w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Csrft")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Access-Control-Max-Age", "86400")

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek file: %w", err)
	}

	http.ServeContent(w, nil, filename, time.Now(), file)
	return nil
	// _, err := w.Write(jsonResponse)
	// if err != nil {
	// 	//logrus.Info(err.Error())
	// 	return
	// }
}
