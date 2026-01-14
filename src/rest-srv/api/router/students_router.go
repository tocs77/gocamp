package router

import (
	"net/http"
	"rest-srv/api/handlers"
)

func registerStudentRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students", handlers.AddStudentHandler)
	mux.HandleFunc("POST /students/", handlers.AddStudentHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("PATCH /students/", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students", handlers.DeleteStudentsHandler)
	mux.HandleFunc("DELETE /students/", handlers.DeleteStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteStudentHandler)
}
