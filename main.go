package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//		@title			Manka Api
//		@version		1.0
//		@description	Manga search
//	 @BasePath	/

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	router := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4000", "http://localhost:3000", "https://golang-on-koyeb-mankago.koyeb.app/", "https://manka-next.vercel.app"},
	})

	// handlerM := handler.NewMangaHandler(db, rdb)
	router.HandleFunc("GET /yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	router.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	// router.HandleFunc("GET /mangas", handlerM.Mangas)
	// router.HandleFunc("GET /manga", handlerM.Manga)
	// router.HandleFunc("GET /manga/{name}/{chapter}", handlerM.Chapter)

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}
	server := http.Server{
		Addr: ":" + PORT,
		// Handler: middleware.Logging(c.Handler(router)),
		Handler: c.Handler(router),
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
