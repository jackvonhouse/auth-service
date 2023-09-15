package transport

import (
	"encoding/json"
	"net/http"
)

func Response(
	w http.ResponseWriter,
	data interface{},
) {

	if err := json.NewEncoder(w).Encode(&data); err != nil {
		Error(
			w,
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		)
	}
}

func Error(
	w http.ResponseWriter,
	statusCode int,
	message string,
) {

	w.WriteHeader(statusCode)

	if message != "" {
		json.NewEncoder(w).Encode(
			map[string]string{
				"error": message,
			},
		)
	}
}
