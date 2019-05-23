package middleware

import (
	"net/http"

	"github.com/winded/tyomaa/backend/db"
	"github.com/winded/tyomaa/backend/util"
	"github.com/winded/tyomaa/backend/util/context"
	"github.com/winded/tyomaa/shared/api"
)

// Authentication middleware finds the User from the database
// associated with the token claim, if it exists.
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Get(r)
		if ctx.Claims.UserID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		var user db.User
		if err := db.Instance.First(&user, ctx.Claims.UserID).Error; err == nil {
			ctx.User = user
			r = context.Set(r, ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// Authorization rejects any request with non-valid token.
func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Get(r)
		if ctx.User.ID == 0 {
			util.WriteApiError(w, api.Error(http.StatusForbidden, "Login required"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
