package middleware

import (
	"context"
	"github/imfropz/go-ddd/common/util"
	"github/imfropz/go-ddd/internal/domain/repository"
	"net/http"
)

func AuthenticationHandler(next http.Handler, userRepository repository.UserRepository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := util.RemoveBearer(r.Header.Get("Authorization"))
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := util.ValidateAccessToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := userRepository.FindByEmail(claims.Email)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), util.AccessTokenClaims{}, util.AccessTokenClaims{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		}))
		next.ServeHTTP(w, r)
	})
}
