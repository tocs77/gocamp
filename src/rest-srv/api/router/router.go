package router

import (
	"net/http"
	"rest-srv/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/teachers/", handlers.TeacherHandler)
	mux.HandleFunc("/teachers", handlers.TeacherHandler)
	mux.HandleFunc("/students", handlers.StudentHandler)
	mux.HandleFunc("/execs", handlers.ExecHandler)
	mux.HandleFunc("/", handlers.RootHandler)
	return mux
}
