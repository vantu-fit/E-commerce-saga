package http

// import (
// 	"context"
// 	"errors"
// 	"net/http"
// 	"regexp"
// 	"time"

// 	"github.com/jackc/pgx/v5"
// 	"github.com/rs/zerolog/log"
// 	"github.com/vantu-fit/saga-pattern/internal/account/token"
// )

// const (
// 	AuthorizationHeader = "Authorization"
// 	AuthorizationPrefix = "Bearer"
// )

// var UrlCustomer = regexp.MustCompile(`^/api/v1/account/customer/([a-zA-Z0-9-]+)$`).MatchString

// func (app *HTTPGatewayServer) authMiddleware(next http.Handler, maker token.Maker) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Info().Msg("auth middleware :"+ r.URL.Path)
// 		if !UrlCustomer(r.URL.Path) {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			app.errJSON(w, errors.New("missing auth header"), http.StatusUnauthorized)
// 			return
// 		}

// 		if len(authHeader) <= len(AuthorizationPrefix) {
// 			app.errJSON(w, errors.New("invalid auth header"), http.StatusUnauthorized)
// 			return
// 		}

// 		token := authHeader[len(AuthorizationPrefix)+1:]
// 		if token == "" {
// 			app.errJSON(w, errors.New("missing token"), http.StatusUnauthorized)
// 			return
// 		}

// 		log.Info().Msg("token:"+ token)
// 		payload, err := maker.VerifyToken(token)
// 		if err != nil {
// 			app.errJSON(w, err, http.StatusUnauthorized)
// 			return
// 		}

// 		if payload.ID.String() == "" {
// 			app.errJSON(w, errors.New("missing session id"), http.StatusUnauthorized)
// 			return
// 		}

// 		session, err := app.store.GetSessionById(context.Background(), payload.ID)
// 		if err != nil {
// 			if err == pgx.ErrNoRows {
// 				app.errJSON(w, errors.New("session not found"), http.StatusUnauthorized)
// 				return
// 			}
// 			app.errJSON(w, err, http.StatusInternalServerError)
// 			return
// 		}

// 		if session.ExpiresAt.Before(time.Now()) {
// 			app.errJSON(w, errors.New("session expired"), http.StatusUnauthorized)
// 			return
// 		}

// 		if session.IsBlocked {
// 			app.errJSON(w, errors.New("session is blocked"), http.StatusUnauthorized)
// 			return
// 		}

// 		if session.UserID.String() != payload.UserID.String() {
// 			app.errJSON(w, errors.New("invalid session"), http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
