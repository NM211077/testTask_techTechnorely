package middleware

import (
	//"errors"
	"net/http"

	//"github.com/NM211077/testTask_techTechnorely/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
