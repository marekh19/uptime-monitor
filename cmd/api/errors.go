package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: at path: %s via method: %s, error: %s", r.URL.Path, r.Method, err)

	writeJSONError(w, http.StatusInternalServerError, "The server has encountered a problem and could not complete your request.")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request: at path: %s via method: %s, error: %s", r.URL.Path, r.Method, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("resource not found: at path: %s via method: %s, error: %s", r.URL.Path, r.Method, err)

	writeJSONError(w, http.StatusNotFound, "The requested resource could not be found.")
}
