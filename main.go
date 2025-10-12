package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/NurulloMahmud/article_hub/internal/app"
	"github.com/NurulloMahmud/article_hub/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "application port")
	flag.Parse()

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	// routes
	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	app.Logger.Printf("application running on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
