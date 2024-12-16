package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/exxxception/rpn/pkg/rpn"
)

type resultResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(data)
}

// CalculateHandler handles the /api/v1/calculate endpoint.
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var request struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("JSON:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	expr := strings.TrimSpace(request.Expression)

	result, err := rpn.Calc(expr)
	if err != nil {
		log.Println("expression is not valid:", err)
		if err := writeJSON(w, errorResponse{"Expression is not valid"}, http.StatusUnprocessableEntity); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	log.Println(expr, "expression calculate:", result)
	if err := writeJSON(w, resultResponse{result}, http.StatusOK); err != nil {
		log.Println("JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
