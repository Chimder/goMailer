package main

import (
	"context"
	"goMailer/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "goMailer/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title MAilere
// @version 1.0
// @description Mailer Api
// @BasePath /

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	googleHandler := handler.NewGoogleHandler()
	r.Route("/google", func(r chi.Router) {
		r.Post("/reg", googleHandler.RegGoogleAcc)
		r.Get("/delete", googleHandler.DeleteGoogleCookie)
		r.Get("/session", googleHandler.GetGoogleSession)
		r.Get("/messages", googleHandler.MessagesAndContent)
	})

	tempHandler := handler.NewTempHandler()
	r.Route("/temp", func(r chi.Router) {
		r.Get("/reg", tempHandler.RegTempEmail)
		r.Get("/message", tempHandler.GetTempMessage)
		r.Get("/messages", tempHandler.GetTempMessages)
		r.Get("/session", tempHandler.GetTempSession)
		r.Delete("/delete", tempHandler.DeleteTempSession)
	})

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "4000"
	}

	server := &http.Server{
		Addr:    ":" + PORT,
		Handler: r,
	}
	go shutdownMonitor(ctx, server)

	log.Println("Listening...")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", PORT, err)
	}
}
func shutdownMonitor(ctx context.Context, server *http.Server) {
	c, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Println("Waiting for shutdown signal...")
	select {
	case <-c.Done():
		log.Println("Shutdown signal received, shutting down server...")
	case <-ctx.Done():
		log.Println("Context done, shutting down server...")
	}

	ctxShutdown, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

// func shutdownMonitor(ctx context.Context, server *http.Server) {
// 	c, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
// 	defer stop()

// 	<-c.Done()
// 	log.Println("Shutting down gracefully, press Ctrl+C again to force")

// 	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	log.Println("Attempting to shutdown server...")
// 	if err := server.Shutdown(ctxShutdown); err != nil {
// 		log.Printf("Error shutting down server: %v", err)
// 		log.Fatalf("Server forced to shutdown: %v", err)
// 	}

// 	log.Println("Server exiting")
// }
