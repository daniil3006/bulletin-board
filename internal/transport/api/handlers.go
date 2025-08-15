package api

import (
	"bulletin-board/internal/ad"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func GetAll(store ad.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ads, err := store.GetAll(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ads)
	}
}

func GetByID(store ad.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		oneAd, err := store.GetByID(r.Context(), id)
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

func Create(store ad.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var ad ad.Ad
		err := json.NewDecoder(r.Body).Decode(&ad)

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid id")
		}

		if err = isValidateAd(ad); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid ad")
			return
		}

		ad, err = store.Create(r.Context(), ad)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create ad")
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ad)
	}
}

func Update(store ad.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var newAdd ad.Ad

		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.NewDecoder(r.Body).Decode(&newAdd)

		if err = isValidateAd(newAdd); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid ad")
			return
		}

		updatedAd, err := store.Update(r.Context(), newAdd, id)
		if err != nil {
			if errors.Is(err, ad.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(updatedAd)
	}
}

func Delete(store ad.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id, err := strconv.Atoi(params["id"])
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid id")
			return
		}
		err = store.Delete(r.Context(), id)
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

func isValidateAd(checkAd ad.Ad) error {
	if checkAd.Price <= 0 {
		return ad.ErrInvalidAd
	}
	return nil
}
