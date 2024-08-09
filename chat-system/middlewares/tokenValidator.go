package middlewares

import (
	"context"
	"log"
	"net/http"

	"chat-system/database"
	"chat-system/types"
)

func TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		log.Println("Authorization token:", token)

		userId, valid, _ := database.ValidateToken(token)
		if token == "" || !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Println("Setting userID in context:", userId)
		ctx := context.WithValue(r.Context(), types.UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
