package main

import (
	"goMailer/handler"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//		@title			Manka Api
//		@version		1.0
//		@description	Manga search
//	 @BasePath	/

func main() {
	r := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4000", "http://localhost:3000"},
	})

	r.HandleFunc("GET /yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	r.HandleFunc("GET /swag/", httpSwagger.WrapHandler)

	r.HandleFunc("GET /google/reg", handler.RegGoogleAcc)
	r.HandleFunc("GET /google/delete", handler.DeleteGoogleCookie)
	r.HandleFunc("GET /google/session", handler.GetGoogleSession)

	// router.HandleFunc("GET /manga", handlerM.Manga)
	// router.HandleFunc("GET /manga/{name}/{chapter}", handlerM.Chapter)

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}
	server := http.Server{
		Addr: ":" + PORT,
		// Handler: middleware.Logging(c.Handler(router)),
		Handler: c.Handler(r),
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
