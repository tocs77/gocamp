package router

import (
	"net/http"
	"rest-srv/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	// Register routes with and without trailing slash
	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.AddTeacherHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeacherHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeacherHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)
	mux.HandleFunc("GET /teachers/{id}/students", handlers.GetTeacherStudentsHandler)
	mux.HandleFunc("GET /teachers/{id}/studentsCount", handlers.GetTeacherStudentsCountHandler)

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

	mux.HandleFunc("GET /execs", handlers.ExecHandler)
	mux.HandleFunc("POST /execs", handlers.ExecHandler)
	mux.HandleFunc("PUT /execs", handlers.ExecHandler)
	mux.HandleFunc("DELETE /execs", handlers.ExecHandler)
	mux.HandleFunc("GET /", handlers.RootHandler)
	return mux
}
