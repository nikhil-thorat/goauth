package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/nikhil-thorat/goauth/configs"
	"github.com/nikhil-thorat/goauth/internals/auth"
	"github.com/nikhil-thorat/goauth/types"
	"github.com/nikhil-thorat/goauth/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var user types.RegisterUserPayload

	err := utils.ParseJSON(r, &user)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("INVALID PAYLOAD : %v", errors))
		return
	}

	existingUser, _ := h.store.GetUserByEmail(user.Email)
	if existingUser != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("USER ALREADY EXISTS"))
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  hashedPassword,
	})

	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload

	err := utils.ParseJSON(r, &user)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(user)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("INVALID PAYLOAD : %v", errors))
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	if !auth.ComparePassword(u.Password, []byte(user.Password)) {
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("INVALID CREDENTIALS : %v", err))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})

}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(auth.UserKey)

	user, err := h.store.GetUserByID(userId.(int))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, user)

}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(auth.UserKey)

	var updateDetails types.UpdateUserDetailsPayload
	err := utils.ParseJSON(r, &updateDetails)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	err = utils.Validate.Struct(updateDetails)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("INVALID PAYLOAD : %v", errors))
		return

	}

	err = h.store.UpdateUserDetails(userId.(int), updateDetails.FirstName, updateDetails.LastName)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)

}

func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {

	var updatePassword types.UpdateUserPasswordPayload
	err := utils.ParseJSON(r, &updatePassword)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, err)
		return

	}

	err = utils.Validate.Struct(updatePassword)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorJSON(w, http.StatusBadRequest, fmt.Errorf("INVALID PAYLOAD : %v", errors))
		return
	}

	userId := r.Context().Value(auth.UserKey)

	user, err := h.store.GetUserByID(userId.(int))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return

	}

	isValidPassword := auth.ComparePassword(user.Password, []byte(updatePassword.OldPassword))
	if !isValidPassword {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, fmt.Errorf("INVALID CREDENTIALS"))
		return
	}

	hashedPassword, err := auth.HashPassword(updatePassword.NewPassword)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdateUserPassword(userId.(int), hashedPassword)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {

	userId := r.Context().Value(auth.UserKey)

	err := h.store.DeleteUser(userId.(int))
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}
