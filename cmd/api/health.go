package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := WriteJSON(w, http.StatusOK, "Ok")
	if err != nil {
		log.Print(err.Error())
	}
}
