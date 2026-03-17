package router

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func respond(w http.ResponseWriter, payload interface{}, err error) {
	if err != nil {
		respondWithError(w, 500, err.Error())
	} else {
		respondWithJSON(w, 200, payload)
	}
}

func badRequest(w http.ResponseWriter, err error) {
	slog.Error("error processing request body", slog.Any("error", err))
	respondWithError(w, http.StatusBadRequest, err.Error())
}

func badIdentifierError(w http.ResponseWriter, val interface{}, err error) {
	slog.Error("bad request, could not parse identifier", slog.Any("identifier", val), slog.Any("error", err))
	respondWithError(w, http.StatusBadRequest, err.Error())
}

// respondwithError return error message
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"message": msg})
}

// respondWithJSON write json response format
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	slog.Info("response", slog.Any("payload", payload))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func paramToInt(param string) (int, error) {
	return strconv.Atoi(param)
}
