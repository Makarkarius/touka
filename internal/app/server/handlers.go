package server

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

func handleGet(s *server, w http.ResponseWriter, r *http.Request) {
	s.logger.Info("start handling get request")
	defer s.logger.Info("finish handling get request")

	switch r.Method {
	case http.MethodPost:
		request := DefaultBatch()
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			s.logger.Info("error decoding request", zap.Error(err))
			return
		}
		s.logger.Info("decoded batch request", zap.Any("request", request))
		s.requestQueue <- request
		http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusOK)
	default:
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusNotImplemented)
		s.logger.Info("incorrect method", zap.String("method", r.Method))
	}
}
