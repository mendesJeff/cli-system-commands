package responses

import (
	"encoding/json"
	"log"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, datas interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if datas != nil {
		if error := json.NewEncoder(w).Encode(datas); error != nil {
			log.Fatal(error)
		}
	}
}

func Error(w http.ResponseWriter, statusCode int, error error) {
	JSON(w, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: error.Error(),
	})

}
