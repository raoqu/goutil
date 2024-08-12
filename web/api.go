package web

import "encoding/json"

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

func NewAPI[T any, R any](callback WebAPI[T, R]) func(string) (string, error) {
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
