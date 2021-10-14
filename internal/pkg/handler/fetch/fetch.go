/*
 * @Author: Fredy Gunawan
 * @Date: 14/10/21 13.48
 */

package fetch

import (
	usecaseFetch "github.com/keleeeep/test/internal/pkg/usecase/fetch"
	"github.com/keleeeep/test/internal/pkg/utils"
	"net/http"
)

type Handler interface {
	Fetch(w http.ResponseWriter, r *http.Request)
	Aggregate(w http.ResponseWriter, r *http.Request)
	CheckToken(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	uc usecaseFetch.Usecase
}

func NewHandler(uc usecaseFetch.Usecase) Handler {
	return &handler{uc: uc}
}

func (h *handler) Fetch(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	resp, err := h.uc.Fetch(&tokenString)
	if err != nil {
		utils.RespondErrWithJSON(w, http.StatusUnprocessableEntity, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func (h *handler) Aggregate(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")

	resp, err := h.uc.Aggregate(&tokenString)
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