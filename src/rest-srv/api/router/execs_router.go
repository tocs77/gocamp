package router

import (
	"net/http"
	"rest-srv/api/handlers"
)

func registerExecsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs", handlers.ExecHandler)
	mux.HandleFunc("GET /execs/", handlers.ExecHandler)
	mux.HandleFunc("POST /execs", handlers.ExecHandler)
	mux.HandleFunc("POST /execs/", handlers.ExecHandler)
	mux.HandleFunc("PATCH /execs", handlers.ExecHandler)
	mux.HandleFunc("PATCH /execs/", handlers.ExecHandler)
	mux.HandleFunc("DELETE /execs", handlers.ExecHandler)
	mux.HandleFunc("DELETE /execs/", handlers.ExecHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.ExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.ExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.ExecHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.ExecHandler)

	mux.HandleFunc("POST /execs/login", handlers.ExecHandler)
	mux.HandleFunc("POST /execs/logout", handlers.ExecHandler)
	mux.HandleFunc("POST /execs/forgotpassword", handlers.ExecHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetToken}", handlers.ExecHandler)

}
