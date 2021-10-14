/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.42
 */

package user

import (
	"encoding/json"
	"fmt"
	"github.com/keleeeep/test/internal/pkg/model"
	usecaseUser "github.com/keleeeep/test/internal/pkg/usecase/user"
	"github.com/keleeeep/test/internal/pkg/utils"
	"io/ioutil"
	"net/http"
)

type Handler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	CheckToken(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	uc usecaseUser.Usecase
}

func NewHandler(uc usecaseUser.Usecase) Handler {
	return &handler{uc: uc}
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("read request body failed: %v", err)

		utils.RespondErrWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		err = fmt.Errorf("unmarshal json failed: %v", err)

		utils.RespondErrWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := h.uc.CreateUser(r.Context(), &user)
	if err != nil {
		utils.RespondErrWithJSON(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("read request body failed: %v", err)

		utils.RespondErrWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	var user model.User
	err = json.Unmarshal(body, &user)
	if err != nil {
		err = fmt.Errorf("unmarshal json failed: %v", err)

		utils.RespondErrWithJSON(w, http.StatusInternalServerError, err)
		return
	}

	resp, err := h.uc.Login(r.Context(), &user)
	if err != nil {
		utils.RespondErrWithJSON(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *handler) CheckToken(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	resp, err := h.uc.CheckToken(&tokenString)
	if err != nil {
		utils.RespondErrWithJSON(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}