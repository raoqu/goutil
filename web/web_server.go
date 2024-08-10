package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WebAPI func(string) (string, error)

type Web struct {
	Address    string
	Port       int
	AssetsPath string
	Callbacks  map[string]WebAPI
}

func WebServer(address string, port int, path string) Web {
	var server = Web{
		Address:    address,
		Port:       port,
		AssetsPath: path,
		Callbacks:  make(map[string]WebAPI),
	}
	go http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)

	return server
}

func WebAPIRegister[T any, R any](w *Web, service string, callback func(T) (R, error)) {
	w.Callbacks[service] = func(param string) (string, error) {
		var obj T

		err := json.Unmarshal([]byte(param), &obj)
		if err != nil {
			return "", err
		}

		result, err := callback(obj)
		if err != nil {
			return "", err
		}

		response, err := json.Marshal(result)
		if err != nil {
			return "", err
		}

		return string(response), nil
	}
}
