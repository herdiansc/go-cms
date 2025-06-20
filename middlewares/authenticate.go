package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/services"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Verifying Token.....")

		svc := services.NewTokenVerifyServices()
		code, res := svc.Verify(r.Header.Get("Authorization"))
		if code != http.StatusOK {
			log.Printf("Failed to verify token\n")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store the custom data in the request context
		ctx := context.WithValue(r.Context(), models.AuthVerifyCtxKey, res.Data)
		r = r.WithContext(ctx)

		log.Printf("Auth Res: %+v\n", res)

		next.ServeHTTP(w, r)
	})
}
