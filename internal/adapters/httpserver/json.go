package httpserver

import (
	"encoding/json"
	"net/http"
)

type logger interface {
	Error(string, ...any)
}

func WriteJSON(
	status int,
	response any, writer http.ResponseWriter,
	logger logger,
) {
	payload, err := json.Marshal(response)
	if err != nil {
		logger.Error("json.Marshal:", "err", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(status)
	_, _ = writer.Write(payload)
}
