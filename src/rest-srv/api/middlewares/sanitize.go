package middlewares

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"rest-srv/utility"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//sanitize the path
		sanitizedPath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.URL.Path = sanitizedPath.(string)

		//sanitize the query params
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for k, v := range params {
			sanitizedKey, err := clean(k)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var sanitizedValues []string
			for _, val := range v {
				sanitizedValue, err := clean(val)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				sanitizedValues = append(sanitizedValues, sanitizedValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
		}
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()

		//sanitize request body
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type not supported", http.StatusUnsupportedMediaType)
			return
		}
		if r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bodyString := strings.TrimSpace(string(bodyBytes))
			if len(bodyString) != 0 {
				sanitizedBody, err := clean(bodyString)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(strings.NewReader(sanitizedBody.(string)))
			}
		}
		next.ServeHTTP(w, r)
	})
}

func clean(data any) (any, error) {
	switch v := data.(type) {
	case map[string]any:
		for k, val := range v {
			v[k] = sanitizeValue(val)
		}
		return v, nil
	case []any:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:
		return nil, utility.ErrorHandler(fmt.Errorf("invalid data type: %T", data), "internal server error")
	}
}

func sanitizeValue(data any) any {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)
	case map[string]any:
		for k, val := range v {
			v[k] = sanitizeValue(val)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
		return v
	default:
		return v
	}
}

func sanitizeString(data string) string {
	return bluemonday.UGCPolicy().Sanitize(data)
}
