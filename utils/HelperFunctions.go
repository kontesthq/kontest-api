package utils

import (
	"context"
	"fmt"
	"kontest-api/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// InitalizeDatabase Initialize the database connection with default values
func InitalizeDatabase(
	dbPort string,
	dbHost string,
	user string,
	password string,
	dbname string,
	sslmode string,
) {

	if dbErr := database.Connect(dbname, dbPort, dbHost, user, password, sslmode); dbErr != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", dbErr)
		return
	}
}

func HandleShutdown(server *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Shutting down server...")
		// Create a context with a timeout to allow for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Server shutdown failed: %+v", err)
		}
		log.Println("Server exited properly")
	}()
}
