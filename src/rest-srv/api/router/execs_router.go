package router

import (
	"net/http"
	"rest-srv/api/handlers"
)

func registerExecsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("GET /execs/", handlers.GetExecHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecHandler)
	mux.HandleFunc("POST /execs/", handlers.AddExecHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecHandler)
	mux.HandleFunc("PATCH /execs/", handlers.PatchExecHandler)
	mux.HandleFunc("DELETE /execs", handlers.DeleteExecHandler)
	mux.HandleFunc("DELETE /execs/", handlers.DeleteExecHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.GetExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.UpdateExecHandler)

	mux.HandleFunc("POST /execs/login", handlers.LoginExecHandler)
	mux.HandleFunc("POST /execs/logout", handlers.LogoutExecHandler)
	// mux.HandleFunc("POST /execs/forgotpassword", handlers.ForgotPasswordExecHandler)
	// mux.HandleFunc("POST /execs/resetpassword/reset/{resetToken}", handlers.ResetPasswordExecHandler)
}
