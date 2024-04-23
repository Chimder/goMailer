package main

import (
	"goMailer/handler"
	"log"
	"net/http"
	"os"

	_ "goMailer/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//		@title			MAilere
//		@version		1.0
//		@description	Mailer Api
//	  @BasePath	/

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4000", "http://localhost:3000"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	r.Mount("/swag", httpSwagger.WrapHandler)

	googleHandler := handler.GoogleHandler{}
	r.Route("/google", func(r chi.Router) {
		r.Post("/reg", googleHandler.RegGoogleAcc)
		r.Get("/delete", googleHandler.DeleteGoogleCookie)
		r.Get("/session", googleHandler.GetGoogleSession)
		r.Get("/messages", googleHandler.MessagesAndContent)
	})

	tempHandler := handler.TempHandler{}
	r.Route("/temp", func(r chi.Router) {
		r.Get("/reg", tempHandler.RegTempEmail)
		r.Get("/message", tempHandler.GetTempMessage)
		r.Get("/messages", tempHandler.GetTempMessages)
		// r.Get("/delete", googleHandler.DeleteGoogleCookie)
		// r.Get("/session", googleHandler.GetGoogleSession)
	})

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}
	server := http.Server{
		Addr:    ":" + PORT,
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
