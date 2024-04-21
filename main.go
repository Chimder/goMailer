package main

import (
	"goMailer/handler"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

//		@title			Mailer Api
//		@version		1.0
//		@description	manage your gmail or temp mail
//	  @BasePath	/

func main() {

	r := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4000", "http://localhost:3000"},
	})

	// handlerM := handler.NewMangaHandler(db, rdb)

	r.HandleFunc("GET /yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/swagger.yaml")
	})
	r.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// r.HandleFunc("GET /cookie/", func(w http.ResponseWriter, r *http.Request) {
	// 	cookie := http.Cookie{Name: "test", Value: "lxxx", Path: "/", MaxAge: 3600, Secure: true}
	// 	http.SetCookie(w, &cookie)

	// 	w.Write([]byte("set cccxc"))
	// })

	// r.HandleFunc("GET /cookieD/", func(w http.ResponseWriter, r *http.Request) {
	// 	c := http.Cookie{
	// 		Name:   "test",
	// 		Value:  "",
	// 		MaxAge: -1,
	// 	}
	// 	http.SetCookie(w, &c)
	// 	w.Write([]byte("dell cccxc"))
	// })

	r.HandleFunc("GET /google/reg", handler.RegAcc)

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
