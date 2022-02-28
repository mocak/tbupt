package json

import (
	"encoding/json"
	"io"
	"net/http"
)

// Response prepares json response
// Sets content type and response body
func Response(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(v)
	if err != nil {
		Error(w, "Unexpected Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(bytes)
}

// DecodeBody reads the json-encoded request body
// and stores it in the value pointed to by v
// if fails, writes error response and return error
func DecodeBody(w http.ResponseWriter, r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		switch err {
		case io.EOF:
			return nil
		default:
			Error(w, err.Error(), http.StatusBadRequest)
			return err
		}
	}
	return nil
}

// Error prepares json error response
func Error(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
