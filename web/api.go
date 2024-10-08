package web

import (
	"encoding/json"
	"io"
	"net/http"
)

type WebAPI[T any, R any] func(T) (R, error)
type StdAPI func(string) (string, error)

type Any struct{}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewResponse[T interface{}](sucess bool, message string, data T) Response {
	return Response{
		Success: sucess,
		Message: message,
		Data:    data,
	}
}

const INVALID_RESPONSE = `{"success":false}`

func NewAPI[T any, R any](callback WebAPI[T, R]) StdAPI {
	return func(input string) (string, error) {
		var value T
		err := json.Unmarshal([]byte(input), &value)
		if err != nil {
			return INVALID_RESPONSE, err
		}

		result, err := callback(value)
		message := ""
		if err != nil {
			message = err.Error()
		}

		resp := NewResponse(err == nil, message, result)
		resultJSON, err := json.Marshal(resp)
		if err != nil {
			return INVALID_RESPONSE, err
		}

		return string(resultJSON), nil
	}
}

func RegisterAPI[T any, R any](s *Server, endpoint string, callback WebAPI[T, R]) {
	api := NewAPI(callback)
	s.API(endpoint, api)
}

func apiHandler(api StdAPI, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		req := string(body)
		resp, err := api(req)
		if err != nil {
			http.Error(w, "Failed process endpoint", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
