package api

import (
	"bulletin-board/internal/ad"
	"bulletin-board/internal/ad/dto"
	"bulletin-board/internal/ad/service"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ads, err := h.service.GetAll(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ads)
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

		oneAd, err := h.service.GetByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, ad.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(oneAd)
	}
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var requestAd dto.RequestAd
		err := json.NewDecoder(r.Body).Decode(&requestAd)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		userId, ok := r.Context().Value("user_id").(int)
		if !ok {
			writeJSONError(w, http.StatusForbidden, "forbidden")
			return
		}
		requestAd.UserID = userId

		ad, err := h.service.Create(r.Context(), requestAd)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ad)
	}
}

func (h *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var requestAd dto.RequestAd

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.NewDecoder(r.Body).Decode(&requestAd)

		updatedAd, err := h.service.Update(r.Context(), requestAd, id)
		if err != nil {
			if errors.Is(err, ad.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else if errors.Is(err, ad.ErrForbidden) {
				writeJSONError(w, http.StatusForbidden, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(updatedAd)
	}
}

func (h *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid id")
			return
		}
		err = h.service.Delete(r.Context(), id)
		if err != nil {
			if errors.Is(err, ad.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
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
