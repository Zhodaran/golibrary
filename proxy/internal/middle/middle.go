package middle

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"studentgit.kata.academy/Zhodaran/go-kata/controller"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/auth"
)

func TokenAuthMiddleware(resp controller.Responder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				resp.ErrorUnauthorized(w, errors.New("missing authorization token"))
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			_, err := auth.TokenAuth.Decode(token)
			if err != nil {
				resp.ErrorUnauthorized(w, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func ProxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
			return
		}
		proxyURL, _ := url.Parse("http://hugo:1313")
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)
		proxy.ServeHTTP(w, r)
	})
}
