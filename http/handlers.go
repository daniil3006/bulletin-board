package http

import (
	"bulletin-board/domain"
	"bulletin-board/storage"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func AllAds(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ads, err := store.GetAll()
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ads)
	}
}

func GetById(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		ID, err := strconv.ParseInt(params["id"], 10, 64)

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		ad, err := store.GetById(ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ad)
	}
}

func Create(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var ad domain.Ad
		err := json.NewDecoder(r.Body).Decode(&ad)

		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid id")
		}

		if err = isValidateAd(ad); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid ad")
			return
		}

		ad, err = store.Create(ad)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to create ad")
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ad)
	}
}

func Update(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var newAdd domain.Ad
		err := json.NewDecoder(r.Body).Decode(&newAdd)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err = isValidateAd(newAdd); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid ad")
			return
		}

		ad, err := store.Update(newAdd)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				writeJSONError(w, http.StatusNotFound, err.Error())
			} else {
				writeJSONError(w, http.StatusInternalServerError, err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ad)
	}
}

func Delete(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		ID, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid id")
			return
		}
		err = store.Delete(ID)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
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

func isValidateAd(ad domain.Ad) error {
	if ad.Price <= 0 {
		return storage.ErrInvalidAd
	}
	return nil
}
