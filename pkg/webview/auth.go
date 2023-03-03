package webview

import (
	"context"
	"net/http"
)

type key int

const (
	keyUserId key = iota
)

func (s *Server) isAuthorized(endpoint http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("mvcAuthToken")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := s.service.ParseToken(context.Background(), cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
		endpoint(w, r.WithContext(context.WithValue(r.Context(), keyUserId, id)))
	}
}
