package router

import (
	"net/http"
	"rest-srv/api/handlers"
)

func MainRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Register root handler
	mux.HandleFunc("GET /", handlers.RootHandler)

	// Register execs routes
	mux.HandleFunc("GET /execs", handlers.ExecHandler)
	mux.HandleFunc("POST /execs", handlers.ExecHandler)
	mux.HandleFunc("PUT /execs", handlers.ExecHandler)
	mux.HandleFunc("DELETE /execs", handlers.ExecHandler)

	// Register routes from other files
	registerStudentRoutes(mux)
	registerTeacherRoutes(mux)

	return mux
}
