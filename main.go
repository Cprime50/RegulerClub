package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Cprime50/RegulerClub/database"
	"github.com/Cprime50/RegulerClub/migrations"
	routes "github.com/Cprime50/RegulerClub/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env file
	loadEnv()

	// load db config and connection
	loadDatabase()

	// close db
	defer database.DbClose()

	// migrate
	migrations.Migrate()

	// start server
	serveApplication()
}

func loadEnv() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Fatal("Error loading dev.env file", err)
	}
	log.Println("dev.env file loaded successfully")
}

func loadDatabase() {
	database.InitDb()
}

func serveApplication() {
	router := gin.Default()
	// import routes
	routes.SetupRoutes(router)

	srv := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		//log.Println("server running on https://localhost:8082")
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 2 seconds.")
	}
	log.Println("Server exiting")
}
