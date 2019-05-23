package middleware

import (
	"net/http"

	"github.com/winded/tyomaa/backend/util"
	"github.com/winded/tyomaa/backend/util/context"
	"github.com/winded/tyomaa/backend/util/token"
	"github.com/winded/tyomaa/shared/api"
)

func TokenSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tk := r.Header.Get("X-Access-Token")
		if tk == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.Get(r)

		claims, err := token.Verify(tk)
		if err != nil {
			util.WriteApiError(w, api.Error(http.StatusBadRequest, "Invalid or malformed token"))
			return
		}

		ctx.Token = tk
		ctx.Claims = claims

		r = context.Set(r, ctx)
		next.ServeHTTP(w, r)
	})
}
