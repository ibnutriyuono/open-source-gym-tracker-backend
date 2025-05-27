package response

import (
	"encoding/json"
	"net/http"
)

type Meta struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message,omitempty"`
}

type JSONResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := JSONResponse{
		Meta: Meta{
			Success:    statusCode >= 200 && statusCode < 300,
			StatusCode: statusCode,
			Message:    message,
		},
		Data: data,
	}
	
	json.NewEncoder(w).Encode(resp)
}
