package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		requestBody string
		expr        string
		statusCode  int
	}{
		{
			name:        "invalid HTTP method",
			method:      "GET",
			requestBody: "",
			expr:        "",
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:        "invalid JSON request body",
			method:      "POST",
			requestBody: "invalid json",
			expr:        "",
			statusCode:  http.StatusInternalServerError,
		},
		{
			name:        "valid request with invalid expression",
			method:      "POST",
			requestBody: `{"expression": "invalid expr"}`,
			expr:        "invalid expr",
			statusCode:  http.StatusUnprocessableEntity,
		},
		{
			name:        "valid request with valid expression",
			method:      "POST",
			requestBody: `{"expression": "1+1"}`,
			expr:        "1+1",
			statusCode:  http.StatusOK,
		},
		{
			name:        "error handling for rpn.Calc() error",
			method:      "POST",
			requestBody: `{"expression": "invalid expr"}`,
			expr:        "invalid expr",
			statusCode:  http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, "/api/v1/calculate", bytes.NewBuffer([]byte(test.requestBody)))
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()
			CalculateHandler(w, req)

			if w.Code != test.statusCode {
				t.Errorf("expected status code %d, got %d", test.statusCode, w.Code)
			}
		})
	}
}
