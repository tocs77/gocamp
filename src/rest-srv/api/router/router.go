package router

import (
	"net/http"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	registerStudentRoutes(mux)
	registerTeacherRoutes(mux)
	registerExecsRoutes(mux)

	return mux
}
