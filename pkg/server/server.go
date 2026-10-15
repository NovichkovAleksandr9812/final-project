package server

import (
	"final-project/pkg/api"
	"net/http"
)

type Server struct {
	Server *http.Server
}

func CreateServer() *Server {
	mux := http.NewServeMux()

	api.Init(mux)
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	server := &http.Server{
		Addr:    ":7540",
		Handler: mux,
	}

	return &Server{
		Server: server,
	}
}
