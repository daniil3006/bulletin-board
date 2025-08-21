package api

import (
	"bulletin-board/internal/user"
	"bulletin-board/internal/user/dto"
	"bulletin-board/internal/user/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: &service}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		users, err := h.service.GetAll(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(users)
	}
}

func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		oneUser, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(oneUser)
	}
}

func (h *Handler) GetUsersAds() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		ads, err := h.service.GetUsersAds(r.Context(), id)
		if err != nil {
			if errors.Is(err, user.ErrUserNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ads)
	}
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var requestUser dto.RequestUser
		err := json.NewDecoder(r.Body).Decode(&requestUser)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		responseUser, err := h.service.Create(r.Context(), requestUser)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(responseUser)
	}
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var requestUser dto.RequestUser
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.NewDecoder(r.Body).Decode(&requestUser)

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := h.service.Update(r.Context(), requestUser, id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		err = h.service.Delete(r.Context(), id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
	log.Printf("Status: %d | Message: %s", status, message)
}
