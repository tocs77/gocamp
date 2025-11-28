package middlewares

import (
	"net/http"
	"slices"
	"strings"
)

type HPPOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	WhiteList                   []string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckBody && isNeedToCheck(r, options.CheckBodyOnlyForContentType) {
				//filter the body params
				filterBodyParams(r, options.WhiteList)
			}
			if options.CheckQuery && r.URL.Query() != nil {
				filterQueryParams(r, options.WhiteList)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isNeedToCheck(r *http.Request, contentType string) bool {
	if r.Method != http.MethodPost {
		return false
	}
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, whiteList []string) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	for k, v := range r.Form {
		if !slices.Contains(whiteList, k) {
			r.Form.Del(k)
			continue
		}
		if len(v) > 0 {
			r.Form.Set(k, v[0])
		}
	}

}

func filterQueryParams(r *http.Request, whiteList []string) {
	queryParams := r.URL.Query()
	for k, v := range queryParams {
		if !slices.Contains(whiteList, k) {
			queryParams.Del(k)
			continue
		}
		if len(v) > 0 {
			queryParams.Set(k, v[0])
		}
	}
	r.URL.RawQuery = queryParams.Encode()
}
