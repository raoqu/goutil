package web

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type Server struct {
	Local             bool
	Port              int
	Dir               string
	Endpoint          map[string]StdAPI
	WSSEndpointCreate map[string]WSSInstance
}

func NewServer(port int, local bool, dir string) *Server {
	server := Server{
		Local:             local,
		Port:              port,
		Dir:               dir,
		Endpoint:          make(map[string]StdAPI),
		WSSEndpointCreate: make(map[string]WSSInstance),
	}
	return &server
}

func (s *Server) API(endpoint string, api StdAPI) {
	s.Endpoint[endpoint] = api
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	// s.Mux = mux
	mux.HandleFunc("/api/", func(response http.ResponseWriter, request *http.Request) {
		dynamicHandler(s, response, request)
	})
	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		staticHandler(s, response, request)
	})

	go http.ListenAndServe(s.Address(), mux)
}

func (s *Server) Address() string {
	address := ""
	if s.Local {
		address = "127.0.0.1"
	}
	return fmt.Sprintf("%s:%d", address, s.Port)
}

func dynamicHandler(s *Server, w http.ResponseWriter, r *http.Request) {
	endpoint := r.URL.Path[len("/api/"):]
	if endpoint == "" {
		http.Error(w, "Endpoint is required", http.StatusBadRequest)
		return
	}

	if api, exists := s.Endpoint[endpoint]; exists {
		apiHandler(api, w, r)
	} else if wss, exists := s.WSSEndpointCreate[endpoint]; exists {
		wssHandler(wss, w, r)
	} else {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}
}

func staticHandler(s *Server, w http.ResponseWriter, r *http.Request) {
	assetsDir := s.Dir
	if r.URL.Path == "/" {
		http.ServeFile(w, r, filepath.Join(assetsDir, "index.html"))
		return
	}

	filePath := filepath.Join(assetsDir, r.URL.Path)
	contentType := getContentType(filePath)
	w.Header().Set("Content-Type", contentType)
	http.ServeFile(w, r, filePath)
}

func getContentType(filePath string) string {
	switch filepath.Ext(filePath) {
	case ".js":
		return "application/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}
