package api

import (
	"encoding/json"
	"net/http"
)

func sendJSON(resp http.ResponseWriter, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	resp.Header().Set("Content-Type", "application/json")
	_, err = resp.Write(data)
	if err != nil {
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

func sendJSONError(resp http.ResponseWriter, obj interface{}, code int) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Error(resp, string(data), code)
}
