package requests

import (
	"net/http"

	"github.com/emirpasic/gods/sets/hashset"
)

func AllowedMethodMiddleware(next http.Handler, methods *hashset.Set) http.Handler {
	return http.HandlerFunc(func(respWriter http.ResponseWriter, request *http.Request) {
		//log := request.Context().Value(Logg).(Log)

		if request.Method == http.MethodOptions {
			//log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("preflight")
			SendSimpleResponse(respWriter, request, http.StatusOK, "")
			return
		}
		if !methods.Contains(request.Method) {
			//log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("method not allowed")
			SendSimpleResponse(respWriter, request, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		//log.Logger.WithFields(logrus.Fields{RequestID: log.RequestID}).Info("methods checked")
		next.ServeHTTP(respWriter, request)
	})
}
