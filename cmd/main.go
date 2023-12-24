package main

import (
	"log"
	"net/http"

	"github.com/Isotton1/web-authenticatior/internal/handler"
	"github.com/Isotton1/web-authenticatior/pkg/database"
)

type Server struct{}

func (server Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/signup":
		handler.SignUpHandler(w, r)
	case "/login":
		handler.LogInHandler(w, r)
	default:
		handler.HomeHandler(w, r)
	}
}

func main() {
	dbUrl := "database.db"
	err := database.InitDB(dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	var server Server
	http.ListenAndServe(":3000", server)
}
