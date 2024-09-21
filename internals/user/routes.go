package user

import (
	"net/http"

	"github.com/nikhil-thorat/goauth/internals/auth"
)

func (h *Handler) RegisterRoutes(router *http.ServeMux) {

	router.HandleFunc("POST /user/register", h.Register)
	router.HandleFunc("POST /user/login", h.Login)

	router.HandleFunc("GET /user/profile", auth.WithJWTAuth(h.Profile, h.store))
	router.HandleFunc("PUT /user/update", auth.WithJWTAuth(h.Update, h.store))
	router.HandleFunc("PUT /user/update-password", auth.WithJWTAuth(h.UpdatePassword, h.store))
	router.HandleFunc("DELETE /user/delete", auth.WithJWTAuth(h.Delete, h.store))
}
