package main

import (
	"Outlawed/middleware"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type requestData struct {
	A        int    `json:"a"`
	B        int    `json:"b"`
	Operator string `json:"op"`
}

type dividendError struct {
	dividend float64
}

func (d dividendError) Error() string {
	return fmt.Sprintf("You can not divide %.2f by 0", d.dividend)
}

func handleOperation(w http.ResponseWriter, r *http.Request) {
	var err error
	var data requestData
	var statusCode int
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, `      Error parsing json Data: Not correct format. 
      Correct format:{"a":<int>,"b":<int>, "op":<str>}`,
			http.StatusConflict)
		return
	}

	var result int
	switch data.Operator {
	case "add", "+":
		result = data.A + data.B
	case "subs", "-":
		result = data.A - data.B
	case "mult", "x":
		result = data.A * data.B
	case "div", "/":
		if data.B == 0 {
			err = dividendError{dividend: float64(data.A)}
			statusCode = http.StatusForbidden
		} else {
			result = data.A / data.B
		}
	default:
		err = errors.New("Unknown Operator")
		statusCode = http.StatusBadRequest
	}

	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = fmt.Sprintf("The result of %d %s %d is equal to %d", data.A, data.Operator, data.B, result)
		statusCode = http.StatusOK
	}

	jsonRet := fmt.Sprintf(`{"message": %s}`, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(jsonRet))
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /operation", http.HandlerFunc(handleOperation))
	mux.HandleFunc("GET /operation", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Use: {'a':<int>,'b':<int>, 'op':<str>"}`))
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: middleware.Logging(mux),
	}

	server.ListenAndServe()
}
