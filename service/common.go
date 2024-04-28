package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func WriteJsonResponseWithCode(w http.ResponseWriter, resp []byte, code int) {
	w.WriteHeader(code)
	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, string(resp))
}

func WriteJsonErrorResponseWithCode(w http.ResponseWriter, err error, msg string, code int) {
	out, err := json.Marshal(&ErrorResponse{
		Error:   err.Error(),
		Message: msg,
	})
	if err != nil {
		log.Println("[Error] Failed to process request", err)
	}

	WriteJsonResponseWithCode(w, out, code)
}
