package utility

import (
	"net/http"
	"strconv"
)

func GetPaginationParams(r *http.Request) (int, int) {
	limit := 10
	page := 1
	limitStr := r.URL.Query()["limit"]
	if len(limitStr) > 0 {
		newLimit, err := strconv.Atoi(limitStr[0])

		if err == nil {
			limit = newLimit
		}
	}
	pageStr := r.URL.Query()["page"]
	if len(pageStr) > 0 {
		newPage, err := strconv.Atoi(pageStr[0])
		if err == nil {
			page = newPage
		}
	}
	return limit, page
}
